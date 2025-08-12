package grpc

import (
	"context"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/ramil063/gometrics/cmd/server/handlers/grpc/server"
	metrics "github.com/ramil063/gometrics/internal/grpc/proto"
	"github.com/ramil063/gometrics/internal/hash"
)

func TestClient_Close(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "test 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewGRPCClient(":3202")
			assert.NoError(t, err)
			err = c.Close()
			assert.NoError(t, err)
		})
	}
}

// Mock для gRPC клиента
type mockMetricsServiceClient struct {
	metrics.UnimplementedMetricsServer
	updateMetricsFunc func(context.Context, *metrics.ListMetricsRequest) (*metrics.ListMetricsResponse, error)
}

func (m *mockMetricsServiceClient) UpdateMetrics(ctx context.Context, req *metrics.ListMetricsRequest, opts ...grpc.CallOption) (*metrics.ListMetricsResponse, error) {
	return m.updateMetricsFunc(ctx, req)
}

func TestClient_SendMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Подготовка тестовых данных
	testMetrics := []*metrics.Metric{
		{Id: "cpu_usage", Value: 42.5},
	}
	testIP := "192.168.1.1"
	testHashKey := "secret"

	// Настройка мока
	mockClient := &mockMetricsServiceClient{
		updateMetricsFunc: func(ctx context.Context, req *metrics.ListMetricsRequest) (*metrics.ListMetricsResponse, error) {
			// Проверяем метаданные
			md, ok := metadata.FromOutgoingContext(ctx)
			assert.True(t, ok)
			assert.Equal(t, []string{testIP}, md.Get("x-real-ip"))

			// Проверяем хеш, если нужно
			if testHashKey != "" {
				body, _ := proto.Marshal(&metrics.ListMetricsRequest{Metrics: testMetrics})
				expectedHash := hash.CreateSha256(body, testHashKey)
				assert.Equal(t, []string{expectedHash}, md.Get("hashsha256"))
			}

			// Проверяем переданные метрики
			assert.Equal(t, testMetrics, req.Metrics)

			return &metrics.ListMetricsResponse{}, nil
		},
	}

	// Создаем тестовый клиент
	client := &Client{
		client: mockClient,
	}
	ctx, err := setHashByMetrics(request{IP: testIP}, testMetrics, &SystemConfigFlags{HashKey: testHashKey})
	assert.NoError(t, err)
	// Вызываем тестируемый метод
	err = client.SendMetrics(
		ctx,
		testMetrics,
		[]byte{},
	)

	// Проверяем результат
	assert.NoError(t, err)
}

func TestNewGRPCClient(t *testing.T) {
	// Запускаем тестовый gRPC сервер
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	metrics.RegisterMetricsServer(s, &server.MetricsServer{})

	go func() {
		_ = s.Serve(lis)
	}()
	defer s.Stop()

	// Вызываем тестируемую функцию с адресом тестового сервера
	client, err := NewGRPCClient(lis.Addr().String())

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Проверяем, что соединение работает
	_, err = client.client.UpdateMetrics(context.Background(), &metrics.ListMetricsRequest{})
	assert.NoError(t, err)
}

func TestStartClient(t *testing.T) {
	type args struct {
		ctxGrSh   context.Context
		serversWg *sync.WaitGroup
	}
	tests := []struct {
		args args
		name string
	}{
		{
			name: "test 1",
			args: args{
				ctxGrSh:   context.Background(),
				serversWg: &sync.WaitGroup{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(tt.args.ctxGrSh, time.Second*5)
			tt.args.serversWg.Add(1)
			StartClient(ctx, tt.args.serversWg)
			cancel()
		})
	}
}

func Test_setHashByMetrics(t *testing.T) {
	type args struct {
		flags   *SystemConfigFlags
		r       request
		metrics []*metrics.Metric
	}
	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "test 1",
			args: args{
				r:       request{},
				metrics: []*metrics.Metric{},
				flags:   &SystemConfigFlags{},
			},
			want: "*context.valueCtx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setHashByMetrics(tt.args.r, tt.args.metrics, tt.args.flags)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, reflect.ValueOf(got).Type().String(), "setHashByMetrics(%v, %v, %v)", tt.args.r, tt.args.metrics, tt.args.flags)
		})
	}
}
