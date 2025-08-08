package interceptors

import (
	"context"

	pb "github.com/ramil063/gometrics/internal/grpc/proto"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/security/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// DecryptUnaryInterceptor расшифровывает входящие gRPC сообщения
func DecryptUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Если дешифратор не настроен, пропускаем
	if crypto.GRPCDecryptor == nil {
		return handler(ctx, req)
	}

	request := req.(*pb.ListMetricsRequest)
	// Дешифруем данные
	decryptedData, err := crypto.GRPCDecryptor.Decrypt(request.GetCryptoMetrics())
	if err != nil {
		logger.WriteErrorLog(err.Error(), "Decryption failed")
		return nil, status.Errorf(codes.InvalidArgument, "decryption failed")
	}

	// Десериализуем оригинальный запрос
	var originalReq pb.ListMetricsRequest
	if err = proto.Unmarshal(decryptedData, &originalReq); err != nil {
		logger.WriteErrorLog(err.Error(), "Failed to unmarshal decrypted data")
		return nil, status.Errorf(codes.InvalidArgument, "invalid request format")
	}

	// Вызываем обработчик с дешифрованным запросом
	return handler(ctx, &originalReq)
}
