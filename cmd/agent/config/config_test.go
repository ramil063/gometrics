package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentConfig_loadConfig(t *testing.T) {

	file, _ := os.OpenFile("config_test.json", os.O_WRONLY|os.O_CREATE, 0766)
	_, _ = file.Write([]byte(`{
  "address": "testhost:8080",
  "report_interval": "1s",
  "poll_interval": "1s",
  "hash_key": "test",
  "rate_limit": "1",
  "crypto_key": "/test/test/test.pem"
}`))

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
			path: "config_test.json",
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
	_ = os.Remove("config_test.json")
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
	file, _ := os.OpenFile("config_test.json", os.O_WRONLY|os.O_CREATE, 0766)
	_, _ = file.Write([]byte(`{
  "address": "testhost:8080",
  "report_interval": "1s",
  "poll_interval": "1s",
  "hash_key": "test",
  "rate_limit": "1",
  "crypto_key": "/test/test/test.pem"
}`))
	_ = os.Setenv("CONFIG", "config_test.json")
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
			got, err := GetConfig()
			assert.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
	_ = os.Remove("config_test.json")
}
