package handlers

import (
	"flag"
	"os"
	"testing"

	"github.com/ramil063/gometrics/cmd/agent/config"
	"github.com/stretchr/testify/assert"
)

func TestGetFlags(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		args     []string
		expected struct {
			address        string
			hashKey        string
			cryptoKey      string
			reportInterval int
			pollInterval   int
			rateLimit      int
		}
	}{
		{
			name:    "default values when no flags or env vars",
			envVars: map[string]string{},
			args:    []string{},
			expected: struct {
				address        string
				hashKey        string
				cryptoKey      string
				reportInterval int
				pollInterval   int
				rateLimit      int
			}{
				address:        "localhost:8080",
				hashKey:        "",
				reportInterval: 10,
				pollInterval:   2,
				rateLimit:      1,
				cryptoKey:      "",
			},
		},
		{
			name:    "command line flags override defaults",
			envVars: map[string]string{},
			args: []string{
				"-a", "localhost:9090",
				"-k", "secret",
				"-r", "20",
				"-p", "5",
				"-l", "3",
				"-crypto-key", "secret",
			},
			expected: struct {
				address        string
				hashKey        string
				cryptoKey      string
				reportInterval int
				pollInterval   int
				rateLimit      int
			}{
				address:        "localhost:9090",
				hashKey:        "secret",
				reportInterval: 20,
				pollInterval:   5,
				rateLimit:      3,
				cryptoKey:      "secret",
			},
		},
		{
			name: "env vars override defaults",
			envVars: map[string]string{
				"ADDRESS":         "localhost:7070",
				"KEY":             "envkey",
				"REPORT_INTERVAL": "15",
				"POLL_INTERVAL":   "3",
				"RATE_LIMIT":      "2",
				"CRYPTO_KEY":      "envkey",
			},
			args: []string{},
			expected: struct {
				address        string
				hashKey        string
				cryptoKey      string
				reportInterval int
				pollInterval   int
				rateLimit      int
			}{
				address:        "localhost:7070",
				hashKey:        "envkey",
				reportInterval: 15,
				pollInterval:   3,
				rateLimit:      2,
				cryptoKey:      "envkey",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags before each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Set command line args
			os.Args = append([]string{"cmd"}, tt.args...)

			configMock := config.AgentConfig{}
			flags, err := GetFlags(&configMock)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.address, flags.Address)
			assert.Equal(t, tt.expected.hashKey, flags.HashKey)
			assert.Equal(t, tt.expected.reportInterval, flags.ReportInterval)
			assert.Equal(t, tt.expected.pollInterval, flags.PollInterval)
			assert.Equal(t, tt.expected.rateLimit, flags.RateLimit)
		})
	}
}

func Test_applyFlags(t *testing.T) {
	type args struct {
		flags          *SystemConfigFlags
		address        string
		hashKey        string
		cryptoKey      string
		reportInterval int
		pollInterval   int
		rateLimit      int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test 1",
			args: args{
				flags: &SystemConfigFlags{
					Address:        "defaultLocalhost:8080",
					HashKey:        "defaultHashKey",
					CryptoKey:      "defaultCryptoKey",
					ReportInterval: 1,
					PollInterval:   2,
					RateLimit:      3,
				},
				address:        "notDefaultLocalhost:8080",
				hashKey:        "notDefault",
				cryptoKey:      "notDefault",
				reportInterval: 10,
				pollInterval:   20,
				rateLimit:      30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyFlags(tt.args.flags, tt.args.address, tt.args.reportInterval, tt.args.pollInterval, tt.args.hashKey, tt.args.rateLimit, tt.args.cryptoKey)
			assert.Equal(t, tt.args.flags.Address, tt.args.address)
			assert.Equal(t, tt.args.hashKey, tt.args.hashKey)
			assert.Equal(t, tt.args.reportInterval, tt.args.reportInterval)
			assert.Equal(t, tt.args.pollInterval, tt.args.pollInterval)
			assert.Equal(t, tt.args.rateLimit, tt.args.rateLimit)
		})
	}
}

func Test_applyEnvVars(t *testing.T) {
	type args struct {
		flags   *SystemConfigFlags
		envVars SystemConfigFlags
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test 1",
			args: args{
				flags: &SystemConfigFlags{
					Address:        "defaultLocalhost:8080",
					HashKey:        "defaultHashKey",
					CryptoKey:      "defaultCryptoKey",
					ReportInterval: 1,
					PollInterval:   2,
					RateLimit:      3,
				},
				envVars: SystemConfigFlags{
					Address:        "notDefaultLocalhost:8080",
					HashKey:        "notDefault",
					CryptoKey:      "notDefault",
					ReportInterval: 10,
					PollInterval:   20,
					RateLimit:      30,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyEnvVars(tt.args.flags, tt.args.envVars)
			assert.Equal(t, tt.args.flags.Address, tt.args.envVars.Address)
			assert.Equal(t, tt.args.flags.HashKey, tt.args.envVars.HashKey)
			assert.Equal(t, tt.args.flags.CryptoKey, tt.args.envVars.CryptoKey)
			assert.Equal(t, tt.args.envVars.ReportInterval, tt.args.envVars.ReportInterval)
			assert.Equal(t, tt.args.envVars.PollInterval, tt.args.envVars.PollInterval)
			assert.Equal(t, tt.args.envVars.RateLimit, tt.args.envVars.RateLimit)
		})
	}
}
