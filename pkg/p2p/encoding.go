package p2p

import (
	"encoding/gob"
	"io"
	"log"
)

type Decoder interface {
	Decode(io.Reader, *Payload) error
}

type GOBDecoder struct{}

func GOBEncode(w io.Writer, payload *Payload) error {
	// log.Println("encoding payload ", payload)

	return gob.NewEncoder(w).Encode(payload)
}

func (gdec GOBDecoder) Decode(r io.Reader, v *Payload) error {
	return gob.NewDecoder(r).Decode(v)
}

type DefaultDecoder struct{}

func (dDec DefaultDecoder) Decode(r io.Reader, msg *Payload) error {
	buff := make([]byte, 1024)
	n, err := r.Read(buff)
	if err != nil {
		log.Printf("error decoding: %s", err)
		return err
	}
	msg.Data = buff[:n]
	return nil
}
