//
// primitivitives for high-speed packet statistics processing
//
package main

import (
	"flag"
	"fmt"
	"log"

	dp "github.com/saigs/pktstats/pkg/datapath"
)

const (
	// default max threads
	DEFAULT_MAX_THREADS = 10

	// at every run, process these many packets
	DEFAULT_MAX_PKTS = 10000

	// get stats every % of updates
	DEFAULT_GET_STATS_FREQ = 10
)

type Config struct {
	mode    int
	intf    string
	threads int
	packets int
}

var (
	// config
	cfg Config

	// datapath context
	ctx = dp.CreateDatapathContext()

	// datapath config
	dpCfgList = []*dp.DatapathConfig{
		&dp.DatapathConfig{ // mutex
			ctxType:       CTX_MUTEX,
			updateStatsFn: ctxTest.UpdateStatsMutex,
			getStatsFn:    ctxTest.GetAveragePktSize,
			getStatsFreq:  DEFAULT_GET_STATS_FREQ,
			maxThreads:    DEFAULT_MAX_THREADS,
			maxPkts:       DEFAULT_MAX_PKTS,
		},
		&dp.DatapathConfig{ // atomic
			ctxType:       CTX_ATOMIC,
			updateStatsFn: ctxTest.UpdateStatsAtomic,
			getStatsFn:    ctxTest.GetAveragePktSize,
			getStatsFreq:  DEFAULT_GET_STATS_FREQ,
			maxThreads:    DEFAULT_MAX_THREADS,
			maxPkts:       DEFAULT_MAX_PKTS,
		},
		&dp.DatapathConfig{ // map
			ctxType:       CTX_MAP,
			updateStatsFn: ctxTest.UpdateStatsMap,
			getStatsFn:    ctxTest.GetAveragePktSizeMap,
			getStatsFreq:  DEFAULT_GET_STATS_FREQ,
			maxThreads:    DEFAULT_MAX_THREADS,
			maxPkts:       DEFAULT_MAX_PKTS,
		},
		&dp.DatapathConfig{ // concurrent
			ctxType:       CTX_CONCURRENT,
			updateStatsFn: ctxTest.UpdateStatsAtomicConcurrent,
			getStatsFn:    ctxTest.GetAveragePktSizeConcurrent,
			getStatsFreq:  DEFAULT_GET_STATS_FREQ,
			maxThreads:    DEFAULT_MAX_THREADS,
			maxPkts:       DEFAULT_MAX_PKTS,
		},
	}
)

func parseCmdLine() {
	flag.IntVar(&cfg.mode, "mode", "mutex", "mutex=0 atomic=1 map=2 concurrent=3,")
	flag.StringVar(&cfg.intf, "interface", "en0", "Interface to listen on e.g. \"eth0\"")
	flag.IntVar(&cfg.threads, "threads", DEFAULT_MAX_THREADS, "#concurrent threads")
	flag.IntVar(&cfg.packets, "packets", DEFAULT_MAX_PKTS, "#packets to send")
	flag.Parse()
}

func main() {
	parseCmdLine()
	if cfg.mode < 0 || cfg.mode > 3 {
		panic(fmt.Sprintf("invalid datapath mode %d", cfg.mode))
	}

	dpCfg[cfg.mode].maxThreads = cfg.threads
	dpCfg[cfg.mode].maxPkts = cfg.packets

	log.Printf("Initializing datapath..\n")
	dp.InitDatapath(dpCfg[cfg.mode])
	log.Printf("Running datapath threads %d, pkts %d\n",
		dpCfg[cfg.mode].maxThreads, dpCfg[cfg.mode].maxPkts)
	dp.RunDatapath()
}
