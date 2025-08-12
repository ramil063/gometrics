package interceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/hash"
)

// HashCheckUnaryInterceptor проверяет хеш входящих данных и добавляет хеш к исходящим
func HashCheckUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Если хеш-ключ не установлен, пропускаем проверку
	if handlers.HashKey == "" {
		return handler(ctx, req)
	}
	// Получаем метаданные из контекста
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Error(codes.InvalidArgument, "grpc: metadata is required")
	}

	// 1. Проверка входящего хеша
	if reqBytes, ok := req.([]byte); ok {
		// Получаем хеш из заголовков
		headerHashSHA256 := getFirstValue(md, "hashsha256")
		if headerHashSHA256 == "" {
			return nil, status.Error(codes.InvalidArgument, "grpc: hash is empty")
		}

		// Вычисляем хеш тела запроса
		bodyHashSHA256 := hash.CreateSha256(reqBytes, handlers.HashKey)
		if headerHashSHA256 != bodyHashSHA256 {
			return nil, status.Error(codes.InvalidArgument, "grpc: hash isn't correct")
		}
	}

	// Вызываем обработчик
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}

	// 2. Добавление хеша к исходящим данным
	if respBytes, ok := resp.([]byte); ok {
		// Вычисляем хеш ответа
		respHash := hash.CreateSha256(respBytes, handlers.HashKey)

		// Устанавливаем заголовок с хешем
		header := metadata.Pairs("hashsha256", respHash)
		if err = grpc.SetHeader(ctx, header); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to set header: %v", err)
		}
	}

	return resp, nil
}

// getFirstValue получает первое значение из метаданных по ключу
func getFirstValue(md metadata.MD, key string) string {
	key = strings.ToLower(key)
	if vals := md.Get(key); len(vals) > 0 {
		return vals[0]
	}
	return ""
}
