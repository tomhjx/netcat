package core

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
