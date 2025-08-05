package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentConfig_Get(t *testing.T) {
	type conf struct {
		Address        string
		ReportInterval string
		PollInterval   string
		HashKey        string
		RateLimit      string
		CryptoKey      string
	}
	type wantConf struct {
		Address        string
		HashKey        string
		CryptoKey      string
		ReportInterval int
		PollInterval   int
		RateLimit      int
	}
	tests := []struct {
		name               string
		defaultStringValue string
		conf               conf
		wantConf           wantConf
		defaultIntValue    int
	}{
		{
			name: "test default value",
			conf: conf{
				Address:        "localhost:8080",
				ReportInterval: "1",
				PollInterval:   "1",
				HashKey:        "testhashkey",
				RateLimit:      "1",
				CryptoKey:      "testcryptokey",
			},
			wantConf: wantConf{
				Address:        "localhost:8080",
				ReportInterval: 1,
				PollInterval:   1,
				HashKey:        "testhashkey",
				RateLimit:      1,
				CryptoKey:      "testcryptokey",
			},
			defaultStringValue: "default",
			defaultIntValue:    100,
		},
		{
			name: "test default value",
			conf: conf{},
			wantConf: wantConf{
				Address:        "default",
				ReportInterval: 100,
				PollInterval:   100,
				HashKey:        "default",
				RateLimit:      100,
				CryptoKey:      "default",
			},
			defaultStringValue: "default",
			defaultIntValue:    100,
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
			assert.Equalf(t, tt.wantConf.Address, cfg.GetAddress(tt.defaultStringValue), "GetAddress(%v)", tt.defaultStringValue)
			assert.Equalf(t, tt.wantConf.CryptoKey, cfg.GetCryptoKey(tt.defaultStringValue), "GetCryptoKey(%v)", tt.defaultStringValue)
			assert.Equalf(t, tt.wantConf.HashKey, cfg.GetHashKey(tt.defaultStringValue), "GetHashKey(%v)", tt.defaultStringValue)
			assert.Equalf(t, tt.wantConf.PollInterval, cfg.GetPollInterval(tt.defaultIntValue), "GetPollInterval(%v)", tt.defaultIntValue)
			assert.Equalf(t, tt.wantConf.RateLimit, cfg.GetRateLimit(tt.defaultIntValue), "GetRateLimit(%v)", tt.defaultIntValue)
			assert.Equalf(t, tt.wantConf.ReportInterval, cfg.GetReportInterval(tt.defaultIntValue), "GetReportInterval(%v)", tt.defaultIntValue)
		})
	}
}
