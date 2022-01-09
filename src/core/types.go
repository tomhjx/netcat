package core

import "github.com/google/gopacket"

type Source struct {
	transport gopacket.Flow
	payload   []byte
}

type Resolved struct {
	isClientFlow bool
	srcHost      string
	srcPort      int
	dstHost      string
	dstPort      int
	code         string
	content      string
	seq          uint32
}
