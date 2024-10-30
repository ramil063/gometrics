package filer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ramil063/gometrics/cmd/agent/storage"
)

func TestNewReader(t *testing.T) {
	filename := "../../../internal/storage/files/test.json"
	r, _ := NewReader(filename)

	tests := []struct {
		name     string
		filename string
		want     *Reader
		wantErr  bool
	}{
		{"test 1", filename, r, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, _ := NewWriter(tt.filename)
			var m = storage.NewMonitor()
			err := w.WriteMonitor(&m)
			if err != nil {
				return
			}

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

func TestReader_ReadMonitor(t *testing.T) {

	filename := "../../../internal/storage/files/test.json"

	tests := []struct {
		name     string
		filename string
	}{
		{"test 1", filename},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := NewReader(tt.filename)
			w, _ := NewWriter(tt.filename)

			var m = storage.NewMonitor()
			err := w.WriteMonitor(&m)
			if err != nil {
				return
			}
			got, err := r.ReadMonitor()
			assert.NoError(t, err, "ReadMonitor error")
			assert.Equal(t, got.PollCount, m.PollCount, "equal metrics")
		})
	}
}
