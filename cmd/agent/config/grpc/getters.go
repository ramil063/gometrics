package grpc

import "strconv"

// GetAddress получение параметра Address
func (cfg *AgentConfig) GetAddress(defaultValue string) string {
	if cfg.Address != "" {
		return cfg.Address
	}
	return defaultValue
}

// GetReportInterval получение параметра ReportInterval
func (cfg *AgentConfig) GetReportInterval(defaultValue int) int {
	if cfg.ReportInterval != "0" {
		if val, err := strconv.Atoi(cfg.ReportInterval); err == nil {
			return val
		}
	}
	return defaultValue
}

// GetPollInterval получение параметра PollInterval
func (cfg *AgentConfig) GetPollInterval(defaultValue int) int {
	if cfg.PollInterval != "0" {
		if val, err := strconv.Atoi(cfg.PollInterval); err == nil {
			return val
		}
	}
	return defaultValue
}

// GetHashKey получение параметра HashKey
func (cfg *AgentConfig) GetHashKey(defaultValue string) string {
	if cfg.HashKey != "" {
		return cfg.HashKey
	}
	return defaultValue
}

// GetRateLimit получение параметра RateLimit
func (cfg *AgentConfig) GetRateLimit(defaultValue int) int {
	if cfg.RateLimit != "0" {
		if val, err := strconv.Atoi(cfg.RateLimit); err == nil {
			return val
		}
	}
	return defaultValue
}

// GetCryptoKey получение параметра CryptoKey
func (cfg *AgentConfig) GetCryptoKey(defaultValue string) string {
	if cfg.CryptoKey != "" {
		return cfg.CryptoKey
	}
	return defaultValue
}
