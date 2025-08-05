package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentConfig_Get(t *testing.T) {
	restoreFalse := false
	restoreTrue := true
	type conf struct {
		Restore         *bool
		Address         string
		FileStoragePath string
		DatabaseDSN     string
		HashKey         string
		CryptoKey       string
		StoreInterval   string
		TrustedSubnet   string
	}
	type wantConf struct {
		Restore         *bool
		Address         string
		FileStoragePath string
		DatabaseDSN     string
		HashKey         string
		CryptoKey       string
		TrustedSubnet   string
		StoreInterval   int
	}
	tests := []struct {
		name               string
		defaultStringValue string
		conf               conf
		wantConf           wantConf
		defaultBoolValue   bool
		defaultIntValue    int
	}{
		{
			name: "test default value",
			conf: conf{
				Address:         "localhost:8080",
				FileStoragePath: "testfilepath",
				DatabaseDSN:     "testdatabase",
				HashKey:         "testhashkey",
				CryptoKey:       "testcryptokey",
				TrustedSubnet:   "testtrustedsubnet",
				StoreInterval:   "1",
				Restore:         &restoreFalse,
			},
			wantConf: wantConf{
				Address:         "localhost:8080",
				FileStoragePath: "testfilepath",
				DatabaseDSN:     "testdatabase",
				HashKey:         "testhashkey",
				CryptoKey:       "testcryptokey",
				TrustedSubnet:   "testtrustedsubnet",
				StoreInterval:   1,
				Restore:         &restoreFalse,
			},
			defaultStringValue: "default",
			defaultIntValue:    100,
			defaultBoolValue:   true,
		},
		{
			name: "test default value",
			conf: conf{},
			wantConf: wantConf{
				Address:         "default",
				FileStoragePath: "default",
				DatabaseDSN:     "default",
				HashKey:         "default",
				CryptoKey:       "default",
				TrustedSubnet:   "default",
				StoreInterval:   100,
				Restore:         &restoreTrue,
			},
			defaultStringValue: "default",
			defaultIntValue:    100,
			defaultBoolValue:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &ServerConfig{
				Address:         tt.conf.Address,
				FileStoragePath: tt.conf.FileStoragePath,
				DatabaseDSN:     tt.conf.DatabaseDSN,
				HashKey:         tt.conf.HashKey,
				CryptoKey:       tt.conf.CryptoKey,
				TrustedSubnet:   tt.conf.TrustedSubnet,
				StoreInterval:   tt.conf.StoreInterval,
				Restore:         tt.conf.Restore,
			}
			assert.Equalf(t, tt.wantConf.Address, cfg.GetAddress(tt.defaultStringValue), "GetAddress(%v)", tt.defaultStringValue)
			assert.Equalf(t, tt.wantConf.FileStoragePath, cfg.GetFileStoragePath(tt.defaultStringValue), "GetCryptoKey(%v)", tt.defaultStringValue)
			assert.Equalf(t, tt.wantConf.DatabaseDSN, cfg.GetDatabaseDSN(tt.defaultStringValue), "GetCryptoKey(%v)", tt.defaultStringValue)
			assert.Equalf(t, tt.wantConf.HashKey, cfg.GetHashKey(tt.defaultStringValue), "GetHashKey(%v)", tt.defaultStringValue)
			assert.Equalf(t, tt.wantConf.CryptoKey, cfg.GetCryptoKey(tt.defaultStringValue), "GetCryptoKey(%v)", tt.defaultStringValue)
			assert.Equalf(t, tt.wantConf.StoreInterval, cfg.GetStoreInterval(tt.defaultIntValue), "GetHashKey(%v)", tt.defaultIntValue)
			assert.Equalf(t, *(tt.wantConf.Restore), cfg.GetRestore(tt.defaultBoolValue), "GetHashKey(%v)", tt.defaultBoolValue)
			assert.Equalf(t, tt.wantConf.TrustedSubnet, cfg.GetTrustedSubnet(tt.defaultStringValue), "GetTrustedSubnet(%v)", tt.defaultStringValue)
		})
	}
}
