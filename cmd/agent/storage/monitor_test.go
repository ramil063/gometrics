package storage

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewMonitor(t *testing.T) {
	tests := []struct {
		name string
		want reflect.Kind
	}{
		{"check monitor", reflect.ValueOf(Monitor{}).Kind()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMonitor()
			assert.Equal(t, reflect.ValueOf(m).Kind(), tt.want)
		})
	}
}
