package core

import (
	"errors"
	"io"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/tomhjx/netcat/protocol"
)

type Parser struct {
	protocol protocol.Driver
}

func NewParser(protocol protocol.Driver) *Parser {

	return &Parser{protocol: protocol}
}

func (me *Parser) Resolve(packet gopacket.Packet) *Resolved {
	appLayer := packet.ApplicationLayer()
	if appLayer == nil {
		return nil
	}
	ret := Resolved{}
	// Network Layer
	network := packet.NetworkLayer().NetworkFlow()

	// Transport Layer
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return nil
	}
	tcp, _ := tcpLayer.(*layers.TCP)

	if len(appLayer.Payload()) == 0 {
		return nil
	}

	var (
		code    string
		content string
	)

	if int(tcp.SrcPort) == me.protocol.Port() {
		code, content = me.protocol.ResolveServer(appLayer.Payload())
		ret.isClientFlow = false
	} else {
		code, content = me.protocol.ResolveClient(appLayer.Payload())
		ret.isClientFlow = true
	}

	if code == "" {
		return nil
	}

	ret.srcHost = network.Src().String()
	ret.srcPort = int(tcp.SrcPort)
	ret.dstHost = network.Dst().String()
	ret.dstPort = int(tcp.DstPort)
	ret.code = code
	ret.content = content
	ret.seq = tcp.Seq

	return &ret
}

func (me *Parser) resolvePacketTo(r io.Reader, w io.Writer) (uint8, error) {

	header := make([]byte, 4)
	if n, err := io.ReadFull(r, header); err != nil {
		if n == 0 && err == io.EOF {
			return 0, io.EOF
		}
		return 0, errors.New("ERR : Unknown stream")
	}

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)

	var seq uint8
	seq = header[3]

	if n, err := io.CopyN(w, r, int64(length)); err != nil {
		return 0, errors.New("ERR : Unknown stream")
	} else if n != int64(length) {
		return 0, errors.New("ERR : Unknown stream")
	} else {
		return seq, nil
	}

	return seq, nil
}
