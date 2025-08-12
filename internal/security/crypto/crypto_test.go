package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRSADecryptor(t *testing.T) {
	privateFile, _ := os.OpenFile("priv_test.pem", os.O_WRONLY|os.O_CREATE, 0766)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	bytes, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bytes,
	}
	_ = pem.Encode(privateFile, privateBlock)

	tests := []struct {
		name           string
		privateKeyPath string
		want           string
	}{
		{
			name:           "test1",
			privateKeyPath: "priv_test.pem",
			want:           "*rsa.RsaDecryptor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRSADecryptor(tt.privateKeyPath)
			assert.Equal(t, tt.want, reflect.ValueOf(got).Type().String())
			assert.NotNil(t, got)
			assert.NoError(t, err)
		})
	}
	_ = os.Remove("priv_test.pem")
}

func TestNewRSAEncryptor(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	publicKey := &privateKey.PublicKey
	publicBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	publicFile, _ := os.Create("pub_test.pem")
	defer publicFile.Close()

	publicBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	}

	_ = pem.Encode(publicFile, publicBlock)

	tests := []struct {
		name          string
		publicKeyPath string
		want          string
	}{
		{
			name:          "test1",
			publicKeyPath: "pub_test.pem",
			want:          "*rsa.RsaEncryptor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRSAEncryptor(tt.publicKeyPath)
			assert.Equal(t, tt.want, reflect.ValueOf(got).Type().String())
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
	_ = os.Remove("pub_test.pem")
}

func TestManager_GetGRPCDecryptor(t *testing.T) {
	type fields struct {
		defaultEncryptor Encryptor
		defaultDecryptor Decryptor
		grpcEncryptor    Encryptor
		grpcDecryptor    Decryptor
		mx               sync.RWMutex
	}
	var decryptor Decryptor
	tests := []struct {
		name   string
		fields fields
		want   Decryptor
	}{
		{
			name: "test1",
			fields: fields{
				defaultDecryptor: nil,
				defaultEncryptor: nil,
				grpcEncryptor:    nil,
				grpcDecryptor:    decryptor,
			},
			want: decryptor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cm := &Manager{
				defaultEncryptor: tt.fields.defaultEncryptor,
				defaultDecryptor: tt.fields.defaultDecryptor,
				grpcEncryptor:    tt.fields.grpcEncryptor,
				grpcDecryptor:    tt.fields.grpcDecryptor,
				mx:               tt.fields.mx,
			}
			assert.Equalf(t, tt.want, cm.GetGRPCDecryptor(), "GetGRPCDecryptor()")
		})
	}
}

func TestManager_GetDefaultDecryptor(t *testing.T) {
	type fields struct {
		defaultEncryptor Encryptor
		defaultDecryptor Decryptor
		grpcEncryptor    Encryptor
		grpcDecryptor    Decryptor
		mx               sync.RWMutex
	}
	var decryptor Decryptor
	tests := []struct {
		name   string
		fields fields
		want   Decryptor
	}{
		{
			name: "test1",
			fields: fields{
				defaultEncryptor: nil,
				defaultDecryptor: decryptor,
				grpcEncryptor:    nil,
				grpcDecryptor:    nil,
			},
			want: decryptor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				defaultEncryptor: tt.fields.defaultEncryptor,
				defaultDecryptor: tt.fields.defaultDecryptor,
				grpcEncryptor:    tt.fields.grpcEncryptor,
				grpcDecryptor:    tt.fields.grpcDecryptor,
				mx:               tt.fields.mx,
			}
			assert.Equalf(t, tt.want, cm.GetDefaultDecryptor(), "GetDefaultDecryptor()")
		})
	}
}

func TestManager_GetDefaultEncryptor(t *testing.T) {
	type fields struct {
		defaultEncryptor Encryptor
		defaultDecryptor Decryptor
		grpcEncryptor    Encryptor
		grpcDecryptor    Decryptor
		mx               sync.RWMutex
	}
	var encryptor Encryptor
	tests := []struct {
		name   string
		fields fields
		want   Encryptor
	}{
		{
			name: "test1",
			fields: fields{
				defaultEncryptor: encryptor,
				defaultDecryptor: nil,
				grpcEncryptor:    nil,
				grpcDecryptor:    nil,
			},
			want: encryptor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				defaultEncryptor: tt.fields.defaultEncryptor,
				defaultDecryptor: tt.fields.defaultDecryptor,
				grpcEncryptor:    tt.fields.grpcEncryptor,
				grpcDecryptor:    tt.fields.grpcDecryptor,
				mx:               tt.fields.mx,
			}
			assert.Equalf(t, tt.want, cm.GetDefaultEncryptor(), "GetDefaultEncryptor()")
		})
	}
}

func TestManager_GetGRPCEncryptor(t *testing.T) {
	type fields struct {
		defaultEncryptor Encryptor
		defaultDecryptor Decryptor
		grpcEncryptor    Encryptor
		grpcDecryptor    Decryptor
		mx               sync.RWMutex
	}
	var encryptor Encryptor
	tests := []struct {
		name   string
		fields fields
		want   Encryptor
	}{
		{
			name: "test1",
			fields: fields{
				defaultEncryptor: nil,
				defaultDecryptor: nil,
				grpcEncryptor:    encryptor,
				grpcDecryptor:    nil,
			},
			want: encryptor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				defaultEncryptor: tt.fields.defaultEncryptor,
				defaultDecryptor: tt.fields.defaultDecryptor,
				grpcEncryptor:    tt.fields.grpcEncryptor,
				grpcDecryptor:    tt.fields.grpcDecryptor,
				mx:               tt.fields.mx,
			}
			assert.Equalf(t, tt.want, cm.GetGRPCEncryptor(), "GetGRPCEncryptor()")
		})
	}
}

func TestManager_SetDefaultDecryptor(t *testing.T) {
	type fields struct {
		defaultEncryptor Encryptor
		defaultDecryptor Decryptor
		grpcEncryptor    Encryptor
		grpcDecryptor    Decryptor
		mx               sync.RWMutex
	}
	type args struct {
		decr Decryptor
	}
	var decryptor Decryptor
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "test1",
			fields: fields{},
			args: args{
				decr: decryptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				defaultEncryptor: tt.fields.defaultEncryptor,
				defaultDecryptor: tt.fields.defaultDecryptor,
				grpcEncryptor:    tt.fields.grpcEncryptor,
				grpcDecryptor:    tt.fields.grpcDecryptor,
				mx:               tt.fields.mx,
			}
			cm.SetDefaultDecryptor(tt.args.decr)
			assert.Equalf(t, tt.args.decr, cm.GetDefaultDecryptor(), "GetDefaultDecryptor()")
		})
	}
}

func TestManager_SetDefaultEncryptor(t *testing.T) {
	type fields struct {
		defaultEncryptor Encryptor
		defaultDecryptor Decryptor
		grpcEncryptor    Encryptor
		grpcDecryptor    Decryptor
		mx               sync.RWMutex
	}
	type args struct {
		enc Encryptor
	}
	var encryptor Encryptor
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "test1",
			fields: fields{},
			args: args{
				enc: encryptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				defaultEncryptor: tt.fields.defaultEncryptor,
				defaultDecryptor: tt.fields.defaultDecryptor,
				grpcEncryptor:    tt.fields.grpcEncryptor,
				grpcDecryptor:    tt.fields.grpcDecryptor,
				mx:               tt.fields.mx,
			}
			cm.SetDefaultEncryptor(tt.args.enc)
			assert.Equal(t, tt.args.enc, cm.GetDefaultEncryptor(), "GetDefaultEncryptor()")
		})
	}
}

func TestManager_SetGRPCDecryptor(t *testing.T) {
	type fields struct {
		defaultEncryptor Encryptor
		defaultDecryptor Decryptor
		grpcEncryptor    Encryptor
		grpcDecryptor    Decryptor
		mx               sync.RWMutex
	}
	type args struct {
		decr Decryptor
	}
	var decryptor Decryptor
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "test1",
			fields: fields{},
			args: args{
				decr: decryptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				defaultEncryptor: tt.fields.defaultEncryptor,
				defaultDecryptor: tt.fields.defaultDecryptor,
				grpcEncryptor:    tt.fields.grpcEncryptor,
				grpcDecryptor:    tt.fields.grpcDecryptor,
				mx:               tt.fields.mx,
			}
			cm.SetGRPCDecryptor(tt.args.decr)
			assert.Equal(t, tt.args.decr, cm.GetGRPCDecryptor(), "GetGRPCDecryptor()")
		})
	}
}

func TestManager_SetGRPCEncryptor(t *testing.T) {
	type fields struct {
		defaultEncryptor Encryptor
		defaultDecryptor Decryptor
		grpcEncryptor    Encryptor
		grpcDecryptor    Decryptor
		mx               sync.RWMutex
	}
	type args struct {
		enc Encryptor
	}
	var encryptor Encryptor
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "SetGRPCEncryptor",
			fields: fields{},
			args: args{
				enc: encryptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				defaultEncryptor: tt.fields.defaultEncryptor,
				defaultDecryptor: tt.fields.defaultDecryptor,
				grpcEncryptor:    tt.fields.grpcEncryptor,
				grpcDecryptor:    tt.fields.grpcDecryptor,
				mx:               tt.fields.mx,
			}
			cm.SetGRPCEncryptor(tt.args.enc)
			assert.Equalf(t, tt.args.enc, cm.GetGRPCEncryptor(), "SetGRPCEncryptor()")
		})
	}
}

func TestNewCryptoManager(t *testing.T) {
	tests := []struct {
		name string
		want *Manager
	}{
		{
			name: "NewCryptoManager",
			want: &Manager{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewCryptoManager(), "NewCryptoManager()")
		})
	}
}
