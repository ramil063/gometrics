package interceptors

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func NewTrustedIPInterceptor(trustedSubnet string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Если не задана доверенная подсеть, пропускаем проверку
		if trustedSubnet == "" {
			return handler(ctx, req)
		}

		// Получаем IP из контекста
		clientIP, err := getClientIP(ctx)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "failed to get client IP: %v", err)
		}

		// Проверяем IP
		if !isIPTrusted(trustedSubnet, clientIP) {
			return nil, status.Errorf(codes.PermissionDenied, "IP %s is not trusted", clientIP)
		}

		return handler(ctx, req)
	}
}

// getClientIP извлекает IP адрес клиента из контекста
func getClientIP(ctx context.Context) (string, error) {
	// 1. Проверяем X-Real-IP в метаданных (как в HTTP)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if realIPs := md.Get("x-real-ip"); len(realIPs) > 0 {
			return realIPs[0], nil
		}
	}

	// 2. Если нет в метаданных, берем из peer информации
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", status.Error(codes.PermissionDenied, "could not get peer info")
	}

	if addr, ok := p.Addr.(*net.TCPAddr); ok {
		return addr.IP.String(), nil
	}

	return "", status.Error(codes.PermissionDenied, "could not extract IP from peer")
}

// isIPTrusted проверяет принадлежность IP к доверенной подсети
func isIPTrusted(trustedSubnet, ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	_, subnet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return false
	}

	return subnet.Contains(ip)
}
