//
// primitivitives for high-speed packet statistics processing
//
package datapath

import (
	"fmt"
	"log"
	"sync"
	"math/rand"
	"time"

	pkt "github.com/saigs/pktstats/pkg/pkt"
	que "github.com/saigs/pktstats/pkg/queue"
)

const (
	// various datapath locking context
	CTX_MUTEX = iota
	CTX_ATOMIC
	CTX_MAP
	CTX_CONCURRENT
)

type (
	// locking mechanism
	DatapathContextType int

	// stat structure
	PktStat struct {
		avgPktSize      uint64 // rounded average so far
		count           uint64 // count of packets seen so far
		lastUpdateCount uint64 // last update count of packet
	}

	// dynamic datapath context
	DatapathContext struct {
		stat   PktStat        // stat struct
		smap   map[int]uint64 // packet sizes -> count
		cstat  []PktStat      // stat used by atomic concurrent
		rwlock sync.RWMutex   // read-write protection mutex
	}

	// update/get function types
	updateStatsFunc func(th int, pktSize int)
	getStatsFunc    func() uint64

	// config for datapath
	DatapathConfig struct {
		CtxType       DatapathContextType // context/lock type
		UpdateStatsFn updateStatsFunc     // update stats routine
		GetStatsFn    getStatsFunc        // get stats function
		GetStatsFreq  int                 // get stats every nth updates e.g. 1/2, 1/3,...
		MaxThreads    int                 // max supported threads
		MaxPkts       int                 // max packets handled total
	}
)

var (
	ctx   *DatapathContext // datapath context
	dpCfg *DatapathConfig  // datapath config
	pq    []*que.Queue     // packet queue
)

func (t DatapathContextType) String() string {
	return [...]string{"mutex", "atomic", "map", "concurrent"}[t]
}

func CreateDatapathContext() *DatapathContext {
	ctx = &DatapathContext{}
	return GetRunningContext()
}

func GetRunningContext() *DatapathContext {
	return ctx
}

func validateConfig(c *DatapathConfig) error {
	if c == nil {
		return fmt.Errorf("invalid datapath config nil")
	}
	if c.CtxType < CTX_MUTEX || c.CtxType > CTX_CONCURRENT {
		return fmt.Errorf("invalid CtxType %d", c.CtxType)
	}
	if c.MaxPkts <= 1 {
		return fmt.Errorf("invalid MaxPkts %d", c.MaxPkts)
	}
	if c.MaxThreads < 1 {
		return fmt.Errorf("invalid MaxThreads %d", c.MaxThreads)
	}
	if c.UpdateStatsFn == nil || c.GetStatsFn == nil {
		return fmt.Errorf("invalid nil update/get functions")
	}
	return nil
}

func GetCtxType() DatapathContextType {
	if dpCfg == nil {
		return 0
	}
	return dpCfg.CtxType
}
func GetStatsFreq() int {
	if dpCfg == nil {
		return 0
	}
	return dpCfg.GetStatsFreq
}

func GetStatsFn() getStatsFunc {
	if dpCfg == nil {
		return nil
	}
	return dpCfg.GetStatsFn
}

func GetMaxThreads() int {
	if dpCfg == nil {
		return 0
	}
	return dpCfg.MaxThreads
}

func GetMaxPkts() int {
	if dpCfg == nil {
		return 0
	}
	return dpCfg.MaxPkts
}

//
// populate packet queue
//
func populatePacketQueue() error {
	if ctx == nil || dpCfg == nil {
		return fmt.Errorf("nil datapath config or context")
	}

	// spray packets across threads
	log.Printf("Populating datapath queues [%d]", dpCfg.MaxThreads)
	pq = make([]*que.Queue, dpCfg.MaxThreads)
	m := dpCfg.MaxPkts / dpCfg.MaxThreads
	for i := 0; i < dpCfg.MaxThreads; i++ {
		if pq[i] = que.NewQueue(m); pq[i] == nil {
			return fmt.Errorf("error: failed creating packet quque %d", i)
		}
		for k := 0; k < m; k++ {
			if err := pq[i].Queue(pkt.MakeRandomTCPPkt()); err != nil {
				return fmt.Errorf("error: failed queueing packet %d to thread %d\n", k, i)
			}
		}
		log.Printf("datapath pkt queue populated: thread %d, pkts %d\n",
			i, pq[i].Len())
	}
	return nil
}

//
// initializes structures, populates datapath queues
//
func InitDatapath(cfg *DatapathConfig) error {
	if err := validateConfig(cfg); err != nil {
		return err
	}

	dpCfg = &DatapathConfig{
		CtxType:       cfg.CtxType,
		UpdateStatsFn: cfg.UpdateStatsFn,
		GetStatsFn:    cfg.GetStatsFn,
		GetStatsFreq:  cfg.GetStatsFreq,
		MaxThreads:    cfg.MaxThreads,
		MaxPkts:       cfg.MaxPkts,
	}

	// initialize context
	ctx.Init()

	// populate packet queue
	populatePacketQueue()

	rand.Seed(time.Now().UTC().UnixNano())
	log.Printf("Init datapath: type %s, threads %d, total pkts %d\n",
		dpCfg.CtxType.String(), dpCfg.MaxThreads, dpCfg.MaxPkts)
	return nil
}

//
// get next packet, returns pkt,size
//
func GetNextPacket(t int) (interface{}, int) {
	if pq[t].IsEmpty() {
		return nil, 0
	}
	if pp, err := pq[t].Dequeue(); err == nil && pp != nil {
		p := pp.(*pkt.Packet)
		return pp, p.GetLength()
	}
	return nil, 0
}

//
// preprocess a packet
//
func PreprocessPacket(th int, p interface{}, sz int) error {
	return nil
}

//
// process a packet, returns size
//
func ProcessPacket(th int, p interface{}, sz int) error {
	dpCfg.UpdateStatsFn(th, sz)
	return nil
}

//
// postprocess a packet
//
func PostprocessPacket(th int, p interface{}, sz int) error {
	return nil
}

func duration(msg string, start time.Time) {
    log.Printf("%v took %v\n", msg, time.Since(start))
}

func track(msg string) (string, time.Time) {
    return msg, time.Now()
}

//
// RunDatapath simulates a datapath experience with threads x packets,
// shares the same packet statistics structure across
//
func RunDatapath() {
	if ctx == nil || dpCfg == nil {
		panic("nil datapath config or context, exiting..\n")
	}

	defer duration(track("RunDatapath"))

	var (
		wg sync.WaitGroup
	)

	// concurrent processing routine
	processPkt := func(th int) {
		defer wg.Done()
		// handle packets from queue
		for ; !pq[th].IsEmpty(); {
			pkt, sz := GetNextPacket(th)
			PreprocessPacket(th, pkt, sz)
			ProcessPacket(th, pkt, sz)
			PostprocessPacket(th, pkt, sz)
		}
	}

	// main processing
	wg.Add(dpCfg.MaxThreads)
	for th := 0; th < dpCfg.MaxThreads; th++ {
		go processPkt(th)
	}
	wg.Wait()
}
