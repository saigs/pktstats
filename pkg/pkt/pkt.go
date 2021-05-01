package pkt

import (
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/ghedo/go.pkt/packet"
	"github.com/ghedo/go.pkt/packet/eth"
	"github.com/ghedo/go.pkt/packet/ipv4"
	"github.com/ghedo/go.pkt/packet/raw"
	"github.com/ghedo/go.pkt/packet/tcp"
)

const (
	DEFAULT_MAX_PKT_SIZE = 32766
)

type Packet struct {
	Pkt *eth.Packet
}

//
// pack a packet into bytes
//
func (pkt *Packet) Pack() ([]byte, error) {
	b := packet.Buffer{}
	if err := pkt.Pkt.Pack(&b); err != nil {
		log.Fatal("error: packing failure")
		return b.Bytes(), err
	}
	return b.Bytes(), nil
}

func (pkt *Packet) GetLength() uint16 {
	return pkt.Pkt.GetLength()
}

func (pkt *Packet) String() string {
	return pkt.Pkt.String()
}

//
// Make a TCP packet with random payload
//
func MakeRandomTCPPkt() *Packet {
	// data payload
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, r.Intn(DEFAULT_MAX_PKT_SIZE)+1)
	data := &raw.Packet{
		Data: b,
	}

	// TCP
	tcp := &tcp.Packet{
		SrcPort:    20,
		DstPort:    80,
		Seq:        5400,
		Ack:        432,
		DataOff:    5,
		Flags:      tcp.Syn,
		WindowSize: 8192,
		Urgent:     40,
	}
	tcp.SetPayload(data)

	// IPv4
	ip4 := ipv4.Make()
	ip4.SrcAddr = net.ParseIP("10.1.1.1")
	ip4.DstAddr = net.ParseIP("10.1.1.2")
	ip4.SetPayload(tcp)

	// eth
	hwSrc, _ := net.ParseMAC("5e:16:f7:e4:42:3f")
	hwDst, _ := net.ParseMAC("1e:00:d2:2e:1a:89")
	eth := &eth.Packet{
		SrcAddr: hwSrc,
		DstAddr: hwDst,
		Type:    eth.IPv4,
	}
	eth.SetPayload(ip4)
	return &Packet{
		Pkt: eth,
	}
}
