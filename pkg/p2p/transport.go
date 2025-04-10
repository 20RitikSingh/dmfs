package p2p

import "net"

type Peer interface {
	net.Conn
	Send([]byte) error
}
type Transporter interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan Payload
	Close() error
}
