package handlers

import (
	"os"
	"testing"

	serverConfig "github.com/ramil063/gometrics/cmd/server/config"
	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		args     []string
		expected EnvVars
	}{
		{
			name: "default values",
			args: []string{},
			expected: EnvVars{
				Address:         "localhost:8080",
				StoreInterval:   300,
				FileStoragePath: "internal/storage/files/metrics.json",
				Restore:         false,
				DatabaseDSN:     "",
				HashKey:         "",
				CryptoKey:       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original args and env
			oldArgs := os.Args
			oldEnv := map[string]string{}

			// Restore original values after test
			defer func() {
				os.Args = oldArgs
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
				for k, v := range oldEnv {
					os.Setenv(k, v)
				}
			}()

			// Set test args
			os.Args = append([]string{"cmd"}, tt.args...)

			// Set test env vars
			for k, v := range tt.envVars {
				if oldVal, exists := os.LookupEnv(k); exists {
					oldEnv[k] = oldVal
				}
				os.Setenv(k, v)
			}
			config := serverConfig.ServerConfig{}
			InitFlags(&config)

			assert.Equal(t, tt.expected.Address, MainURL)
			assert.Equal(t, tt.expected.StoreInterval, StoreInterval)
			assert.Equal(t, tt.expected.FileStoragePath, FileStoragePath)
			assert.Equal(t, tt.expected.Restore, Restore)
			assert.Equal(t, tt.expected.DatabaseDSN, DatabaseDSN)
			assert.Equal(t, tt.expected.HashKey, HashKey)
			assert.Equal(t, tt.expected.CryptoKey, CryptoKey)
		})
	}
}
