package gzip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_compressData(t *testing.T) {
	type args struct {
		data []byte
	}
	want := []byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x4a, 0x4, 0x4, 0x0, 0x0, 0xff, 0xff, 0x43, 0xbe, 0xb7, 0xe8, 0x1, 0x0, 0x0, 0x0}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"test 1", args{data: []byte("a")}, want, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompressData(tt.args.data)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}
