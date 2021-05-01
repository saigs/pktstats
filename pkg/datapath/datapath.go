//
// primitivitives for high-speed packet statistics processing
//
package datapath

import (
	"fmt"
	"log"
	"sync"

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
		ctxType       DatapathContextType // context/lock type
		updateStatsFn updateStatsFunc     // update stats routine
		getStatsFn    getStatsFunc        // get stats function
		getStatsFreq  int                 // get stats every nth updates e.g. 1/2, 1/3,...
		maxThreads    int                 // max supported threads
		maxPkts       int                 // max packets handled total
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
	if c.ctxType < CTX_MUTEX || c.ctxType > CTX_CONCURRENT {
		return fmt.Errorf("invalid ctxType %d", c.ctxType)
	}
	if c.maxPkts <= 1 {
		return fmt.Errorf("invalid maxPkts %d", c.maxPkts)
	}
	if c.maxThreads < 1 {
		return fmt.Errorf("invalid maxThreads %d", c.maxThreads)
	}
	if c.updateStatsFn == nil || c.getStatsFn == nil {
		return fmt.Errorf("invalid nil update/get functions")
	}
	return nil
}

func GetStatsFreq() int {
	if dpCfg == nil {
		return 0
	}
	return dpCfg.getStatsFreq
}

func GetStatsFn() getStatsFunc {
	if dpCfg == nil {
		return nil
	}
	return dpCfg.getStatsFn
}

func GetMaxThreads() int {
	if dpCfg == nil {
		return 0
	}
	return dpCfg.maxThreads
}

func GetMaxPkts() int {
	if dpCfg == nil {
		return 0
	}
	return dpCfg.maxPkts
}

//
// populate packet queue
//
func populatePacketQueue() error {
	if ctx == nil || dpCfg == nil {
		return fmt.Errorf("nil datapath config or context")
	}

	// spray packets across threads
	pq = make([]*que.Queue, dpCfg.maxThreads)
	m := int(dpCfg.maxPkts / dpCfg.maxThreads)
	for i := 0; i < dpCfg.maxThreads; i++ {
		if pq[i] = NewQueue(m); pq[i] == nil {
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
		ctxType:       cfg.ctxType,
		updateStatsFn: cfg.updateStatsFn,
		getStatsFn:    cfg.getStatsFn,
		getStatsFreq:  cfg.getStatsFreq,
		maxThreads:    cfg.maxThreads,
		maxPkts:       cfg.maxPkts,
	}

	// initialize context
	ctx.Init()

	// populate packet queue
	populatePacketQueue()

	log.Printf("Init datapath: type %s, threads %d, total pkts %d\n",
		dpCfg.ctxType.String(), dpCfg.maxThreads, dpCfg.maxPkts)
	return nil
}

//
// get next packet, returns pkt,size
//
func GetNextPacket(t int) (interface{}, int) {
	if pq[t].IsEmpty() {
		return nil, 0
	}
	if pp, err := pq[t].Dequeue(); err == nil {
		p := pp.(pkt.Packet)
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
	dpCfg.updateStatsFn(th, sz)
	return nil
}

//
// postprocess a packet
//
func PostprocessPacket(th int, p interface{}, sz int) error {
	return nil
}

//
// RunDatapath simulates a datapath experience with threads x packets,
// shares the same packet statistics structure across
//
func RunDatapath() {
	if ctx == nil || dpCfg == nil {
		panic("nil datapath config or context, exiting..\n")
	}

	var (
		wg sync.WaitGroup
	)

	// concurrent processing routine
	processPkt := func(th int) {
		defer wg.Done()
		// handle packets from queue
		for p := 0; p < len(pq[th]); p++ {
			pkt, sz := GetNextPacket(th)
			PreprocessPacket(th, pkt, sz)
			ProcessPacket(th, pkt, sz)
			PostprocessPacket(th, pkt, sz)
		}
	}

	// main processing
	wg.Add(dpCfg.maxThreads)
	for th := 0; th < dpCfg.maxThreads; th++ {
		go processPkt(th)
	}
	wg.Wait()
}
