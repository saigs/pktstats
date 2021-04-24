//
// primitivitives for high-speed packet statistics processing
//
package pktstats

import (
	"fmt"
	"testing"
)

const (
	MAX_THREADS = 10

	// at every run, process 10,000 packets. Benchmark will xtimes this
	// as many times, so we will run into several million pkts/sec
	MAX_PKTS = 10000
)

var (
	stat PktStat
)

func BenchmarkMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runDatapath(MAX_THREADS, MAX_PKTS, stat.UpdateStatWithLock, stat.GetAveragePktSize)
	}
	fmt.Printf("DP processed %d pkts\n", stat.cnt)
}

func BenchmarkAtomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runDatapath(MAX_THREADS, MAX_PKTS, stat.UpdateStatAtomic, stat.GetAveragePktSize)
	}
	fmt.Printf("DP processed %d pkts\n", stat.cnt)
}

func BenchmarkWithMap(b *testing.B) {
	stat.smap = make(map[int]uint64, 0)
	for i := 0; i < b.N; i++ {
		runDatapath(MAX_THREADS, MAX_PKTS, stat.UpdateStatWithMap, stat.GetAveragePktSizeUsingMap)
	}
	fmt.Printf("DP processed %d pkts\n", stat.cnt)
	stat.ClearStatMap()
}
