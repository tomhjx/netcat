package protocol

type Driver interface {
	ResolveClient(payload []byte) Content
	ResolveServer(payload []byte) Content
	Port() int
	TransportType() string
}

type BaseDriver struct {
	SpecPort          int
	SpecTransportType string
}

func (me *BaseDriver) Port() int {
	return me.SpecPort
}
func (me *BaseDriver) TransportType() string {
	return me.SpecTransportType
}

type Content interface {
	Section() string
}
