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
	// default mode mutex
	DEFAULT_MODE = 0

	// default max threads
	DEFAULT_MAX_THREADS = 10

	// at every run, process these many packets
	DEFAULT_MAX_PKTS = 10000

	// get stats every % of updates
	DEFAULT_GET_STATS_FREQ = 10
)

type Config struct {
	mode    int
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
			CtxType:       dp.CTX_MUTEX,
			UpdateStatsFn: ctx.UpdateStatsMutex,
			GetStatsFn:    ctx.GetAveragePktSize,
			GetStatsFreq:  DEFAULT_GET_STATS_FREQ,
			MaxThreads:    DEFAULT_MAX_THREADS,
			MaxPkts:       DEFAULT_MAX_PKTS,
		},
		&dp.DatapathConfig{ // atomic
			CtxType:       dp.CTX_ATOMIC,
			UpdateStatsFn: ctx.UpdateStatsAtomic,
			GetStatsFn:    ctx.GetAveragePktSize,
			GetStatsFreq:  DEFAULT_GET_STATS_FREQ,
			MaxThreads:    DEFAULT_MAX_THREADS,
			MaxPkts:       DEFAULT_MAX_PKTS,
		},
		&dp.DatapathConfig{ // map
			CtxType:       dp.CTX_MAP,
			UpdateStatsFn: ctx.UpdateStatsMap,
			GetStatsFn:    ctx.GetAveragePktSizeMap,
			GetStatsFreq:  DEFAULT_GET_STATS_FREQ,
			MaxThreads:    DEFAULT_MAX_THREADS,
			MaxPkts:       DEFAULT_MAX_PKTS,
		},
		&dp.DatapathConfig{ // concurrent
			CtxType:       dp.CTX_CONCURRENT,
			UpdateStatsFn: ctx.UpdateStatsAtomicConcurrent,
			GetStatsFn:    ctx.GetAveragePktSizeConcurrent,
			GetStatsFreq:  DEFAULT_GET_STATS_FREQ,
			MaxThreads:    DEFAULT_MAX_THREADS,
			MaxPkts:       DEFAULT_MAX_PKTS,
		},
	}
)

func parseCmdLine() {
	flag.IntVar(&cfg.mode, "mode", DEFAULT_MODE, "mutex=0 atomic=1 map=2 concurrent=3")
	flag.IntVar(&cfg.threads, "threads", DEFAULT_MAX_THREADS, "#concurrent threads")
	flag.IntVar(&cfg.packets, "packets", DEFAULT_MAX_PKTS, "#packets to send")
	flag.Parse()
}

func main() {
	parseCmdLine()
	if cfg.mode < 0 || cfg.mode > 3 {
		panic(fmt.Sprintf("invalid datapath mode %d", cfg.mode))
	}

	dpCfg := dpCfgList[cfg.mode]
	dpCfg.MaxThreads = cfg.threads
	dpCfg.MaxPkts = cfg.packets

	log.Printf("Initializing datapath..\n")
	dp.InitDatapath(dpCfg)
	log.Printf("Running datapath threads %d, pkts %d\n",
		dpCfg.MaxThreads, dpCfg.MaxPkts)
	dp.RunDatapath()
	log.Printf("average packet size %d\n", ctx.GetAveragePacketSize())
}
