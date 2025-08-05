package grpc

import (
	"context"
	"fmt"
	"log"
	"sync"

	grpcAgentConfig "github.com/ramil063/gometrics/cmd/agent/config/grpc"
	pb "github.com/ramil063/gometrics/internal/grpc/proto"
	"github.com/ramil063/gometrics/internal/hash"
	"github.com/ramil063/gometrics/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.MetricsClient
	once   sync.Once
}

func NewGRPCClient(serverAddr string) (*Client, error) {
	credentials := insecure.NewCredentials()
	conn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(credentials),
	)
	if err != nil {
		return nil, fmt.Errorf("NewGRPCClient error: %w", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewMetricsClient(conn),
	}, nil
}

func (c *Client) Close() error {
	c.once.Do(func() {
		c.conn.Close()
	})
	return nil
}

// SendMetrics отправляет массив метрик на сервер
func (c *Client) SendMetrics(r request, metrics []*pb.Metric, flags *SystemConfigFlags) error {
	// Создаем контекст с метаданными
	ctx := context.Background()

	// Создаем метаданные gRPC
	md := metadata.New(map[string]string{
		"x-real-ip": r.IP,
	})

	// Добавляем хеш, если указан ключ
	if flags.HashKey != "" {
		// Сериализуем метрики в байты для хеширования
		body, err := proto.Marshal(&pb.ListMetricsRequest{Metrics: metrics})
		if err != nil {
			return fmt.Errorf("failed to marshal metrics: %w", err)
		}

		hashSha256 := hash.CreateSha256(body, flags.HashKey)
		md.Set("hashsha256", hashSha256) // Добавляем хеш в метаданные
	}

	// Прикрепляем метаданные к контексту
	ctx = metadata.NewOutgoingContext(ctx, md)
	resp, err := c.client.UpdateMetrics(ctx, &pb.ListMetricsRequest{
		Metrics: metrics,
	})

	if err != nil {
		return fmt.Errorf("SendMetrics error: %w", err)
	}
	if resp.Error != "" {
		return fmt.Errorf("SendMetrics response error: %s", resp.Error)
	}
	return err
}

// StartClient запуск gRPC клиента
func StartClient(ctxGrSh context.Context, serversWg *sync.WaitGroup) {
	defer serversWg.Done()
	/**
	 * gRPC client
	 */
	configGRPC, err := grpcAgentConfig.GetConfig()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "config")
	}

	flagsGRPC, err := GetFlags(configGRPC)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "flags")
	}

	address := ""
	if flagsGRPC != nil && flagsGRPC.Address != "" {
		address = flagsGRPC.Address
	}
	grpcClient, err := NewGRPCClient(address)
	if err != nil {
		log.Println("NewGRPCClient error:", err)
	}
	defer grpcClient.Close()

	rGRPC := NewRequest()
	rGRPC.SendMetricsProcess(grpcClient, -1, ctxGrSh, flagsGRPC)
}
