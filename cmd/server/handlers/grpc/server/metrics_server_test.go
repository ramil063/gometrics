package server

import (
	"context"
	"reflect"
	"testing"

	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	metrics "github.com/ramil063/gometrics/internal/grpc/proto"
)

func TestMetricsServer_UpdateMetrics(t *testing.T) {
	type fields struct {
		UnimplementedMetricsServer metrics.UnimplementedMetricsServer
		storage                    server.Storager
	}
	type args struct {
		ctx context.Context
		req *metrics.ListMetricsRequest
	}
	tests := []struct {
		want    *metrics.ListMetricsResponse
		args    args
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				UnimplementedMetricsServer: metrics.UnimplementedMetricsServer{},
				storage:                    server.GetStorage("", ""),
			},
			args: args{
				ctx: context.Background(),
				req: &metrics.ListMetricsRequest{
					Metrics: []*metrics.Metric{},
				},
			},
			want: &metrics.ListMetricsResponse{
				Metrics: []*metrics.Metric{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsServer{
				UnimplementedMetricsServer: tt.fields.UnimplementedMetricsServer,
				storage:                    tt.fields.storage,
			}
			got, err := s.UpdateMetrics(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMetrics() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMetricsServer(t *testing.T) {
	type args struct {
		storage server.Storager
	}
	s := server.GetStorage("", "")
	tests := []struct {
		want *MetricsServer
		args args
		name string
	}{
		{
			name: "TestNewMetricsServer",
			args: args{
				storage: s,
			},
			want: &MetricsServer{
				UnimplementedMetricsServer: metrics.UnimplementedMetricsServer{},
				storage:                    s,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetricsServer(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricsServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapMetricType(t *testing.T) {
	type args struct {
		mType string
	}
	tests := []struct {
		name string
		args args
		want metrics.Metric_MetricType
	}{
		{
			name: "Test gauge metric type",
			args: args{
				mType: "gauge",
			},
			want: metrics.Metric_gauge,
		},
		{
			name: "Test counter metric type",
			args: args{
				mType: "counter",
			},
			want: metrics.Metric_counter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapMetricType(tt.args.mType); got != tt.want {
				t.Errorf("mapMetricType() = %v, want %v", got, tt.want)
			}
		})
	}
}
