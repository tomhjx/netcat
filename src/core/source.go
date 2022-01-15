package core

import "github.com/google/gopacket"

type Source struct {
	transport gopacket.Flow
	payload   []byte
}
