package storage

import (
	"math/rand"
	"runtime"
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
	CPUutilization1,
	GCCPUFraction float64
	NumForcedGC,
	NumGC uint32
	PollCount   models.Counter
	RandomValue models.Gauge
}

type GopsutilMonitor struct {
	TotalMemory,
	FreeMemory uint64
	CPUutilization1 float64
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
	return m
}

func NewGopsutilMonitor() (GopsutilMonitor, error) {
	var m GopsutilMonitor
	// Получаем информацию о памяти
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "gopsutil mem.VirtualMemory")
		return m, err
	}

	// Получаем количество логических CPU
	numCPU, err := cpu.Counts(true)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "gopsutil cpu.Counts")
		return m, err
	}

	// Получаем использование CPU за 1 секунду
	cpuPercent, err := cpu.Percent(time.Second, true)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "gopsutil cpu.Percent")
		return m, err
	}

	sumPercent := 0.0
	// использование CPU для каждого ядра
	for _, percent := range cpuPercent {
		sumPercent += percent
	}

	// TotalMemory - общий объем памяти
	m.TotalMemory = vmStat.Total
	// FreeMemory - свободный объем памяти
	m.FreeMemory = vmStat.Free
	//точное количество — по числу CPU, определяемому во время исполнения
	m.CPUutilization1 = sumPercent / float64(numCPU)

	return m, nil
}
