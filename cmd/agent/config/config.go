package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/ramil063/gometrics/internal/constants"
	"github.com/ramil063/gometrics/internal/logger"
)

type envConfig struct {
	Config     string `env:"CONFIG"`
	GRPCConfig string `env:"GRPC_CONFIG"`
}

// AgentConfig структура для парсинга файла конфигурации
type AgentConfig struct {
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"`
	PollInterval   string `json:"poll_interval"`
	HashKey        string `json:"hash_key"`
	RateLimit      string `json:"rate_limit"`
	CryptoKey      string `json:"crypto_key"`
}

// loadConfig загружает конфигурацию из файла
func (cfg *AgentConfig) loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read the config file %s: %w", path, err)
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to unmurshal the config file %s: %w", path, err)
	}

	return nil
}

// prepareConfig подготавливает параметры конфигурации для дальнейшей работы
func (cfg *AgentConfig) prepareConfig() error {
	reportInterval, err := time.ParseDuration(cfg.ReportInterval)
	if err != nil {
		return fmt.Errorf("failed to parse ReportInterval: %w", err)
	}
	cfg.ReportInterval = strconv.FormatFloat(reportInterval.Seconds(), 'f', 0, 64)

	pollInterval, err := time.ParseDuration(cfg.PollInterval)
	if err != nil {
		return fmt.Errorf("failed to parse PollInterval: %w", err)
	}
	cfg.PollInterval = strconv.FormatFloat(pollInterval.Seconds(), 'f', 0, 64)

	return nil
}

// getConfigName получение названия файла конфигурации
func getConfigName(params ParamsProvider) string {
	// configName имя файла конфигурации
	var configName = ""
	flag.StringVar(&configName, params.GetConsoleKeyShort(), "", "key for configuration")
	flag.StringVar(&configName, params.GetConsoleKeyFull(), "", "key for configuration")

	var ev envConfig
	err := env.Parse(&ev)
	if err != nil {
		logger.WriteErrorLog("failed to parse config vars", "envConfig")
	}

	switch params.GetConfigType() {
	case constants.ConfigHTTPTypeAlias:
		configName = ev.Config
	case constants.ConfigGRPCTypeAlias:
		configName = ev.GRPCConfig
	}

	return configName
}

// GetConfig установка значений конфигурации
func GetConfig(params ParamsProvider) (*AgentConfig, error) {
	configName := getConfigName(params)

	var config AgentConfig
	var err error

	if configName == "" {
		return &config, nil
	}

	if err = config.loadConfig(configName); err != nil {
		return nil, err
	}

	if err = config.prepareConfig(); err != nil {
		return nil, err
	}

	return &config, err
}
