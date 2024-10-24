package storage

import (
	"math/rand"
	"runtime"

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
	TotalAlloc uint64
	GCCPUFraction float64
	NumForcedGC,
	NumGC uint32
	PollCount   models.Counter
	RandomValue models.Gauge
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
