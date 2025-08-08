package interceptors

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// mockUnaryHandler для тестирования
type mockUnaryHandler struct {
	resp interface{}
	err  error
}

func (m *mockUnaryHandler) handle(ctx context.Context, req interface{}) (interface{}, error) {
	return m.resp, m.err
}

// createTestContext создает контекст с тестовыми метаданными
func createTestContext(ip string) context.Context {
	ctx := context.Background()
	if ip != "" {
		ctx = peer.NewContext(ctx, &peer.Peer{
			Addr: &net.TCPAddr{IP: net.ParseIP(ip), Port: 12345},
		})
	}
	return ctx
}

func TestNewTrustedIPInterceptor(t *testing.T) {
	tests := []struct {
		name          string
		trustedSubnet string
		clientIP      string
		errContains   string
		wantErr       bool
		errCode       codes.Code
	}{
		{
			name:          "empty trusted subnet allows all",
			trustedSubnet: "",
			clientIP:      "192.168.1.100",
			wantErr:       false,
		},
		{
			name:          "trusted IP passes",
			trustedSubnet: "192.168.1.0/24",
			clientIP:      "192.168.1.100",
			wantErr:       false,
		},
		{
			name:          "untrusted IP blocked",
			trustedSubnet: "192.168.1.0/24",
			clientIP:      "10.0.0.1",
			wantErr:       true,
			errCode:       codes.PermissionDenied,
			errContains:   "is not trusted",
		},
		{
			name:          "no IP in context fails",
			trustedSubnet: "192.168.1.0/24",
			clientIP:      "",
			wantErr:       true,
			errCode:       codes.PermissionDenied,
			errContains:   "failed to get client IP",
		},
		{
			name:          "invalid IP format fails",
			trustedSubnet: "192.168.1.0/24",
			clientIP:      "invalid-ip",
			wantErr:       true,
			errCode:       codes.PermissionDenied,
		},
		{
			name:          "IPv6 trusted",
			trustedSubnet: "2001:db8::/32",
			clientIP:      "2001:db8::1",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := NewTrustedIPInterceptor(tt.trustedSubnet)
			handler := &mockUnaryHandler{resp: "response", err: nil}

			ctx := createTestContext(tt.clientIP)
			resp, err := interceptor(ctx, "request", nil, handler.handle)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errCode != codes.OK {
					assert.Equal(t, tt.errCode, status.Code(err))
				}
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, "response", resp)
		})
	}
}

func TestGetClientIP_FromMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]string
		expected string
	}{
		{
			name:     "single x-real-ip",
			metadata: map[string]string{"x-real-ip": "192.168.1.1"},
			expected: "192.168.1.1",
		},
		{
			name:     "multiple x-real-ip values",
			metadata: map[string]string{"x-real-ip": "10.0.0.1, 10.0.0.2"},
			expected: "10.0.0.1, 10.0.0.2",
		},
		{
			name:     "case insensitive header",
			metadata: map[string]string{"X-Real-Ip": "172.16.0.1"},
			expected: "172.16.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := metadata.New(tt.metadata)
			ctx := metadata.NewIncomingContext(context.Background(), md)

			ip, err := getClientIP(ctx)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, ip)
		})
	}
}

func TestGetClientIP_Errors(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		errMsg string
	}{
		{
			name:   "no metadata no peer",
			ctx:    context.Background(),
			errMsg: "could not get peer info",
		},
		{
			name:   "invalid peer addr type",
			ctx:    peer.NewContext(context.Background(), &peer.Peer{Addr: mockAddr{}}),
			errMsg: "could not extract IP from peer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, err := getClientIP(tt.ctx)

			assert.Equal(t, "", ip)
			assert.Error(t, err)
			assert.Equal(t, codes.PermissionDenied, status.Code(err))
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

type mockAddr struct{}

func (m mockAddr) Network() string { return "mock" }
func (m mockAddr) String() string  { return "mock" }

func Test_isIPTrusted(t *testing.T) {
	type args struct {
		trustedSubnet string
		ipStr         string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "trusted IP passes",
			args: args{
				trustedSubnet: "192.168.1.0/24",
				ipStr:         "192.168.1.100",
			},
			want: true,
		},
		{
			name: "trusted IP passes",
			args: args{
				trustedSubnet: "192.168.1.0/24",
				ipStr:         "192.168.2.100",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isTrusted := isIPTrusted(tt.args.trustedSubnet, tt.args.ipStr)
			assert.Equal(t, tt.want, isTrusted)
		})
	}
}
