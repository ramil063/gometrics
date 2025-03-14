package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSha256(t *testing.T) {
	type args struct {
		body []byte
		key  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test 1", args{[]byte("test"), "test"}, "88cd2108b5347d973cf39cdf9053d7dd42704876d8c9a9bd8e2d168259d3ddf7"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := CreateSha256(tt.args.body, tt.args.key)
			assert.Equal(t, tt.want, hash)
		})
	}
}
