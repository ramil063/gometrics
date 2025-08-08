package interceptors

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/ramil063/gometrics/internal/grpc/proto"
	"github.com/ramil063/gometrics/internal/security/crypto"
)

// mockDecryptor для тестирования
type mockDecryptor struct {
	decryptFunc func([]byte) ([]byte, error)
}

func (m *mockDecryptor) Decrypt(data []byte) ([]byte, error) {
	return m.decryptFunc(data)
}

// mockHandler для тестирования
type mockHandler struct {
	resp interface{}
	err  error
}

func (m *mockHandler) handle(ctx context.Context, req interface{}) (interface{}, error) {
	return m.resp, m.err
}

func TestDecryptUnaryInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Тестовые данные
	testMetrics := &pb.ListMetricsRequest{
		Metrics: []*pb.Metric{
			{Id: "cpu", Value: 42.5},
		},
	}
	encryptedData, _ := proto.Marshal(testMetrics)

	tests := []struct {
		decryptor     crypto.Decryptor
		req           interface{}
		handlerResp   interface{}
		handlerErr    error
		name          string
		wantErrCode   codes.Code
		wantErr       bool
		checkResponse bool
	}{
		{
			name:          "no decryptor skips decryption",
			decryptor:     nil,
			req:           testMetrics,
			handlerResp:   testMetrics,
			handlerErr:    nil,
			wantErr:       false,
			checkResponse: true,
		},
		{
			name: "successful decryption",
			decryptor: &mockDecryptor{
				decryptFunc: func(data []byte) ([]byte, error) {
					return encryptedData, nil
				},
			},
			req: &pb.ListMetricsRequest{
				CryptoMetrics: encryptedData,
			},
			handlerResp:   testMetrics,
			handlerErr:    nil,
			wantErr:       false,
			checkResponse: true,
		},
		{
			name: "decryption failure",
			decryptor: &mockDecryptor{
				decryptFunc: func(data []byte) ([]byte, error) {
					return nil, assert.AnError
				},
			},
			req: &pb.ListMetricsRequest{
				CryptoMetrics: encryptedData,
			},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "unmarshal failure",
			decryptor: &mockDecryptor{
				decryptFunc: func(data []byte) ([]byte, error) {
					return []byte("invalid-protobuf-data"), nil
				},
			},
			req: &pb.ListMetricsRequest{
				CryptoMetrics: encryptedData,
			},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем и восстанавливаем глобальный дешифратор
			originalDecryptor := crypto.GRPCDecryptor
			crypto.GRPCDecryptor = tt.decryptor
			defer func() { crypto.GRPCDecryptor = originalDecryptor }()

			handler := &mockHandler{
				resp: tt.handlerResp,
				err:  tt.handlerErr,
			}

			resp, err := DecryptUnaryInterceptor(
				context.Background(),
				tt.req,
				nil, // info
				handler.handle,
			)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrCode != codes.OK {
					assert.Equal(t, tt.wantErrCode, status.Code(err))
				}
				return
			}

			assert.NoError(t, err)
			if tt.checkResponse {
				assert.Equal(t, tt.handlerResp, resp)
			}
		})
	}
}
