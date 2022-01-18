package amqp

import (
	"bytes"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/tomhjx/netcat/protocol"
)

type Driver struct {
	protocol.BaseDriver
}

func NewDriver() protocol.Driver {

	return &Driver{
		BaseDriver: protocol.BaseDriver{
			SpecPort:          5672,
			SpecTransportType: "tcp",
		},
	}
}

func (me *Driver) ResolveClient(payload []byte) protocol.Content {
	return me.resolve(payload)
}

func (me *Driver) resolve(payload []byte) protocol.Content {
	if len(payload) == 0 {
		return nil
	}
	r := NewReader(bytes.NewReader(payload))

	var (
		remaining int
		body      []byte
	)

	c := &Content{}

	for {

		frame, err := r.ReadFrame()
		if err == io.EOF {
			log.Debug("frame eof.")
			return c
		} else if err != nil {
			log.Errorf("frame err, read: %s", err)
			continue
		}

		switch f := frame.(type) {
		case *heartbeatFrame:
			c.Heartbeat = f

		case *headerFrame:
			log.Debug("headerFrame")
			// start content state
			c.Header = f
			remaining = int(f.Size)

		case *bodyFrame:
			log.Debug("bodyFrame")
			// continue until terminated
			body = append(body, f.Body...)
			remaining -= len(f.Body)
			if remaining > 0 {
				continue
			}
			c.Body = string(body)

		case *methodFrame:
			log.Debug("methodFrame")
			c.Method = f
		default:
			log.Errorf("unexpected frame: %+v", f)
		}
	}
	return c
}

func (me *Driver) ResolveServer(payload []byte) protocol.Content {
	return me.resolve(payload)
}
