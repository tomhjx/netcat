package core

import "github.com/google/gopacket"

type Source struct {
	transport gopacket.Flow
	payload   []byte
}

type Resolved struct {
	isClientFlow bool
	code         string
	content      string
}
