package main

import "github.com/tomhjx/netcat/service"

// type rabbitStreamFactory struct{}

// type rabbitStream struct {
// 	net, transport gopacket.Flow
// 	r              tcpreader.ReaderStream
// }

// func (factory *rabbitStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
// 	stream := &rabbitStream{
// 		net:       net,
// 		transport: transport,
// 		r:         tcpreader.NewReaderStream(),
// 	}
// 	go stream.Run() // Important... we must guarantee that data from the reader stream is read.

// 	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
// 	return &stream.r
// }

// func (stream *rabbitStream) Run() {
// 	buf := bufio.NewReader(&stream.r)
// 	for {
// 		line, _, err := buf.ReadLine()
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		log.Println(line)
// 	}
// }

// var MyLayerType = gopacket.RegisterLayerType(20, gopacket.LayerTypeMetadata{Name: "MyLayerType", Decoder: gopacket.DecodeFunc(decodeMyLayer)})

// // Implement my layer
// type MyLayer struct {
// 	StrangeHeader []byte
// 	payload       []byte
// }

// func (m MyLayer) LayerType() gopacket.LayerType { return MyLayerType }
// func (m MyLayer) LayerContents() []byte         { return m.StrangeHeader }
// func (m MyLayer) LayerPayload() []byte          { return m.payload }

// // Now implement a decoder... this one strips off the first 4 bytes of the
// // packet.
// func decodeMyLayer(data []byte, p gopacket.PacketBuilder) error {
// 	// Create my layer
// 	p.AddLayer(&MyLayer{data[:4], data[4:]})
// 	// Determine how to handle the rest of the packet
// 	return p.NextDecoder(layers.LayerTypeEthernet)
// }

func main() {

	service.NewProcessor().Run()

	// var (
	// 	pcapfile string = "/work/resources/rabbit.pcap"
	// 	handle   *pcap.Handle
	// 	err      error
	// 	// Will reuse these for each packet
	// 	// ethLayer layers.Ethernet
	// 	// ipLayer  layers.IPv4
	// 	// tcpLayer layers.TCP
	// 	// llcLayer layers.LLC
	// )

	// handle, err = pcap.OpenOffline(pcapfile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer handle.Close()

	// // if err = handle.SetBPFFilter("tcp"); err != nil {
	// // 	log.Fatal("BPF filter error:", err)
	// // }

	// // streamFactory := &rabbitStreamFactory{}
	// // streamPool := tcpassembly.NewStreamPool(streamFactory)
	// // assembler := tcpassembly.NewAssembler(streamPool)

	// // packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	// packetSource := gopacket.NewPacketSource(handle, layers.IPProtocolICMPv4)
	// // packets := packetSource.Packets()
	// // ticker := time.Tick(time.Minute)

	// for {

	// 	packet, err := packetSource.NextPacket()
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	if packet == nil {
	// 		return
	// 	}
	// 	// parser := gopacket.NewDecodingLayerParser(
	// 	// 	layers.LayerTypeEthernet,
	// 	// 	&ipLayer,
	// 	// 	&tcpLayer,
	// 	// 	&llcLayer,
	// 	// 	&ethLayer,
	// 	// )
	// 	// foundLayerTypes := []gopacket.LayerType{}
	// 	// if err := parser.DecodeLayers(packet.Data(), &foundLayerTypes); err != nil {
	// 	// 	log.Println("Trouble decoding layers: ", err)
	// 	// }

	// 	// allLayers := packet.Layers()
	// 	// for _, layer := range allLayers {
	// 	// 	log.Printf("layer:%v\n", layer.LayerType())
	// 	// }

	// 	// Check for errors
	// 	// 判断layer是否存在错误
	// 	if err := packet.ErrorLayer(); err != nil {
	// 		log.Println("Error decoding some part of the packet:", err)
	// 	}

	// 	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	// 	if ethernetLayer != nil {
	// 		log.Println("Ethernet layer detected.")
	// 		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
	// 		log.Println("Source MAC: ", ethernetPacket.SrcMAC)
	// 		log.Println("Destination MAC: ", ethernetPacket.DstMAC)
	// 		// Ethernet type is typically IPv4 but could be ARP or other
	// 		log.Println("Ethernet type: ", ethernetPacket.EthernetType)
	// 	}

	// 	// Let's see if the packet is IP (even though the ether type told us)
	// 	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	// 	if ipLayer != nil {
	// 		log.Println("IPv4 layer detected.")
	// 		ip, _ := ipLayer.(*layers.IPv4)
	// 		// IP layer variables:
	// 		// Version (Either 4 or 6)
	// 		// IHL (IP Header Length in 32-bit words)
	// 		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
	// 		// Checksum, SrcIP, DstIP
	// 		log.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
	// 		log.Println("Protocol: ", ip.Protocol)
	// 	}

	// 	// Let's see if the packet is TCP
	// 	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	// 	if tcpLayer != nil {
	// 		log.Println("TCP layer detected.")
	// 		tcp, _ := tcpLayer.(*layers.TCP)
	// 		// TCP layer variables:
	// 		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
	// 		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
	// 		log.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
	// 		log.Println("Sequence number: ", tcp.Seq)
	// 	}

	// 	// When iterating through packet.Layers() above,
	// 	// if it lists Payload layer then that is the same as
	// 	// this applicationLayer. applicationLayer contains the payload
	// 	decodeApplicationLayer(packet.ApplicationLayer())

	// 	if packet.NetworkLayer() == nil ||
	// 		packet.TransportLayer() == nil ||
	// 		packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
	// 		// log.Println("ERR : Unknown Packet", packet)
	// 		continue
	// 	}

	// 	log.Println("reading.")
	// }

	// for {
	// 	select {
	// 	case packet := <-packets:
	// 		if packet == nil {
	// 			return
	// 		}

	// 		allLayers := packet.Layers()
	// 		for _, layer := range allLayers {
	// 			log.Printf("layer:%v\n", layer.LayerType())
	// 		}

	// 		// Check for errors
	// 		// 判断layer是否存在错误
	// 		if err := packet.ErrorLayer(); err != nil {
	// 			log.Println("Error decoding some part of the packet:", err)
	// 		}

	// 		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	// 		if ethernetLayer != nil {
	// 			log.Println("Ethernet layer detected.")
	// 			ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
	// 			log.Println("Source MAC: ", ethernetPacket.SrcMAC)
	// 			log.Println("Destination MAC: ", ethernetPacket.DstMAC)
	// 			// Ethernet type is typically IPv4 but could be ARP or other
	// 			log.Println("Ethernet type: ", ethernetPacket.EthernetType)
	// 		}

	// 		// Let's see if the packet is IP (even though the ether type told us)
	// 		ipLayer := packet.Layer(layers.LayerTypeIPv4)
	// 		if ipLayer != nil {
	// 			log.Println("IPv4 layer detected.")
	// 			ip, _ := ipLayer.(*layers.IPv4)
	// 			// IP layer variables:
	// 			// Version (Either 4 or 6)
	// 			// IHL (IP Header Length in 32-bit words)
	// 			// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
	// 			// Checksum, SrcIP, DstIP
	// 			log.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
	// 			log.Println("Protocol: ", ip.Protocol)
	// 		}

	// 		// Let's see if the packet is TCP
	// 		tcpLayer := packet.Layer(layers.LayerTypeTCP)
	// 		if tcpLayer != nil {
	// 			log.Println("TCP layer detected.")
	// 			tcp, _ := tcpLayer.(*layers.TCP)
	// 			// TCP layer variables:
	// 			// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
	// 			// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
	// 			log.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
	// 			log.Println("Sequence number: ", tcp.Seq)
	// 		}

	// 		// When iterating through packet.Layers() above,
	// 		// if it lists Payload layer then that is the same as
	// 		// this applicationLayer. applicationLayer contains the payload
	// 		applicationLayer := packet.ApplicationLayer()
	// 		if applicationLayer != nil {
	// 			log.Println("Application layer/Payload found.")
	// 			log.Printf("%s\n", applicationLayer.Payload())
	// 		}
	// 		return

	// 		if packet.NetworkLayer() == nil ||
	// 			packet.TransportLayer() == nil ||
	// 			packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
	// 			// log.Println("ERR : Unknown Packet", packet)
	// 			continue
	// 		}

	// 		log.Println("reading.")
	// 		// tcp := packet.TransportLayer().(*layers.TCP)
	// 		// assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

	// 		// case <-ticker:
	// 		// assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
	// 	}
	// }

}
