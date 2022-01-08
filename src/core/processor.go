package core

import (
	"log"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
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
			log.Println("ERR : Unknown Packet -_-")
			continue
		}

		applicationLayer := packet.ApplicationLayer()
		if applicationLayer == nil {
			continue
		}

		log.Println("...")

		// tcp := packet.TransportLayer().(*layers.TCP)
		// assembler.AssembleWithTimestamp(
		// 	packet.NetworkLayer().NetworkFlow(),
		// 	tcp, packet.Metadata().Timestamp,
		// )
		// m := mysql.NewInstance()
		// if m == nil {
		// 	log.Fatal("mysql instance error.")
		// }
		// netflow := packet.NetworkLayer().NetworkFlow()

		// m.ResolveLocal(netflow, bytes.NewReader(applicationLayer.Payload()))
	}
}

func (proc *Processor) Run() {
	pcapfile := "/work/resources/mysql2.pcap"
	concurrency := 3
	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	sources := make(chan *Source, 100)
	resolveds := make(chan *Resolved, 10)

	// input
	go func() {
		defer wg.Done()
		ier := NewInputer()
		ier.RegisterReadTrigger(func(s *Source) {
			// log.Println("send source......")
			sources <- s
			// log.Println("sent source.")
		})
		// ier.RegisterReadDoneTrigger(func() {
		// 	close(sources)
		// 	log.Println("close sources chan")
		// })
		ier.ReadOffline(pcapfile)
	}()

	// parse
	go func() {
		defer wg.Done()
		for {
			select {
			case source := <-sources:
				// time.Sleep(1 * time.Second)
				// log.Println("parse: ", len(source.payload))
				parser := NewParser()
				resolveds <- parser.Resolve(source)
			}
		}
	}()

	// output
	go func() {
		defer wg.Done()
		for {
			select {
			case resolved := <-resolveds:
				log.Println("output: ", resolved)
			}
		}
	}()

	// handleMysql()

	wg.Wait()

}
