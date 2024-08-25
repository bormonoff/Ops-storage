package collector

import (
	"runtime"
)

type counter string

const (
	alloc         counter = "Alloc"
	buckHashSys   counter = "BuckHashSys"
	frees         counter = "Frees"
	gcCPUFraction counter = "GCCPUFraction"
	gcSys         counter = "GCSys"
	heapAlloc     counter = "HeapAlloc"
	heapIdle      counter = "HeapIdle"
	heapInuse     counter = "HeapInuse"
	heapObjects   counter = "HeapObjects"
	heapReleased  counter = "HeapReleased"
	heapSys       counter = "HeapSys"
	lastGC        counter = "LastGC"
	lookups       counter = "Lookups"
	mCacheInuse   counter = "MCacheInuse"
	mCacheSys     counter = "MCacheSys"
	mSpanInuse    counter = "MSpanInuse"
	mSpanSys      counter = "MSpanSys"
	mallocs       counter = "Mallocs"
	nextGC        counter = "NextGC"
	numForcedGC   counter = "NumForcedGC"
	numGC         counter = "NumGC"
	otherSys      counter = "OtherSys"
	pauseTotalNs  counter = "PauseTotalNs"
	stackInuse    counter = "StackInuse"
	stackSys      counter = "StackSys"
	sys           counter = "Sys"
	totalAlloc    counter = "TotalAlloc"
)

type runtimeStats struct {
	UintStats  map[counter]uint64
	FloatStats map[counter]float64
}

type Collection struct {
	memStats runtime.MemStats

	RuntimeStats runtimeStats
	PollCount    int
	RandomVal    float64
}

func NewCollection() Collection {
	collection := Collection{
		memStats: runtime.MemStats{},
		RuntimeStats: runtimeStats{
			UintStats: map[counter]uint64{
				alloc:        0,
				buckHashSys:  0,
				frees:        0,
				gcSys:        0,
				heapAlloc:    0,
				heapIdle:     0,
				heapInuse:    0,
				heapObjects:  0,
				heapReleased: 0,
				heapSys:      0,
				lastGC:       0,
				lookups:      0,
				mCacheInuse:  0,
				mCacheSys:    0,
				mSpanInuse:   0,
				mSpanSys:     0,
				mallocs:      0,
				nextGC:       0,
				numForcedGC:  0,
				numGC:        0,
				otherSys:     0,
				pauseTotalNs: 0,
				stackInuse:   0,
				stackSys:     0,
				sys:          0,
				totalAlloc:   0,
			},
			FloatStats: map[counter]float64{
				gcCPUFraction: 0.0,
			},
		},
	}
	return collection
}

func (c *Collection) RefreshStats() {
	runtime.ReadMemStats(&c.memStats)

	c.RuntimeStats.UintStats[alloc] = c.memStats.Alloc
	c.RuntimeStats.UintStats[buckHashSys] = c.memStats.BuckHashSys
	c.RuntimeStats.UintStats[frees] = c.memStats.Frees
	c.RuntimeStats.UintStats[gcSys] = c.memStats.GCSys
	c.RuntimeStats.UintStats[heapAlloc] = c.memStats.HeapAlloc
	c.RuntimeStats.UintStats[heapIdle] = c.memStats.HeapIdle
	c.RuntimeStats.UintStats[heapInuse] = c.memStats.HeapInuse
	c.RuntimeStats.UintStats[heapObjects] = c.memStats.HeapObjects
	c.RuntimeStats.UintStats[heapReleased] = c.memStats.HeapReleased
	c.RuntimeStats.UintStats[heapSys] = c.memStats.HeapSys
	c.RuntimeStats.UintStats[lastGC] = c.memStats.LastGC
	c.RuntimeStats.UintStats[lookups] = c.memStats.Lookups
	c.RuntimeStats.UintStats[mCacheInuse] = c.memStats.MCacheInuse
	c.RuntimeStats.UintStats[mCacheSys] = c.memStats.MCacheSys
	c.RuntimeStats.UintStats[mSpanInuse] = c.memStats.MSpanInuse
	c.RuntimeStats.UintStats[mSpanSys] = c.memStats.MSpanSys
	c.RuntimeStats.UintStats[mallocs] = c.memStats.Mallocs
	c.RuntimeStats.UintStats[nextGC] = c.memStats.NextGC
	c.RuntimeStats.UintStats[otherSys] = c.memStats.OtherSys
	c.RuntimeStats.UintStats[pauseTotalNs] = c.memStats.PauseTotalNs
	c.RuntimeStats.UintStats[stackInuse] = c.memStats.StackInuse
	c.RuntimeStats.UintStats[stackSys] = c.memStats.StackSys
	c.RuntimeStats.UintStats[sys] = c.memStats.Sys
	c.RuntimeStats.UintStats[totalAlloc] = c.memStats.TotalAlloc
	c.RuntimeStats.UintStats[numGC] = uint64(c.memStats.NumGC)
	c.RuntimeStats.UintStats[numForcedGC] = uint64(c.memStats.NumForcedGC)

	c.RuntimeStats.FloatStats[gcCPUFraction] = c.memStats.GCCPUFraction

	c.PollCount++
}
