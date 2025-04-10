package p2p

import (
	"errors"
	"log"
	"net"
	"reflect"
)

type TCPPeer struct {
	// connection to the peer
	net.Conn
	// is conn outgoing i.e we dial to it
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
	}
}

func (peer *TCPPeer) Send(data []byte) error {
	_, err := peer.Write(data)
	return err
}

type TCPTransporterOptions struct {
	ListenAddr string
	Handshake  Handshake
	Decoder    Decoder
}

type TCPTransporter struct {
	listener net.Listener
	TCPTransporterOptions

	payloadch chan Payload
	// to be run after handshake (init)
	OnPeer func(Peer) error
}

func OnPeer(peer Peer) error {
	log.Printf("verify peer logic\n")
	return nil
}

func NewTCPTransporter(opts TCPTransporterOptions) *TCPTransporter {
	return &TCPTransporter{
		TCPTransporterOptions: opts,
		payloadch:             make(chan Payload),
		OnPeer:                OnPeer,
	}
}

func (t *TCPTransporter) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConnection(conn, true)
	return nil
}

func (t *TCPTransporter) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	// accept connections
	go func() {
		for {
			conn, err := t.listener.Accept()
			if errors.Is(err, net.ErrClosed) {
				return
			}
			if err != nil {
				log.Printf("TCP Transport: Accept error: %v", err)
				continue
			}

			go t.handleConnection(conn, false)
		}
	}()
	return nil
}

func (t *TCPTransporter) handleConnection(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		log.Printf("closing connection %s", err)
	}()

	peer := NewTCPPeer(conn, outbound)
	if outbound {
		log.Printf("new outgoing connection: %+v \n", peer)
	} else {
		log.Printf("new incommig connection: %+v \n", peer)
	}

	//handshaking
	if err = t.Handshake(peer); err != nil {
		log.Printf("handshake error: %s", err)
		return
	} else {
		log.Println("handshake successful!!")
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			log.Println("error OnPeer ", err)
			return
		}
	}
	//read loop
	if t.payloadch == nil {
		log.Println("Error: RPC channel is not initialized")
		return
	}
	for {
		payload := Payload{}
		if err = t.Decoder.Decode(conn, &payload); err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
				return
			}
			log.Printf("TCP error %s\n", err)
			continue
		}

		// log.Printf("payload incoming: %+v", payload)
		t.payloadch <- payload
	}

}

func (t *TCPTransporter) Consume() <-chan Payload {
	return t.payloadch
}

func (t *TCPTransporter) Close() error {
	return t.listener.Close()
}
