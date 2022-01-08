package service

import (
	"bytes"
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/tomhjx/netcat/amqp"
	"github.com/tomhjx/netcat/mysql"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func decodeAMQPApplicationLayer(applicationLayer gopacket.ApplicationLayer) {
	if applicationLayer == nil {
		return
	}
	log.Println("Application layer/Payload found.")
	// log.Printf("%s\n", hex.Dump(applicationLayer.Payload()))
	// log.Printf("%s\n", applicationLayer.Payload())

	r := amqp.NewReader(bytes.NewReader(applicationLayer.Payload()))
	// log.Println(len(applicationLayer.Payload()))
	// r := amqp.NewReader(strings.NewReader(string(hex.Dump(applicationLayer.Payload()))))
	frame, err := r.ReadFrame()
	if err != nil {
		log.Printf("frame err, read: %s", err)
		return
	}
	log.Println(frame)
}

func handleAMQP() {
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
		decodeAMQPApplicationLayer(packet.ApplicationLayer())
		return

		if packet.NetworkLayer() == nil ||
			packet.TransportLayer() == nil ||
			packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			// log.Println("ERR : Unknown Packet", packet)
			continue
		}

		log.Println("reading.")
	}
}

type ProtocolStreamFactory struct {
}

type ProtocolStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (m *ProtocolStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {

	//init stream struct
	stm := &ProtocolStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}

	//new stream
	fmt.Println("# Start new stream:", net, transport)

	//decode packet
	mysql.NewInstance().ResolveStream(net, transport, &(stm.r))

	return &(stm.r)
}

func handleMysql() {
	var (
		pcapfile string = "/work/resources/mysql.pcap"
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
	//capture
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// set up assembly
	// streamFactory := &ProtocolStreamFactory{}
	// streamPool := tcpassembly.NewStreamPool(streamFactory)
	// assembler := tcpassembly.NewAssembler(streamPool)

	// loop until ctrl+z

	for {

		packet, err := packetSource.NextPacket()
		if err != nil {
			log.Println("Read Packet ERR:", err)
			return
		}
		if packet == nil {
			return
		}
		if packet.NetworkLayer() == nil ||
			packet.TransportLayer() == nil ||
			packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			fmt.Println("ERR : Unknown Packet -_-")
			continue
		}

		applicationLayer := packet.ApplicationLayer()
		if applicationLayer == nil {
			continue
		}

		// tcp := packet.TransportLayer().(*layers.TCP)
		// assembler.AssembleWithTimestamp(
		// 	packet.NetworkLayer().NetworkFlow(),
		// 	tcp, packet.Metadata().Timestamp,
		// )
		m := mysql.NewInstance()
		if m == nil {
			log.Fatal("mysql instance error.")
		}
		netflow := packet.NetworkLayer().NetworkFlow()

		m.ResolveLocal(netflow, bytes.NewReader(applicationLayer.Payload()))
	}
}

func (proc *Processor) Run() {
	// handleAMQP()
	handleMysql()

}
