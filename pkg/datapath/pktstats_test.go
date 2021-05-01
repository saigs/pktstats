//
// datapath module defines datapath static configuration, and dynamic
// context of datapath operating statistics
//
package datapath

import (
	"log"
	"testing"
)

const (
	// simulate threads
	TEST_MAX_THREADS = 10

	// at every run, process x packets. Benchmark will xtimes this
	// as many times, so we will run into several million pkts/sec
	TEST_MAX_PKTS = 100000

	// get stats every % of updates e.g. 1/10
	TEST_GET_STATS_FREQ = 1000
)

var (
	// datapath context
	ctxTest = CreateDatapathContext()

	// datapath config
	dpCfgList = []*DatapathConfig{
		&DatapathConfig{ // mutex
			ctxType:       CTX_MUTEX,
			updateStatsFn: ctxTest.UpdateStatsMutex,
			getStatsFn:    ctxTest.GetAveragePktSize,
			getStatsFreq:  TEST_GET_STATS_FREQ,
			maxThreads:    TEST_MAX_THREADS,
			maxPkts:       TEST_MAX_PKTS,
		},
		&DatapathConfig{ // atomic
			ctxType:       CTX_ATOMIC,
			updateStatsFn: ctxTest.UpdateStatsAtomic,
			getStatsFn:    ctxTest.GetAveragePktSize,
			getStatsFreq:  TEST_GET_STATS_FREQ,
			maxThreads:    TEST_MAX_THREADS,
			maxPkts:       TEST_MAX_PKTS,
		},
		&DatapathConfig{ // map
			ctxType:       CTX_MAP,
			updateStatsFn: ctxTest.UpdateStatsMap,
			getStatsFn:    ctxTest.GetAveragePktSizeMap,
			getStatsFreq:  TEST_GET_STATS_FREQ,
			maxThreads:    TEST_MAX_THREADS,
			maxPkts:       TEST_MAX_PKTS,
		},
		&DatapathConfig{ // concurrent
			ctxType:       CTX_CONCURRENT,
			updateStatsFn: ctxTest.UpdateStatsAtomicConcurrent,
			getStatsFn:    ctxTest.GetAveragePktSizeConcurrent,
			getStatsFreq:  TEST_GET_STATS_FREQ,
			maxThreads:    TEST_MAX_THREADS,
			maxPkts:       TEST_MAX_PKTS,
		},
	}
)

func InitTest(ctype DatapathContextType) *DatapathConfig {
	d := dpCfgList[ctype]
	InitDatapath(d)
	return d
}

func BenchmarkMutex(b *testing.B) {
	InitTest(CTX_MUTEX)
	for i := 0; i < b.N; i++ {
		RunDatapath()
	}
	ctx := GetRunningContext()
	log.Printf("mutex datapath processed %d pkts, average size %d\n",
		ctx.GetCount(), GetStatsFn())
}

func BenchmarkAtomic(b *testing.B) {
	InitTest(CTX_ATOMIC)
	for i := 0; i < b.N; i++ {
		RunDatapath()
	}
	ctx := GetRunningContext()
	log.Printf("atomic datapath processed %d pkts, average size %d\n",
		ctx.GetCount(), GetStatsFn())
}

func BenchmarkMap(b *testing.B) {
	InitTest(CTX_MAP)
	for i := 0; i < b.N; i++ {
		RunDatapath()
	}
	ctx := GetRunningContext()
	log.Printf("map datapath processed %d pkts, average size %d\n",
		ctx.GetCount(), GetStatsFn())
}

func BenchmarkConcurrent(b *testing.B) {
	InitTest(CTX_CONCURRENT)
	for i := 0; i < b.N; i++ {
		RunDatapath()
	}
	ctx := GetRunningContext()
	log.Printf("concurrent datapath processed %d pkts, average size %d\n",
		ctx.GetCount(), GetStatsFn())
}
