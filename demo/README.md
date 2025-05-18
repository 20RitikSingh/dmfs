# DMFS Demo: Usage & Examples

This directory contains the demo CLI and configuration for the Distributed Mesh File System (DMFS). It provides a practical way to run, test, and interact with a DMFS cluster on your local machine.

---

## Contents
- `cmd/` — CLI commands for interacting with the DMFS cluster:
  - `upload.go`: Upload files to the mesh.
  - `read.go`: Retrieve files from the mesh.
  - `delete.go`: Remove files from the mesh.
  - `root.go`: CLI entry point, configuration, and cluster setup.
- `config.yaml` — Example configuration for running a multi-node cluster locally.

---

## How the Demo Works

### 1. Cluster Configuration
- The demo uses `config.yaml` to define multiple DMFS nodes, each with its own listening address, data directory, and bootstrap peers.
- Example:
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
- Each server is started with its own configuration, and nodes connect to each other to form a mesh.

### 2. CLI Commands
- The CLI is built using [Cobra](https://github.com/spf13/cobra) and provides the following commands:
  - `upload`: Upload a file to the mesh.
  - `read`: Download a file from the mesh using its key.
  - `delete`: Remove a file from the mesh using its key.
- Example usage:
  ```sh
  ./bin/dmfs upload --file ./myfile.txt
  ./bin/dmfs read --key <file-key> --output ./output.txt
  ./bin/dmfs delete --key <file-key>
  ```
- All commands are run against the local cluster as defined in `config.yaml`.

### 3. Running the Demo
- Build the project:
  ```sh
  make build
  ```
- Start the cluster:
  ```sh
  make run
  ```
- Use the CLI to interact with the mesh as shown above.

### 4. Customizing the Demo
- Edit `config.yaml` to add/remove nodes or change ports and data directories.
- Extend CLI functionality by adding new commands in `cmd/`.

---

## File Reference
- `cmd/root.go`: Loads configuration, initializes the cluster, and registers CLI commands.
- `cmd/upload.go`, `cmd/read.go`, `cmd/delete.go`: Implement the respective CLI operations.
- `config.yaml`: Cluster and node configuration.

---

## Troubleshooting
- Ensure all ports in `config.yaml` are available.
- Data directories must exist or be creatable by the process.
- For debugging, add `fmt.Println` or logging in CLI commands.

---

## Further Reading
- See the [main project README](../README.md) for architecture and implementation details.
- Explore code comments in each file for deeper understanding.

---

*This demo is intended for local development, testing, and learning. For production or distributed deployment, review and adapt configuration and security settings as needed.*
