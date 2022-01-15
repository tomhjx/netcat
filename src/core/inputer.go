package core

import (
	"log"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/tomhjx/netcat/protocol"
)

type Inputer struct {
	protocol protocol.Driver
}

var (
	readTriggers     = []func(gopacket.Packet){}
	readDoneTriggers = []func(){}
)

func NewInputer(protocol protocol.Driver) *Inputer {
	return &Inputer{protocol: protocol}
}
func (me *Inputer) RegisterReadTrigger(trigger func(gopacket.Packet)) {
	readTriggers = append(readTriggers, trigger)
}

func (me *Inputer) RegisterReadDoneTrigger(trigger func()) {
	readDoneTriggers = append(readDoneTriggers, trigger)
}

func (me *Inputer) callReadTrigger(s gopacket.Packet) {
	for _, trigger := range readTriggers {
		trigger(s)
	}
}

func (me *Inputer) callReadDoneTrigger() {
	for _, trigger := range readDoneTriggers {
		trigger()
	}
}

func (me *Inputer) ReadOffline(pcapfile string) {

	var (
		handle *pcap.Handle
		err    error
		// Will reuse these for each packet
		// ethLayer layers.Ethernet
		// ipLayer  layers.IPv4
		// tcpLayer layers.TCP
		// llcLayer layers.LLC
	)

	handle, err = pcap.OpenOffline(pcapfile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	handle.SetBPFFilter(me.protocol.TransportType() + " and port " + strconv.Itoa(me.protocol.Port()))
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {

		if packet == nil {
			continue
		}
		if packet.NetworkLayer() == nil ||
			packet.TransportLayer() == nil ||
			packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			continue
		}

		applicationLayer := packet.ApplicationLayer()
		if applicationLayer == nil {
			continue
		}
		me.callReadTrigger(packet)
	}
	me.callReadDoneTrigger()

}
