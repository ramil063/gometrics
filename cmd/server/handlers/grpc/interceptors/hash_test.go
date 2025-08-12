package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/hash"
)

// Вспомогательная функция для создания контекста с метаданными
func createContextWithMetadata(md metadata.MD) context.Context {
	return metadata.NewIncomingContext(context.Background(), md)
}

func TestHashCheckUnaryInterceptor(t *testing.T) {
	// Устанавливаем тестовый хеш-ключ
	originalHashKey := handlers.HashKey
	handlers.HashKey = "test-secret-key"
	defer func() { handlers.HashKey = originalHashKey }()

	tests := []struct {
		ctx           context.Context
		req           interface{}
		handlerResp   interface{}
		handlerErr    error
		name          string
		wantErrCode   codes.Code
		wantErr       bool
		wantHeader    bool
		checkResponse bool
	}{
		{
			name: "string response without hash",
			ctx: createContextWithMetadata(metadata.Pairs(
				"hashsha256", hash.CreateSha256([]byte("test-data"), "test-secret-key"),
			)),
			req:           []byte("test-data"),
			handlerResp:   "response",
			handlerErr:    nil,
			wantErr:       false,
			wantHeader:    false,
			checkResponse: true,
		},
		{
			name:        "missing metadata",
			ctx:         context.Background(),
			req:         []byte("test-data"),
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "empty hash",
			ctx: createContextWithMetadata(metadata.Pairs(
				"hashsha256", "",
			)),
			req:         []byte("test-data"),
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "invalid hash",
			ctx: createContextWithMetadata(metadata.Pairs(
				"hashsha256", "invalid-hash",
			)),
			req:         []byte("test-data"),
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "non-bytes request",
			ctx: createContextWithMetadata(metadata.Pairs(
				"hashsha256", hash.CreateSha256([]byte("test-data"), "test-secret-key"),
			)),
			req:         "string-request",
			handlerResp: "string-response",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return tt.handlerResp, tt.handlerErr
			}

			resp, err := HashCheckUnaryInterceptor(tt.ctx, tt.req, nil, handler)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrCode != codes.OK {
					assert.Equal(t, tt.wantErrCode, status.Code(err))
				}
				return
			}
			require.NoError(t, err)

			if tt.checkResponse {
				assert.Equal(t, tt.handlerResp, resp)
			}

			if tt.wantHeader {
				md, ok := metadata.FromOutgoingContext(tt.ctx)
				require.True(t, ok)

				var expectedHash string
				switch v := tt.handlerResp.(type) {
				case []byte:
					expectedHash = hash.CreateSha256(v, handlers.HashKey)
				case string:
					expectedHash = hash.CreateSha256([]byte(v), handlers.HashKey)
				}
				assert.Equal(t, []string{expectedHash}, md.Get("hashsha256"))
			}
		})
	}
}

func Test_getFirstValue(t *testing.T) {
	type args struct {
		md  metadata.MD
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 1",
			args: args{
				md: metadata.MD{
					"hashsha256": {"test-data", "test-secret-key", "test-data-2"},
				},
				key: "hashsha256",
			},
			want: "test-data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFirstValue(tt.args.md, tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}
