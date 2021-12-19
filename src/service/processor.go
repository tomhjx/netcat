package service

import (
	"bytes"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/tomhjx/netcat/amqp"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func decodeApplicationLayer(applicationLayer gopacket.ApplicationLayer) {
	if applicationLayer == nil {
		return
	}
	log.Println("Application layer/Payload found.")
	// log.Printf("%s\n", hex.Dump(applicationLayer.Payload()))

	r := amqp.NewReader(bytes.NewReader(applicationLayer.Payload()))
	frame, err := r.ReadFrame()
	if err != nil {
		// log.Println(err)
		return
	}
	log.Println(frame)

}

func (proc *Processor) Run() {
	var (
		pcapfile string = "/work/resources/rabbit.pcap"
		handle   *pcap.Handle
		err      error
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

	// if err = handle.SetBPFFilter("tcp"); err != nil {
	// 	log.Fatal("BPF filter error:", err)
	// }

	// streamFactory := &rabbitStreamFactory{}
	// streamPool := tcpassembly.NewStreamPool(streamFactory)
	// assembler := tcpassembly.NewAssembler(streamPool)

	// packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource := gopacket.NewPacketSource(handle, layers.IPProtocolICMPv4)
	// packets := packetSource.Packets()
	// ticker := time.Tick(time.Minute)

	for {

		packet, err := packetSource.NextPacket()
		if err != nil {
			log.Println(err)
			return
		}
		if packet == nil {
			return
		}
		// parser := gopacket.NewDecodingLayerParser(
		// 	layers.LayerTypeEthernet,
		// 	&ipLayer,
		// 	&tcpLayer,
		// 	&llcLayer,
		// 	&ethLayer,
		// )
		// foundLayerTypes := []gopacket.LayerType{}
		// if err := parser.DecodeLayers(packet.Data(), &foundLayerTypes); err != nil {
		// 	log.Println("Trouble decoding layers: ", err)
		// }

		// allLayers := packet.Layers()
		// for _, layer := range allLayers {
		// 	log.Printf("layer:%v\n", layer.LayerType())
		// }

		// Check for errors
		// 判断layer是否存在错误
		if err := packet.ErrorLayer(); err != nil {
			log.Println("Error decoding some part of the packet:", err)
		}

		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethernetLayer != nil {
			log.Println("Ethernet layer detected.")
			ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
			log.Println("Source MAC: ", ethernetPacket.SrcMAC)
			log.Println("Destination MAC: ", ethernetPacket.DstMAC)
			// Ethernet type is typically IPv4 but could be ARP or other
			log.Println("Ethernet type: ", ethernetPacket.EthernetType)
		}

		// Let's see if the packet is IP (even though the ether type told us)
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			log.Println("IPv4 layer detected.")
			ip, _ := ipLayer.(*layers.IPv4)
			// IP layer variables:
			// Version (Either 4 or 6)
			// IHL (IP Header Length in 32-bit words)
			// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
			// Checksum, SrcIP, DstIP
			log.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
			log.Println("Protocol: ", ip.Protocol)
		}

		// Let's see if the packet is TCP
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			log.Println("TCP layer detected.")
			tcp, _ := tcpLayer.(*layers.TCP)
			// TCP layer variables:
			// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
			// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
			log.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
			log.Println("Sequence number: ", tcp.Seq)
		}

		// When iterating through packet.Layers() above,
		// if it lists Payload layer then that is the same as
		// this applicationLayer. applicationLayer contains the payload
		decodeApplicationLayer(packet.ApplicationLayer())

		if packet.NetworkLayer() == nil ||
			packet.TransportLayer() == nil ||
			packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			// log.Println("ERR : Unknown Packet", packet)
			continue
		}

		log.Println("reading.")
	}
}
