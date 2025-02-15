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
}

type GopsutilMonitor struct {
	TotalMemory,
	FreeMemory uint64
	CPUutilization map[int]models.Gauge
}

func NewMonitor() Monitor {
	var m Monitor
	var rtm runtime.MemStats
	// Read full mem stats
	runtime.ReadMemStats(&rtm)

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
	m.InitCPUutilizationValue()
	return m
}

func NewGopsutilMonitor() (GopsutilMonitor, error) {
	var gm GopsutilMonitor
	gm.InitCPUutilizationValue()
	// Получаем информацию о памяти
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "gopsutil mem.VirtualMemory")
		return gm, err
	}

	// Получаем количество логических CPU
	_, err = cpu.Counts(true)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "gopsutil cpu.Counts")
		return gm, err
	}

	// Получаем использование CPU за 1 секунду
	cpuPercent, err := cpu.Percent(time.Second, true)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "gopsutil cpu.Percent")
		return gm, err
	}

	// использование CPU для каждого ядра
	for key, percent := range cpuPercent {
		gm.StoreCPUutilizationValue(key, models.Gauge(percent))
	}

	// TotalMemory - общий объем памяти
	gm.TotalMemory = vmStat.Total
	// FreeMemory - свободный объем памяти
	gm.FreeMemory = vmStat.Free
	//точное количество — по числу CPU, определяемому во время исполнения

	return gm, nil
}

func (m *Monitor) InitCPUutilizationValue() {
	var mx sync.RWMutex
	mx.Lock()
	m.CPUutilization = make(map[int]models.Gauge)
	defer mx.Unlock()
}

func (m *Monitor) StoreCPUutilizationValue(key int, value models.Gauge) {
	var mx sync.RWMutex
	mx.Lock()
	m.CPUutilization[key] = value
	mx.Unlock()
}

func (m *Monitor) GetAllCPUutilization() map[int]models.Gauge {
	var mx sync.RWMutex

	mx.RLock()
	mapCopy := make(map[int]models.Gauge, len(m.CPUutilization))
	for key, val := range m.CPUutilization {
		mapCopy[key] = val
	}
	mx.RUnlock()

	return mapCopy
}

func (gm *GopsutilMonitor) InitCPUutilizationValue() {
	var mx sync.RWMutex
	mx.Lock()
	gm.CPUutilization = make(map[int]models.Gauge)
	mx.Unlock()
}

func (gm *GopsutilMonitor) StoreCPUutilizationValue(key int, value models.Gauge) {
	var mx sync.RWMutex
	mx.Lock()
	gm.CPUutilization[key] = value
	mx.Unlock()
}

func (gm *GopsutilMonitor) GetAllCPUutilization() map[int]models.Gauge {
	var mx sync.RWMutex

	mx.RLock()
	mapCopy := make(map[int]models.Gauge, len(gm.CPUutilization))
	for key, val := range gm.CPUutilization {
		mapCopy[key] = val
	}
	mx.RUnlock()

	return mapCopy
}
