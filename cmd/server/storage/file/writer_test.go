package file

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWriter(t *testing.T) {
	filename := "../../../internal/storage/files/test.json"
	tests := []struct {
		want     *Writer
		name     string
		filename string
	}{
		{
			name:     "test 1",
			filename: filename,
			want:     &Writer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWriter(tt.filename)
			assert.NoError(t, err, "Error New Writer")
			assert.Equal(t, reflect.ValueOf(tt.want).Kind(), reflect.ValueOf(got).Kind())
		})
	}
}

func TestWriter_Close(t *testing.T) {
	filename := "../../../internal/storage/files/test.json"

	tests := []struct {
		name     string
		filename string
	}{
		{"test 1", filename},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := NewWriter(tt.filename)
			assert.NoError(t, err, "Error New Writer")
			assert.NoError(t, w.Close(), "Error Writer close")
		})
	}
}
