package wsock

import "net"

type MsgShape struct {
	msg string
}

type WebSocketChatWorker struct {
	conn  net.Conn
	msgCH chan MsgShape
	Idx   uint64
}
