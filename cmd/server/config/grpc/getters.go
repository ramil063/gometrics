package grpc

import "strconv"

// GetAddress получение параметра Address
func (cfg *ServerConfig) GetAddress(defaultValue string) string {
	if cfg.Address != "" {
		return cfg.Address
	}
	return defaultValue
}

// GetFileStoragePath получение параметра FileStoragePath
func (cfg *ServerConfig) GetFileStoragePath(defaultValue string) string {
	if cfg.FileStoragePath != "" {
		return cfg.FileStoragePath
	}
	return defaultValue
}

// GetDatabaseDSN получение параметра DatabaseDSN
func (cfg *ServerConfig) GetDatabaseDSN(defaultValue string) string {
	if cfg.DatabaseDSN != "" {
		return cfg.DatabaseDSN
	}
	return defaultValue
}

// GetHashKey получение параметра HashKey
func (cfg *ServerConfig) GetHashKey(defaultValue string) string {
	if cfg.HashKey != "" {
		return cfg.HashKey
	}
	return defaultValue
}

// GetCryptoKey получение параметра CryptoKey
func (cfg *ServerConfig) GetCryptoKey(defaultValue string) string {
	if cfg.CryptoKey != "" {
		return cfg.CryptoKey
	}
	return defaultValue
}

// GetStoreInterval получение параметра StoreInterval
func (cfg *ServerConfig) GetStoreInterval(defaultValue int) int {
	if cfg.StoreInterval != "0" {
		if val, err := strconv.Atoi(cfg.StoreInterval); err == nil {
			return val
		}
	}
	return defaultValue
}

// GetRestore получение параметра Restore
func (cfg *ServerConfig) GetRestore(defaultValue bool) bool {
	if cfg.Restore != nil {
		return *cfg.Restore
	}
	return defaultValue
}

// GetTrustedSubnet получение параметра TrustedSubnet
func (cfg *ServerConfig) GetTrustedSubnet(defaultValue string) string {
	if cfg.TrustedSubnet != "" {
		return cfg.TrustedSubnet
	}
	return defaultValue
}
