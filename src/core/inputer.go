package core

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Inputer struct{}

var (
	readTriggers     = []func(*Source){}
	readDoneTriggers = []func(){}
)

func NewInputer() *Inputer {
	return &Inputer{}
}
func (ier *Inputer) RegisterReadTrigger(trigger func(*Source)) {
	readTriggers = append(readTriggers, trigger)
}

func (ier *Inputer) RegisterReadDoneTrigger(trigger func()) {
	readDoneTriggers = append(readDoneTriggers, trigger)
}

func (ier *Inputer) callReadTrigger(src *Source) {
	for _, trigger := range readTriggers {
		trigger(src)
	}
}

func (ier *Inputer) callReadDoneTrigger() {
	for _, trigger := range readDoneTriggers {
		trigger()
	}
}

func (ier *Inputer) ReadOffline(pcapfile string) {

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
	handle.SetBPFFilter("tcp and port 3306")
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
		var src = Source{}
		src.transport = packet.NetworkLayer().NetworkFlow()
		src.payload = applicationLayer.Payload()
		ier.callReadTrigger(&src)
	}
	ier.callReadDoneTrigger()
}
