// Package storage основные структуры для работы с данными метрик
package storage

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

// Monitor содержит все метрики
type Monitor struct {
	Alloc,
	BuckHashSys,
	Frees,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	Mallocs,
	NextGC,
	OtherSys,
	PauseTotalNs,
	StackInuse,
	Sys,
	StackSys,
	MSpanInuse,
	MSpanSys,
	TotalMemory,
	FreeMemory,
	TotalAlloc uint64
	GCCPUFraction float64
	NumForcedGC,
	NumGC uint32
	PollCount      models.Counter
	RandomValue    models.Gauge
	CPUutilization map[int]models.Gauge
	mx             sync.RWMutex
}

func NewMonitor() *Monitor {
	var m Monitor
	SetMetricsToMonitor(&m)
	return &m
}

// SetMetricsToMonitor заполняет монитор основными метриками
func SetMetricsToMonitor(m *Monitor) {
	var rtm runtime.MemStats
	// Read full mem stats
	runtime.ReadMemStats(&rtm)
	m.mx.Lock()
	// Misc memory stats
	m.Alloc = rtm.Alloc
	m.TotalAlloc = rtm.TotalAlloc
	m.Sys = rtm.Sys
	m.Mallocs = rtm.Mallocs
	m.Frees = rtm.Frees

	m.StackSys = rtm.StackSys
	m.MSpanInuse = rtm.MSpanInuse
	m.MSpanSys = rtm.MSpanSys
	m.BuckHashSys = rtm.BuckHashSys
	m.Frees = rtm.Frees
	m.GCCPUFraction = rtm.GCCPUFraction
	m.GCSys = rtm.GCSys
	m.HeapAlloc = rtm.HeapAlloc
	m.HeapIdle = rtm.HeapIdle
	m.HeapInuse = rtm.HeapInuse
	m.HeapObjects = rtm.HeapObjects
	m.HeapReleased = rtm.HeapReleased
	m.HeapSys = rtm.HeapSys
	m.LastGC = rtm.LastGC
	m.Lookups = rtm.Lookups
	m.MCacheInuse = rtm.MCacheInuse
	m.MCacheSys = rtm.MCacheSys
	m.Mallocs = rtm.Mallocs
	m.NextGC = rtm.NextGC
	m.NumForcedGC = rtm.NumForcedGC
	m.OtherSys = rtm.OtherSys
	m.PauseTotalNs = rtm.PauseTotalNs
	m.StackInuse = rtm.StackInuse
	m.Sys = rtm.Sys
	m.TotalAlloc = rtm.TotalAlloc
	m.NumGC = rtm.NumGC
	m.RandomValue = models.Gauge(rand.Float64())
	m.mx.Unlock()
	m.InitCPUutilizationValue()
}

// SetGopsutilMetricsToMonitor заполняет метриками gopsutil
func SetGopsutilMetricsToMonitor(m *Monitor) error {
	// Получаем информацию о памяти

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "gopsutil mem.VirtualMemory")
		return err
	}

	// Получаем количество логических CPU
	_, err = cpu.Counts(true)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "gopsutil cpu.Counts")
		return err
	}

	// Получаем использование CPU за 1 наносекунду
	cpuPercent, err := cpu.Percent(10*time.Millisecond, true)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "gopsutil cpu.Percent")
		return err
	}

	// использование CPU для каждого ядра
	for key, percent := range cpuPercent {
		m.StoreCPUutilizationValue(key, models.Gauge(percent))
	}

	// TotalMemory - общий объем памяти
	m.TotalMemory = vmStat.Total
	// FreeMemory - свободный объем памяти
	m.FreeMemory = vmStat.Free
	//точное количество — по числу CPU, определяемому во время исполнения

	return nil
}

// InitCPUutilizationValue инициализирует мапу хранящую данные по CPU
func (m *Monitor) InitCPUutilizationValue() {
	m.mx.Lock()
	m.CPUutilization = make(map[int]models.Gauge)
	m.mx.Unlock()
}

func (m *Monitor) StoreCountValue(value int) {
	m.mx.Lock()
	m.PollCount = models.Counter(value)
	m.mx.Unlock()
}

func (m *Monitor) GetCountValue() models.Counter {
	m.mx.RLock()
	val := m.PollCount
	m.mx.RUnlock()
	return val
}

func (m *Monitor) StoreCPUutilizationValue(key int, value models.Gauge) {
	m.mx.Lock()
	m.CPUutilization[key] = value
	m.mx.Unlock()
}

func (m *Monitor) GetAllCPUutilization() map[int]models.Gauge {
	m.mx.RLock()
	mapCopy := make(map[int]models.Gauge, len(m.CPUutilization))
	for key, val := range m.CPUutilization {
		mapCopy[key] = val
	}
	m.mx.RUnlock()
	return mapCopy
}
