package core

import "github.com/tomhjx/netcat/protocol"

type Resolved struct {
	IsClientFlow bool
	SrcHost      string
	SrcPort      int
	DstHost      string
	DstPort      int
	Content      protocol.Content
	Seq          uint32
}
