package core

import (
	"encoding/json"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	log "github.com/sirupsen/logrus"
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
		content      protocol.Content
		isClientFlow bool
	)

	if int(tcp.SrcPort) == me.protocol.Port() {
		log.Info("come from server.")
		content = me.protocol.ResolveServer(appLayer.Payload())
		isClientFlow = false
	} else {
		log.Info("come from client.")
		content = me.protocol.ResolveClient(appLayer.Payload())
		isClientFlow = true
	}

	if content == nil {
		return nil
	}

	ret := &Resolved{
		IsClientFlow: isClientFlow,
		SrcHost:      network.Src().String(),
		SrcPort:      int(tcp.SrcPort),
		DstHost:      network.Dst().String(),
		DstPort:      int(tcp.DstPort),
		Content:      content,
		Seq:          tcp.Seq,
	}

	cjson, _ := json.Marshal(ret)
	log.Debug(content.Section() + ":" + string(cjson))

	return ret
}
