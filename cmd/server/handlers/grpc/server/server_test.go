package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	grpcHandlers "github.com/ramil063/gometrics/cmd/server/handlers/grpc"
	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	"github.com/ramil063/gometrics/internal/security/crypto"
)

func TestGetGRPCServer1(t *testing.T) {
	tests := []struct {
		want    *grpc.Server
		name    string
		wantErr bool
	}{
		{
			name:    "test 1",
			want:    &grpc.Server{},
			wantErr: false,
		},
	}
	flags := &grpcHandlers.ServerConfigFlags{}
	storage := server.NewMemStorage()
	manager := crypto.NewCryptoManager()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGRPCServer(flags, storage, manager)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGRPCServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.IsType(t, tt.want, got)
		})
	}
}

func TestPrepareServerEnvironment(t *testing.T) {
	tests := []struct {
		want1 server.Storager
		want  *grpcHandlers.ServerConfigFlags
		want2 *crypto.Manager
		name  string
	}{
		{
			name: "test 1",
			want: &grpcHandlers.ServerConfigFlags{
				Address:         "localhost:3202",
				FileStoragePath: "internal/storage/files/grpc/metrics.json",
				StoreInterval:   300,
			},
			want1: server.NewFileStorage(),
			want2: crypto.NewCryptoManager(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := PrepareServerEnvironment()
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "PrepareServerEnvironment()")
			assert.Equalf(t, tt.want1, got1, "PrepareServerEnvironment()")
			assert.Equalf(t, tt.want2, got2, "PrepareServerEnvironment()")
		})
	}
}
