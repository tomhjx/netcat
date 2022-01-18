package core

import (
	"log"
	"sync"

	"github.com/google/gopacket"
	"github.com/tomhjx/netcat/protocol"
	amqpProtocol "github.com/tomhjx/netcat/protocol/amqp"
	mysqlProtocol "github.com/tomhjx/netcat/protocol/mysql"
	redisProtocol "github.com/tomhjx/netcat/protocol/redis"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func tryMysql() (string, protocol.Driver) {
	return "/work/resources/mysql.pcap", mysqlProtocol.NewDriver()
}

func tryRabbit() (string, protocol.Driver) {
	return "/work/resources/rabbit.pcap", amqpProtocol.NewDriver()
}

func tryRedis() (string, protocol.Driver) {
	return "/work/resources/redis.pcap", redisProtocol.NewDriver()
}

func (proc *Processor) Run() {

	// pcapfile, pd := tryMysql()
	pcapfile, pd := tryRabbit()
	// pcapfile, pd := tryRedis()
	concurrency := 3
	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	sources := make(chan gopacket.Packet, 100)
	resolveds := make(chan *Resolved, 10)
	parser := NewParser(pd)
	ier := NewInputer(pd)

	// input
	go func() {
		defer wg.Done()
		ier.RegisterReadTrigger(func(s gopacket.Packet) {
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

				resolved := parser.Resolve(source)
				if resolved != nil {
					resolveds <- parser.Resolve(source)
				}

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
