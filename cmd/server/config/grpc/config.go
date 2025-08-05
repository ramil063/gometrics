package grpc

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/ramil063/gometrics/internal/logger"
)

type envConfig struct {
	Config string `env:"GRPC_CONFIG"`
}

// ServerConfig структура для парсинга файла конфигурации
type ServerConfig struct {
	Restore         *bool  `json:"restore"`
	Address         string `json:"address"`
	FileStoragePath string `json:"store_file"`
	DatabaseDSN     string `json:"database_dsn"`
	HashKey         string `json:"hash_key"`
	CryptoKey       string `json:"crypto_key"`
	StoreInterval   string `json:"store_interval"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

// loadConfig загружает конфигурацию из файла
func (cfg *ServerConfig) loadConfig(path string) error {
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
func (cfg *ServerConfig) prepareConfig() error {
	storeInterval, err := time.ParseDuration(cfg.StoreInterval)
	if err != nil {
		return fmt.Errorf("failed to parse ReportInterval: %w", err)
	}
	cfg.StoreInterval = strconv.FormatFloat(storeInterval.Seconds(), 'f', 0, 64)

	return nil
}

// getConfigName получение названия файла конфигурации
func getConfigName() string {
	// configName имя файла конфигурации
	var configName = ""
	flag.StringVar(&configName, "grpc-c", "", "key for configuration")
	flag.StringVar(&configName, "grpc-config", "", "key for configuration")

	var ev envConfig
	err := env.Parse(&ev)
	if err != nil {
		logger.WriteErrorLog("failed to parse config vars", "envConfig")
	}

	if ev.Config != "" {
		configName = ev.Config
	}

	return configName
}

// GetConfig установка значений конфигурации
func GetConfig() (*ServerConfig, error) {
	configName := getConfigName()

	var config ServerConfig
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
