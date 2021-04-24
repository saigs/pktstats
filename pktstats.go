//
// primitivitives for high-speed packet statistics processing
//
package pktstats

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MAX_PKT_SIZE               = 32766
	GET_AVERAGE_SAMPLE_PERCENT = 10
)

var (
	debug = false
)

type (
	updateFn func(pktSize int)
	getFn    func() uint64
)

type PktStat struct {
	pktSize uint64
	cnt     uint64

	// map
	smap map[int]uint64

	// lock
	wlock sync.Mutex
}

// update stats with lock
func (p *PktStat) UpdateStatWithLock(pktSize int) {
	p.wlock.Lock()
	p.pktSize += uint64(pktSize)
	p.cnt++
	p.wlock.Unlock()
}

// update stats atomic
func (p *PktStat) UpdateStatAtomic(pktSize int) {
	atomic.AddUint64(&p.pktSize, uint64(pktSize))
	atomic.AddUint64(&p.cnt, 1)
}

// update stats with map
func (p *PktStat) UpdateStatWithMap(pktSize int) {
	p.wlock.Lock()
	p.smap[pktSize]++
	p.wlock.Unlock()
}

func (p *PktStat) ClearStatMap() {
	p.wlock.Lock()
	for k := range p.smap {
		delete(p.smap, k)
	}
	p.wlock.Unlock()
}

// get average packet size
func (p *PktStat) GetAveragePktSize() uint64 {
	if p.cnt == 0 {
		return 0
	}
	ret := p.pktSize / p.cnt
	if debug {
		fmt.Printf("size average %d [%d pkts]\n", ret, p.cnt)
	}
	return ret
}

func (p *PktStat) GetAveragePktSizeUsingMap() uint64 {
	if len(p.smap) == 0 {
		return 0
	}

	var (
		s uint64
		c uint64
	)

	p.wlock.Lock()
	for k, v := range p.smap {
		s += uint64(k) * v
		c += uint64(v)
	}
	p.wlock.Unlock()
	ret := s / c
	if debug {
		fmt.Printf("size average %d [%d pkts]\n", ret, c)
	}
	return ret
}

//
// runDatapath simulates a datapath experience with threads x packets
//
func runDatapath(maxThreads, maxPkts int, updateF updateFn, getF getFn) {
	maxPkts /= maxThreads
	if maxPkts <= 1 {
		maxPkts = 1
	}
	var (
		wg    sync.WaitGroup
		every = 100 / GET_AVERAGE_SAMPLE_PERCENT
	)
	wg.Add(maxThreads)
	for t := 0; t < maxThreads; t++ {
		go func() {
			var (
				s = rand.New(rand.NewSource(time.Now().UnixNano()))
				r = rand.New(s)
			)
			defer wg.Done()
			for p := 0; p < maxPkts; p++ {
				updateF(r.Intn(MAX_PKT_SIZE) + 1)
				if p%every == 0 { // get every-th time
					getF()
				}
			}
		}()
	}
	wg.Wait()
}
