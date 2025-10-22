# kvStore

Simple distributed key-value store showcasing multi-node clustering and consistent hashing. It demonstrates how keys are deterministically routed to nodes in a small cluster using a hash ring. Current scope is educational/minimal.

## Features

- Multiple nodes form a logical cluster
- Consistent hashing ring (Murmur3) for key-to-node routing
- Simple in-memory KV store with RPC-based forwarding
- Minimal CLI to run nodes and interact (put/get)

Planned (not implemented yet): gossip-based membership, Raft for consensus/leadership, and data replication.

## Architecture

- Consistent hashing: `internal/cluster` builds a ring of node hashes and maps keys to the next clockwise node.
- Transport: nodes expose a Go net/rpc service; a node that receives a request forwards to the responsible peer if needed.
- Storage: `internal/node` provides a thread-safe in-memory map as the KV store.

Key files:

- `cmd/kvnode/main.go`: node process (id/addr/peers flags)
- `cmd/kvcli/cli.go`: CLI client (put/get)
- `internal/cluster/*`: consistent hash ring and node selection
- `internal/node/{node.go,handler.go}`: RPC wiring and KV handlers


## Quick start

Run three nodes in separate terminals. The peers list must include the node itself.

Terminal 1:

```zsh
go run ./cmd/kvnode --id=node-1 --addr=:8001 --peers=:8001,:8002,:8003
```

Terminal 2:

```zsh
go run ./cmd/kvnode --id=node-2 --addr=:8002 --peers=:8001,:8002,:8003
```

Terminal 3:

```zsh
go run ./cmd/kvnode --id=node-3 --addr=:8003 --peers=:8001,:8002,:8003
```

Interact using the CLI from a fourth terminal:

```zsh
# Set a value
go run ./cmd/kvcli --peers=:8001,:8002,:8003 put user:42 alice

# Get it back
go run ./cmd/kvcli --peers=:8001,:8002,:8003 get user:42
```

Under the hood, `kvcli` uses the consistent hash ring to pick the responsible node; if a command hits a non-owner, the node forwards the request to the correct peer.

## Building binaries (optional)

```zsh
go build -o bin/kvnode ./cmd/kvnode
go build -o bin/kvcli  ./cmd/kvcli
```

## Behavior and limitations

- Data is in-memory only; restarts lose data.
- No replication, durability, or quorum; single owner per key.
- No membership changes at runtime; peers are static via flags.
- No virtual nodes; each node has a single hash point.
- No retries/backoff or timeouts beyond simple RPC dial timeout.

## Extensibility roadmap

- Gossip-based membership and failure detection
- Raft for leader election, log replication, and write durability
- Read/write replication and consistency controls
- Virtual nodes for better load distribution and smoother rebalancing
- Persistence and snapshotting


