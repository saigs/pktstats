//
// primitivitives for high-speed packet statistics processing
//
package datapath

import (
	"log"
	"sync/atomic"
)

var (
	debug = false
)

func (p *DatapathContext) GetCount() uint64 {
	return uint64(p.stat.count)
}

func (p *DatapathContext) GetLastUpdateCount() uint64 {
	return uint64(p.stat.lastUpdateCount)
}

func (p *DatapathContext) GetAveragePacketSize() uint64 {
	switch GetCtxType() {
	case CTX_MUTEX:
		return uint64(p.stat.avgPktSize)

	case CTX_ATOMIC:
		return uint64(p.stat.avgPktSize)

	case CTX_MAP:
		return p.GetAveragePktSizeMap()

	case CTX_CONCURRENT:
		return p.GetAveragePktSizeConcurrent()
	}
	return 0
}

// update stats with mutex lock
func (p *DatapathContext) UpdateStatsMutex(th int, pktSize int) {
	p.rwlock.Lock()
	c := p.GetCount()
	if c < 1 {
		p.stat.avgPktSize = uint64(pktSize)
	} else {
		p.stat.avgPktSize = (p.stat.avgPktSize*c + uint64(pktSize)) / (c + 1)
	}
	p.stat.count++
	p.rwlock.Unlock()

	p.ObtainStats(th)
}

// update stats atomic
func (p *DatapathContext) UpdateStatsAtomic(th int, pktSize int) {
	var (
		s, c uint64
	)
	s = atomic.LoadUint64(&p.stat.avgPktSize)
	c = atomic.LoadUint64(&p.stat.count)
	atomic.StoreUint64(&p.stat.avgPktSize, (s*c+uint64(pktSize))/(c+1))
	atomic.AddUint64(&p.stat.count, 1)

	p.ObtainStats(th)
}

// update stats with map
func (p *DatapathContext) UpdateStatsMap(th int, pktSize int) {
	p.rwlock.Lock()
	p.smap[pktSize]++
	p.stat.count++
	p.rwlock.Unlock()

	p.ObtainStats(th)
}

// update stats atomic concurrent
func (p *DatapathContext) UpdateStatsAtomicConcurrent(th int, pktSize int) {
	var (
		s, c uint64
	)
	s = atomic.LoadUint64(&p.cstat[th].avgPktSize)
	c = atomic.LoadUint64(&p.cstat[th].count)
	atomic.StoreUint64(&p.cstat[th].avgPktSize, (s*c+uint64(pktSize))/(c+1))
	atomic.AddUint64(&p.cstat[th].count, 1)

	p.ObtainStats(th)
}

func (p *DatapathContext) Init() {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()
	p.stat.avgPktSize = 0
	p.stat.count = 0
	p.stat.lastUpdateCount = 0
	p.smap = make(map[int]uint64, 0)
	p.cstat = make([]PktStat, GetMaxThreads())
}

func (p *DatapathContext) Reset() {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()
	for k := range p.smap {
		delete(p.smap, k)
	}
	p.stat.avgPktSize = 0
	p.stat.count = 0
	p.stat.lastUpdateCount = 0
	p.cstat = p.cstat[:0]
}

// get average packet size
func (p *DatapathContext) GetAveragePktSize() uint64 {
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()
	if p.GetCount() <= 0 {
		return 0
	}
	ret := p.GetAveragePacketSize()
	if debug {
		log.Printf("average pkt size %d, #pkts %d\n", ret, p.GetCount())
	}
	return uint64(ret)
}

// get average packet size using map
func (p *DatapathContext) GetAveragePktSizeMap() uint64 {
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()

	if len(p.smap) == 0 {
		return 0
	}

	var (
		s uint64
		c uint64
	)

	for k, v := range p.smap {
		s += uint64(k) * v
		c += uint64(v)
	}
	ret := s / c
	if debug {
		log.Printf("size average %d [%d pkts]\n", ret, c)
	}
	return ret
}

// get average packet size concurrent
func (p *DatapathContext) GetAveragePktSizeConcurrent() uint64 {
	var (
		s uint64
	)
	for i := 0; i < GetMaxThreads(); i++ {
		s += atomic.LoadUint64(&p.cstat[i].avgPktSize)
	}
	ret := s / uint64(GetMaxThreads())
	if debug {
		log.Printf("size average %d [%d pkts]\n", ret, p.GetCount())
	}
	return ret
}

func (p *DatapathContext) ObtainStats(th int) {
	// time to obtain stats?
	if (p.GetCount() - p.GetLastUpdateCount()) >= uint64(GetStatsFreq()) {
		if fn := GetStatsFn(); fn != nil {
			fn()
		}
		p.stat.lastUpdateCount = p.stat.count
		if debug {
			log.Printf("stats prev [%d] current [%d]\n",
				p.stat.lastUpdateCount, p.GetCount())
		}
	}
}
