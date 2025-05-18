# DMFS: Distributed Mesh File System

## Overview
DMFS is a distributed mesh file system written in Go. It enables decentralized file storage and retrieval across a network of peer nodes, supporting encryption, chunking, and peer-to-peer (P2P) communication. The system is designed for reliability, scalability, and security.

## Architecture & Design

### Core Components
- **Server**: Each node runs a file server (`server.FileServer`) that manages local storage and communicates with peers.
- **P2P Layer**: Handles peer discovery, message passing, and network bootstrapping using TCP connections (`pkg/p2p`).
- **Chunking & Encryption**: Files are split into chunks and encrypted before storage/transmission (`pkg/chunk`, `pkg/encryption`).
- **Store**: Abstracts file storage and retrieval on disk (`store/`).
- **CLI**: User interacts via a command-line interface built with Cobra (`demo/cmd`).

### Data Flow
1. **Startup**: Each server reads its configuration from `demo/config.yaml` and joins the mesh network, optionally bootstrapping to known peers.
2. **File Operations**:
   - **Upload**: Files are chunked, encrypted, and distributed. Metadata (key, chunk info) is shared with peers.
   - **Read**: Chunks are requested from the network, decrypted, and reassembled.
   - **Delete**: Deletion signals are broadcast to all peers.
3. **P2P Communication**: Uses TCP for peer connections, with handshake and message encoding/decoding.
4. **Replication**: When a file is uploaded, it is broadcast to all connected peers for redundancy.

### Encryption & Chunking
- **Encryption**: AES-GCM (256-bit) is used for chunk encryption. The key must be 32 bytes.
- **Chunking**: Files are split into fixed-size chunks for efficient transfer and storage.
- **Hashing**: Blake3 is used to generate unique file keys.

## Configuration
Edit `demo/config.yaml` to define the cluster:
```yaml
servers:
  - listenAddr: "127.0.0.1:8000"
    root: "./data/server0"
    nodes: []
  - listenAddr: "127.0.0.1:8001"
    root: "./data/server1"
    nodes:
      - ":8000"
  - listenAddr: "127.0.0.1:8002"
    root: "./data/server2"
    nodes:
      - ":8000"
      - ":8001"
```

## Building & Running

### Prerequisites
- Go 1.18+

### Build
```pwsh
make build
```

### Run
```pwsh
make run
```

### Test
```pwsh
make test
```

### Clean
```pwsh
make clean
```

## Usage
- Use the CLI (`dmfs`) to upload, read, and delete files.
- The CLI commands are defined in `demo/cmd/upload.go`, `read.go`, and `delete.go`.

## Demo & CLI Reference
For detailed instructions on running the demo cluster, using the CLI, and customizing your local setup, see [`demo/README.md`](demo/README.md).

## Key Files & Directories
- `main.go`: Entry point.
- `server/`: File server logic.
- `pkg/p2p/`: Peer-to-peer networking ([see detailed documentation](pkg/p2p/README.md)).
- `pkg/chunk/`: File chunking logic.
- `pkg/encryption/`: Encryption and hashing ([see cipher choices](pkg/encryption/cipherCoices.md)).
- `store/`: Local file storage abstraction.
- `demo/`: CLI and configuration.

## How It Works
- Each server starts, connects to peers, and listens for file operations.
- File operations are broadcast to all peers for redundancy.
- Files are encrypted and chunked before storage.
- The system is resilient to node failures as files are replicated.

## Extending DMFS
- Add new commands to the CLI by editing `demo/cmd/`.
- Change chunk size or encryption by modifying `pkg/chunk` and `pkg/encryption`.
- Add new peer discovery mechanisms in `pkg/p2p`.

## License
MIT License.

---
*For more details, see code comments and each package's documentation.*
