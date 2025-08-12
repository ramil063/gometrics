package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	pb "github.com/ramil063/gometrics/internal/grpc/proto"
	"github.com/ramil063/gometrics/internal/models"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer

	storage server.Storager
}

// NewMetricsServer получение нового сервера для обновления метрик
func NewMetricsServer(storage server.Storager) *MetricsServer {
	return &MetricsServer{
		storage: storage,
	}
}

// UpdateMetrics основная функция обновления метрик
func (s *MetricsServer) UpdateMetrics(ctx context.Context, req *pb.ListMetricsRequest) (*pb.ListMetricsResponse, error) {
	// 1. Конвертируем protobuf -> models.Metrics
	metrics := make([]models.Metrics, 0, len(req.GetMetrics()))
	for _, pbMetric := range req.GetMetrics() {
		m := models.Metrics{
			ID:    pbMetric.GetId(),
			MType: pbMetric.GetType().String(),
		}

		switch pbMetric.GetType() {
		case pb.Metric_gauge:
			val := pbMetric.GetValue()
			m.Value = &val
		case pb.Metric_counter:
			delta := pbMetric.GetDelta()
			m.Delta = &delta
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unknown metric type: %v", pbMetric.GetType())
		}
		metrics = append(metrics, m)
	}

	// 2. Вызываем логику обработки
	result, err := server.UpdateMetrics(s.storage, metrics)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update metrics failed: %v", err)
	}

	// 3. Конвертируем результат обратно в protobuf
	pbResults := make([]*pb.Metric, 0, len(result))
	for _, m := range result {
		pbMetric := &pb.Metric{
			Id:   m.ID,
			Type: mapMetricType(m.MType),
		}

		switch m.MType {
		case "gauge":
			if m.Value != nil {
				pbMetric.Value = *m.Value
			}
		case "counter":
			if m.Delta != nil {
				pbMetric.Delta = *m.Delta
			}
		}
		pbResults = append(pbResults, pbMetric)
	}

	return &pb.ListMetricsResponse{
		Metrics: pbResults,
		Error:   "",
	}, nil
}

// Вспомогательная функция для конвертации типа
func mapMetricType(mType string) pb.Metric_MetricType {
	switch mType {
	case "gauge":
		return pb.Metric_gauge
	case "counter":
		return pb.Metric_counter
	default:
		return pb.Metric_gauge
	}
}
