package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigParams(t *testing.T) {
	type args struct {
		consoleKeyShort string
		consoleKeyFull  string
		configType      string
	}
	tests := []struct {
		args args
		want ParamsProvider
		name string
	}{
		{
			name: "TestNewConfigParams",
			args: args{
				consoleKeyShort: "default",
				consoleKeyFull:  "default",
				configType:      "default",
			},
			want: &Params{
				consoleKeyShort: "default",
				consoleKeyFull:  "default",
				configType:      "default",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewConfigParams(tt.args.consoleKeyShort, tt.args.consoleKeyFull, tt.args.configType), "NewConfigParams(%v, %v, %v)", tt.args.consoleKeyShort, tt.args.consoleKeyFull, tt.args.configType)
		})
	}
}

func TestParams_GetConfigType(t *testing.T) {
	type fields struct {
		consoleKeyShort string
		consoleKeyFull  string
		configType      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestGetConfigType",
			fields: fields{
				consoleKeyShort: "default",
				consoleKeyFull:  "default",
				configType:      "default",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{
				consoleKeyShort: tt.fields.consoleKeyShort,
				consoleKeyFull:  tt.fields.consoleKeyFull,
				configType:      tt.fields.configType,
			}
			assert.Equalf(t, tt.want, p.GetConfigType(), "GetConfigType()")
		})
	}
}

func TestParams_GetConsoleKeyFull(t *testing.T) {
	type fields struct {
		consoleKeyShort string
		consoleKeyFull  string
		configType      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestGetConsoleKeyFull",
			fields: fields{
				consoleKeyShort: "default",
				consoleKeyFull:  "default",
				configType:      "default",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{
				consoleKeyShort: tt.fields.consoleKeyShort,
				consoleKeyFull:  tt.fields.consoleKeyFull,
				configType:      tt.fields.configType,
			}
			assert.Equalf(t, tt.want, p.GetConsoleKeyFull(), "GetConsoleKeyFull()")
		})
	}
}

func TestParams_GetConsoleKeyShort(t *testing.T) {
	type fields struct {
		consoleKeyShort string
		consoleKeyFull  string
		configType      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestGetConsoleKeyShort",
			fields: fields{
				consoleKeyShort: "default",
				consoleKeyFull:  "default",
				configType:      "default",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{
				consoleKeyShort: tt.fields.consoleKeyShort,
				consoleKeyFull:  tt.fields.consoleKeyFull,
				configType:      tt.fields.configType,
			}
			assert.Equalf(t, tt.want, p.GetConsoleKeyShort(), "GetConsoleKeyShort()")
		})
	}
}
