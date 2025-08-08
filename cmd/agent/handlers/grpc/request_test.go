package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/cmd/server/handlers/grpc/server"
	pb "github.com/ramil063/gometrics/internal/grpc/proto"
	"github.com/ramil063/gometrics/internal/models"
)

func TestConvertToProto(t *testing.T) {
	type args struct {
		metrics []models.Metrics
	}
	gaugeVal := float64(1.1)
	tests := []struct {
		name string
		args args
		want []*pb.Metric
	}{
		{
			name: "test 1",
			args: args{
				metrics: []models.Metrics{
					{ID: "met1", MType: "gauge", Value: &gaugeVal},
				},
			},
			want: []*pb.Metric{
				{Id: "met1", Type: pb.Metric_gauge, Value: 1.1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ConvertToProto(tt.args.metrics), "ConvertToProto(%v)", tt.args.metrics)
		})
	}
}

func TestNewRequest(t *testing.T) {
	tests := []struct {
		want Requester
		name string
	}{
		{
			name: "test 1",
			want: request{
				IP: "127.0.1.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewRequest(), "NewRequest()")
		})
	}
}

func TestSendMetricsByGRPC(t *testing.T) {
	type args struct {
		c       Clienter
		flags   *SystemConfigFlags
		monitor *storage.Monitor
		r       request
	}
	client, err := NewGRPCClient(":3202")
	assert.NoError(t, err)
	tests := []struct {
		args args
		name string
	}{
		{
			name: "test 1",
			args: args{
				r: request{
					IP: "127.0.1.1",
				},
				c:       client,
				monitor: &storage.Monitor{},
				flags:   &SystemConfigFlags{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendMetricsByGRPC(tt.args.r, tt.args.c, tt.args.monitor, tt.args.flags)
		})
	}
}

func Test_request_SendMetricsProcess(t *testing.T) {
	type fields struct {
		IP string
	}
	type args struct {
		c        Clienter
		ctxGrSh  context.Context
		flags    *SystemConfigFlags
		maxCount int
	}
	client, err := NewGRPCClient(":3202")
	assert.NoError(t, err)
	tests := []struct {
		args   args
		fields fields
		name   string
	}{
		{
			name: "test 1",
			fields: fields{
				IP: "127.0.0.1",
			},
			args: args{
				ctxGrSh:  context.Background(),
				maxCount: 1,
				flags: &SystemConfigFlags{
					Address:        "localhost:3202",
					PollInterval:   2,
					ReportInterval: 10,
					RateLimit:      1,
				},
				c: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request{
				IP: tt.fields.IP,
			}
			r.SendMetricsProcess(tt.args.c, tt.args.maxCount, tt.args.ctxGrSh, tt.args.flags)
		})
	}
}

func Test_request_getOutboundIP(t *testing.T) {
	type fields struct {
		IP string
	}
	tests := []struct {
		fields fields
		name   string
		want   string
	}{
		{
			name: "test 1",
			fields: fields{
				IP: "127.0.1.1",
			},
			want: "127.0.1.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request{
				IP: tt.fields.IP,
			}
			got, err := r.getOutboundIP()
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "getOutboundIP()")
		})
	}
}

func Test_retryToSendMetrics(t *testing.T) {
	// Запускаем тестовый gRPC сервер
	lis, err := net.Listen("tcp", ":3202")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMetricsServer(s, &server.MetricsServer{})

	go func() {
		_ = s.Serve(lis)
	}()
	defer s.Stop()

	type args struct {
		c       Clienter
		metrics []*pb.Metric
		flags   *SystemConfigFlags
		r       request
		tries   []int
	}
	client, err := NewGRPCClient(":3202")
	assert.NoError(t, err)
	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
		args    args
	}{
		{
			name: "test 1",
			args: args{
				r: request{
					IP: "127.0.0.1",
				},
				c:       client,
				metrics: []*pb.Metric{},
				tries:   []int{0},
				flags:   &SystemConfigFlags{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = retryToSendMetrics(tt.args.r, tt.args.c, tt.args.metrics, tt.args.tries, tt.args.flags)
			assert.NoError(t, err)
		})
	}
}
