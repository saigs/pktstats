package pkt

import (
	"fmt"
	"testing"
)

func TestPack(t *testing.T) {
	p := MakeRandomTCPPkt()
	if _, err := p.Pack(); err != nil {
		t.Fatalf("Error packing packet: %s", err)
	}
	fmt.Printf("Packet: %v\n", p.String())
}
