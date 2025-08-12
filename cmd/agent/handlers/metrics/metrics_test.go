package metrics

import (
	"sync"
	"testing"

	"github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/stretchr/testify/assert"
)

func BenchmarkCollectMetricsRequestBodies(b *testing.B) {
	var m storage.Monitor

	for i := 0; i < b.N; i++ {
		CollectMetricsRequestBodies(&m)
	}
}

func BenchmarkCollectMonitorMetrics(b *testing.B) {
	var m storage.Monitor
	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		CollectMonitorMetrics(1, &m, &wg)
	}
}

func BenchmarkCollectGopsutilMetrics(b *testing.B) {
	var m storage.Monitor
	var wg sync.WaitGroup
	storage.SetMetricsToMonitor(&m)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CollectGopsutilMetrics(&m, &wg)
	}
}

func TestCollectMetricsRequestBodies(t *testing.T) {
	tests := []struct {
		name    string
		monitor *storage.Monitor
		want    []byte
	}{
		{
			name:    "test 1",
			monitor: &storage.Monitor{},
			want:    []byte("[{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":0},{\"id\":\"BuckHashSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"Frees\",\"type\":\"gauge\",\"value\":0},{\"id\":\"GCSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapAlloc\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapIdle\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapInuse\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapObjects\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapReleased\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"LastGC\",\"type\":\"gauge\",\"value\":0},{\"id\":\"Lookups\",\"type\":\"gauge\",\"value\":0},{\"id\":\"MCacheInuse\",\"type\":\"gauge\",\"value\":0},{\"id\":\"MCacheSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"Mallocs\",\"type\":\"gauge\",\"value\":0},{\"id\":\"NextGC\",\"type\":\"gauge\",\"value\":0},{\"id\":\"OtherSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"PauseTotalNs\",\"type\":\"gauge\",\"value\":0},{\"id\":\"StackInuse\",\"type\":\"gauge\",\"value\":0},{\"id\":\"Sys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"StackSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"MSpanInuse\",\"type\":\"gauge\",\"value\":0},{\"id\":\"MSpanSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"TotalMemory\",\"type\":\"gauge\",\"value\":0},{\"id\":\"FreeMemory\",\"type\":\"gauge\",\"value\":0},{\"id\":\"TotalAlloc\",\"type\":\"gauge\",\"value\":0},{\"id\":\"GCCPUFraction\",\"type\":\"gauge\",\"value\":0},{\"id\":\"NumForcedGC\",\"type\":\"gauge\",\"value\":0},{\"id\":\"NumGC\",\"type\":\"gauge\",\"value\":0},{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":0},{\"id\":\"RandomValue\",\"type\":\"gauge\",\"value\":0}]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodies := CollectMetricsRequestBodies(tt.monitor)
			assert.Equal(t, tt.want, bodies)
		})
	}
}
