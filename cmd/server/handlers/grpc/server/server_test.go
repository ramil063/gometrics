package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGRPCServer()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGRPCServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.IsType(t, tt.want, got)
		})
	}
}
