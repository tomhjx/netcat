package amqp

type Content struct {
	Heartbeat *heartbeatFrame
	Method    *methodFrame
	Header    *headerFrame
	Body      string
}

func (me Content) Section() string {
	return "amqp.content"
}
