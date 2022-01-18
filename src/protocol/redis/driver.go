package redis

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"

	"github.com/tomhjx/netcat/protocol"
)

type Driver struct {
	protocol.BaseDriver
}

func NewDriver() protocol.Driver {

	return &Driver{
		BaseDriver: protocol.BaseDriver{
			SpecPort:          6379,
			SpecTransportType: "tcp",
		},
	}
}

func (me *Driver) ResolveClient(payload0 []byte) protocol.Content {
	var (
		cmd      string
		cmdCount int
	)
	cmd = ""
	r := bytes.NewReader(payload0)
	buf := bufio.NewReader(r)
	for {
		line, _, _ := buf.ReadLine()
		if len(line) == 0 {
			buff := make([]byte, 1)
			_, err := r.Read(buff)
			if err == io.EOF {
				break
			}
		}

		//Filtering useless data
		if !strings.HasPrefix(string(line), "*") {
			continue
		}

		//run
		l := string(line[1])
		cmdCount, _ = strconv.Atoi(l)
		for j := 0; j < cmdCount*2; j++ {
			c, _, _ := buf.ReadLine()
			if j&1 == 0 {
				continue
			}
			cmd += " " + string(c)
		}
		cmd += "\r\n"
	}
	// return "REDIS_REQ_CMD", cmd
	return nil
}

func (me *Driver) ResolveServer(payload0 []byte) protocol.Content {
	// return "REDIS_RESP_RES", string(payload0)
	return nil

}
