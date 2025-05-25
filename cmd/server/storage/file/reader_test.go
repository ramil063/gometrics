package file

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReader(t *testing.T) {
	filename := "../../../internal/storage/files/test.json"
	r, _ := NewReader(filename)

	tests := []struct {
		want     *Reader
		name     string
		filename string
		wantErr  bool
	}{
		{
			want:     r,
			name:     "test 1",
			filename: filename,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReader(tt.filename)
			assert.NoError(t, err, "NewReader error")
			assert.Equal(t, reflect.ValueOf(tt.want).Kind(), reflect.ValueOf(got).Kind(), "NewReader error in data")
		})
	}
}

func TestReader_Close(t *testing.T) {

	filename := "../../../internal/storage/files/test.json"
	tests := []struct {
		name     string
		filename string
	}{
		{"test 1", filename},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := NewReader(filename)
			assert.NoError(t, r.Close())
		})
	}
}
