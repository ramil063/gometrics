package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/ramil063/gometrics/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestAgentConfig_loadConfig(t *testing.T) {

	file, _ := os.CreateTemp("", "config_test.json")
	_, _ = file.Write([]byte(`{
  "address": "testhost:8080",
  "report_interval": "1s",
  "poll_interval": "1s",
  "hash_key": "test",
  "rate_limit": "1",
  "crypto_key": "/test/test/test.pem"
}`))
	defer os.Remove(file.Name())

	type conf struct {
		Address        string
		ReportInterval string
		PollInterval   string
		HashKey        string
		RateLimit      string
		CryptoKey      string
	}
	tests := []struct {
		name string
		path string
		conf conf
	}{
		{
			name: "Load Config",
			path: file.Name(),
			conf: conf{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AgentConfig{
				Address:        tt.conf.Address,
				ReportInterval: tt.conf.ReportInterval,
				PollInterval:   tt.conf.PollInterval,
				HashKey:        tt.conf.HashKey,
				RateLimit:      tt.conf.RateLimit,
				CryptoKey:      tt.conf.CryptoKey,
			}
			err := cfg.loadConfig(tt.path)
			assert.NoError(t, err)
		})
	}
}

func TestAgentConfig_prepareConfig(t *testing.T) {
	type conf struct {
		Address        string
		ReportInterval string
		PollInterval   string
		HashKey        string
		RateLimit      string
		CryptoKey      string
	}
	tests := []struct {
		name string
		conf conf
	}{
		{
			name: "Prepare Config",
			conf: conf{
				Address:        "testhost:8080",
				ReportInterval: "1s",
				PollInterval:   "2h",
				HashKey:        "test",
				RateLimit:      "1",
				CryptoKey:      "/test/test/test.pem",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AgentConfig{
				Address:        tt.conf.Address,
				ReportInterval: tt.conf.ReportInterval,
				PollInterval:   tt.conf.PollInterval,
				HashKey:        tt.conf.HashKey,
				RateLimit:      tt.conf.RateLimit,
				CryptoKey:      tt.conf.CryptoKey,
			}
			err := cfg.prepareConfig()
			assert.NoError(t, err)
			assert.Equal(t, tt.conf.Address, cfg.Address)
			assert.Equal(t, "1", cfg.ReportInterval)
			assert.Equal(t, "7200", cfg.PollInterval)
			assert.Equal(t, tt.conf.HashKey, cfg.HashKey)
			assert.Equal(t, tt.conf.RateLimit, cfg.RateLimit)
			assert.Equal(t, tt.conf.CryptoKey, cfg.CryptoKey)
		})
	}
}

func TestGetConfig(t *testing.T) {
	file, _ := os.CreateTemp("", "config_test.json")
	_, _ = file.Write([]byte(`{
  "address": "testhost:8080",
  "report_interval": "1s",
  "poll_interval": "1s",
  "hash_key": "test",
  "rate_limit": "1",
  "crypto_key": "/test/test/test.pem"
}`))
	defer os.Remove(file.Name())

	_ = os.Setenv("CONFIG", file.Name())
	defer os.Unsetenv("CONFIG")

	tests := []struct {
		want *AgentConfig
		name string
	}{
		{
			name: "Get Config",
			want: &AgentConfig{
				Address:        "testhost:8080",
				ReportInterval: "1",
				PollInterval:   "1",
				HashKey:        "test",
				RateLimit:      "1",
				CryptoKey:      "/test/test/test.pem",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := NewConfigParams(
				constants.ConfigHTTPConsoleShortKey,
				constants.ConfigHTTPConsoleFullKey,
				constants.ConfigHTTPTypeAlias)
			got, err := GetConfig(params)
			assert.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
