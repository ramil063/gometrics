package handlers

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		args     []string
		expected struct {
			mainURL        string
			hashKey        string
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
				mainURL        string
				hashKey        string
				reportInterval int
				pollInterval   int
				rateLimit      int
			}{
				mainURL:        "localhost:8080",
				hashKey:        "",
				reportInterval: 10,
				pollInterval:   2,
				rateLimit:      1,
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
			},
			expected: struct {
				mainURL        string
				hashKey        string
				reportInterval int
				pollInterval   int
				rateLimit      int
			}{
				mainURL:        "localhost:9090",
				hashKey:        "secret",
				reportInterval: 20,
				pollInterval:   5,
				rateLimit:      3,
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
			},
			args: []string{},
			expected: struct {
				mainURL        string
				hashKey        string
				reportInterval int
				pollInterval   int
				rateLimit      int
			}{
				mainURL:        "localhost:7070",
				hashKey:        "envkey",
				reportInterval: 15,
				pollInterval:   3,
				rateLimit:      2,
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

			// Reset global variables
			MainURL = "localhost:8080"
			HashKey = ""
			ReportInterval = 10
			PollInterval = 2
			RateLimit = 1

			ParseFlags()

			assert.Equal(t, tt.expected.mainURL, MainURL)
			assert.Equal(t, tt.expected.hashKey, HashKey)
			assert.Equal(t, tt.expected.reportInterval, ReportInterval)
			assert.Equal(t, tt.expected.pollInterval, PollInterval)
			assert.Equal(t, tt.expected.rateLimit, RateLimit)
		})
	}
}
