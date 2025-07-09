package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
)

type envConfig struct {
	Config string `env:"CONFIG"`
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
		return err
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	return nil
}

// prepareConfig подготавливает параметры конфигурации для дальнейшей работы
func (cfg *AgentConfig) prepareConfig() error {
	reportInterval, err := time.ParseDuration(cfg.ReportInterval)
	if err != nil {
		return err
	}
	cfg.ReportInterval = strconv.FormatFloat(reportInterval.Seconds(), 'f', 0, 64)

	pollInterval, err := time.ParseDuration(cfg.PollInterval)
	if err != nil {
		return err
	}
	cfg.PollInterval = strconv.FormatFloat(pollInterval.Seconds(), 'f', 0, 64)

	return nil
}

// getConfigName получение названия файла конфигурации
func getConfigName() string {
	// configName имя файла конфигурации
	var configName = ""
	flag.StringVar(&configName, "c", "", "key for configuration")
	flag.StringVar(&configName, "config", "", "key for configuration")

	var ev envConfig
	_ = env.Parse(&ev)

	if ev.Config != "" {
		configName = ev.Config
	}

	return configName
}

// GetConfig установка значений конфигурации
func GetConfig() (*AgentConfig, error) {
	configName := getConfigName()

	var config AgentConfig
	var err error

	if configName == "" {
		return &config, err
	}

	if err = config.loadConfig(configName); err != nil {
		return nil, err
	}

	if err = config.prepareConfig(); err != nil {
		return nil, err
	}

	return &config, err
}
