package p2p

import "net"

// to be improved
type Payload struct {
	Key  string
	Data []byte
	From net.Addr
}
