# DMFS P2P Package Documentation

This document provides a detailed overview of the `pkg/p2p` package, which implements the peer-to-peer (P2P) networking layer for the Distributed Mesh File System (DMFS).

---

## Overview
The `p2p` package is responsible for all network communication between DMFS nodes. It provides abstractions for peers, message encoding/decoding, connection management, and the transport layer (TCP-based by default).

---

## Key Components

### 1. Peer Abstraction
- **Peer Interface**: Represents a network participant. Implemented by `TCPPeer`.
- **TCPPeer**: Wraps a TCP connection, tracks if the connection is outbound, and provides methods for sending data.

### 2. Transport Layer
- **Transporter Interface**: Abstracts network transport (e.g., TCP, UDP). Defines methods for dialing, listening, consuming messages, and closing connections.
- **TCPTransporter**: Implements the transporter interface using TCP sockets. Handles listening, accepting, and dialing peers, as well as message routing.
- **TCPTransporterOptions**: Configuration for TCP transport (listen address, handshake function, decoder).

### 3. Message Handling
- **Payload**: The main message structure exchanged between peers. Contains a key, data (file chunk or control message), and sender address.
- **Encoding/Decoding**: Uses Go's `encoding/gob` for serializing/deserializing messages. Custom decoders can be plugged in.

### 4. Handshaking
- **Handshake Function**: Customizable handshake logic for authenticating or initializing new peer connections.
- **OnPeer Callback**: Invoked after a successful handshake to register or process new peers.

---

## Typical Workflow
1. **Initialization**: Each node creates a `TCPTransporter` with the desired options.
2. **Listening**: The transporter listens for incoming connections and accepts peers.
3. **Dialing**: Nodes can connect to other peers using the `Dial` method.
4. **Handshaking**: Each new connection undergoes a handshake for validation.
5. **Message Exchange**: Peers exchange `Payload` messages (file chunks, delete signals, etc.) via the transporter's channel.
6. **Peer Management**: The server tracks connected peers and can broadcast messages to all.

---

## Example Usage
```go
opts := p2p.TCPTransporterOptions{
    ListenAddr: ":8000",
    Handshake:  func(p2p.Peer) error { return nil },
    Decoder:    p2p.GOBDecoder{},
}
transporter := p2p.NewTCPTransporter(opts)
transporter.ListenAndAccept()
transporter.Dial(":8001")
for payload := range transporter.Consume() {
    // Handle incoming payloads
}
```

---

## File Reference
- `tcp_transport.go`: TCP transport implementation, peer management, connection handling.
- `message.go`: Payload structure and message definitions.
- `encoding.go`: Message encoding/decoding logic.
- `handshaker.go`: Handshake logic and interfaces.
- `transport.go`: Transporter interface and abstractions.

---

## Extending the P2P Layer
- Implement new transporters (e.g., UDP, QUIC) by conforming to the `Transporter` interface.
- Customize handshake logic for authentication or peer validation.
- Extend the `Payload` structure for new message types.

---

## Troubleshooting
- Ensure all peers use compatible encoding/decoding.
- Check for port conflicts and firewall issues.
- Use logging in the transporter and peer methods for debugging connection issues.

---

*For more details, see code comments in each file of the `pkg/p2p` directory.*
