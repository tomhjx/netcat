package amqp

import (
	"bytes"
	"log"

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

// not work

func (me *Driver) ResolveClient(payload []byte) (string, string) {
	if len(payload) == 0 {
		return "", ""
	}
	r := NewReader(bytes.NewReader(payload))
	frame, err := r.ReadFrame()
	if err != nil {
		log.Printf("frame err, read: %s", err)
		return "", ""
	}

	var (
		m         message
		header    *headerFrame
		remaining int
		body      []byte
	)

	switch f := frame.(type) {
	case *heartbeatFrame:
		// drop

	case *headerFrame:
		// start content state
		header = f
		remaining := int(header.Size)
		if remaining == 0 {
			m.(messageWithContent).setContent(header.Properties, nil)
			// return m
			return "RABBIT_HEADER", ""
		}

	case *bodyFrame:
		// continue until terminated
		body = append(body, f.Body...)
		// log.Println("%s", body)
		remaining -= len(f.Body)
		if remaining <= 0 {
			m.(messageWithContent).setContent(header.Properties, body)
			// return m
			return "RABBIT_BODY", ""
		}

	case *methodFrame:
		log.Printf("%s", payload)
		// if reflect.TypeOf(m) == reflect.TypeOf(f.Method) {
		// wantv := reflect.ValueOf(m).Elem()
		// havev := reflect.ValueOf(f.Method).Elem()
		// wantv.Set(havev)
		// if _, ok := m.(messageWithContent); !ok {
		// 	// return m
		// 	return "RABBIT_METHOD", ""
		// }
		// } else {
		// log.Printf("%#v", f)
		// log.Printf("expected method type: %T, got: %T", m, f.Method)
		// }
		me.resolveMethodFrame(f.Method)

	default:
		log.Printf("unexpected frame: %+v", f)
	}

	return "", ""
}

func (me *Driver) resolveMethodFrame(m0 message) (string, string) {
	switch m := m0.(type) {
	case *basicPublish:
		log.Printf("basicPublish: %v", m)
		// hds, ctx := m.getContent()
		// log.Printf("msg id: %s", hds.MessageId)
		// log.Printf("body : %s", ctx)

	default:
		log.Printf("unexpected method type: %T", m0)
	}
	return "", ""
}

func (me *Driver) ResolveServer(payload0 []byte) (string, string) {
	return "", ""
}
