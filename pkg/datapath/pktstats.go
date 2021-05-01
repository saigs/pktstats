//
// primitivitives for high-speed packet statistics processing
//
package datapath

import (
	"log"
	"sync/atomic"
)

var (
	// TBD make this into log.Printf()
	debug = false
)

func (p *DatapathContext) GetCount() int {
	return int(p.stat.count)
}

func (p *DatapathContext) GetLastUpdateCount() int {
	return int(p.stat.lastUpdateCount)
}

func (p *DatapathContext) GetAveragePacketSize() int {
	return int(p.stat.avgPktSize)
}

// update stats with mutex lock
func (p *DatapathContext) UpdateStatsMutex(th int, pktSize int) {
	p.rwlock.Lock()
	if p.GetCount() < 1 {
		p.stat.avgPktSize *= p.stat.count
	} else {
		p.stat.avgPktSize = (p.stat.avgPktSize*p.stat.count + uint64(pktSize)) / (p.stat.count + 1)
	}
	p.stat.count++
	p.rwlock.Unlock()

	p.ObtainStats(th)
}

// update stats atomic
func (p *DatapathContext) UpdateStatsAtomic(th int, pktSize int) {
	atomic.AddUint64(&p.stat.avgPktSize, uint64(pktSize))
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
	atomic.AddUint64(&p.cstat[th].avgPktSize, uint64(pktSize))
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
	ret := uint64(p.GetAveragePacketSize() / p.GetCount())
	if debug {
		log.Printf("size average %d [%d pkts]\n", ret, p.GetCount())
	}
	return ret
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
		s, c uint64
	)
	for i := 0; i < GetMaxThreads(); i++ {
		s += atomic.LoadUint64(&p.cstat[i].avgPktSize)
		c += atomic.LoadUint64(&p.cstat[i].count)
	}
	ret := s / c
	if debug {
		log.Printf("size average %d [%d pkts]\n", ret, p.GetCount())
	}
	return ret
}

func (p *DatapathContext) ObtainStats(th int) {
	// time to obtain stats?
	if (p.GetCount() - p.GetLastUpdateCount()) >= GetStatsFreq() {
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
