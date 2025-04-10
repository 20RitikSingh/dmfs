package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":8080"
	tr := NewTCPTransporter(TCPTransporterOptions{
		ListenAddr: listenAddr,
	})

	if tr == nil {
		t.Fatalf("NewTCPTransport(%s) returned nil", listenAddr)
	}

	assert.Equal(t, tr.ListenAddr, listenAddr)
	assert.Nil(t, tr.ListenAndAccept())
}
