package node

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/mildlybrutal/kvStore/internal/cluster"
)

type Node struct {
	ID      string
	Address string
	kvStore *KVStore
	cluster *cluster.Cluster
	rpcSrv  *rpc.Server
}

func NewNode(id, address string, cl *cluster.Cluster) *Node {
	return &Node{
		ID:      id,
		Address: address,
		kvStore: NewKVStore(),
		cluster: cl,
		rpcSrv:  rpc.NewServer(),
	}
}

func (n *Node) Serve() {
	if err := n.rpcSrv.Register(n.kvStore); err != nil {
		log.Fatalf("[%s] failed to register KVStore RPC: %v", n.ID, err)
	}

	listener, err := net.Listen("tcp", n.Address)
	if err != nil {
		log.Fatalf("[%s] error listening on %s: %v", n.ID, n.Address, err)
	}
	log.Printf("[%s] listening on %s", n.ID, n.Address)

	n.rpcSrv.Accept(listener)
}

func dialRPC(addr string) (*rpc.Client, error) {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(conn), nil
}

func (n *Node) Put(key, value string) error {
	targetNode := n.cluster.GetNodeForKey(key)
	if targetNode.Addr == n.Address {
		args := map[string]string{key: value}
		var reply bool
		return n.kvStore.Put(args, &reply)
	}

	client, err := dialRPC(targetNode.Addr)

	if err != nil {
		return fmt.Errorf("[%s] failed to connect to %s: %v", n.ID, targetNode.Addr, err)
	}
	defer client.Close()

	args := map[string]string{key: value}
	var reply bool
	if err := client.Call("KVStore.Put", args, &reply); err != nil {
		return fmt.Errorf("[%s] forwarding PUT('%s') → [%s] failed: %v",
			n.ID, key, targetNode.Addr, err)
	}

	log.Printf("[%s] forwarded PUT('%s') → [%s]", n.ID, key, targetNode.Addr)
	return nil
}

func (n *Node) Get(key string) (string, error) {
	targetNode := n.cluster.GetNodeForKey(key)

	if targetNode.Addr == n.Address {
		var val string
		if err := n.kvStore.Get(key, &val); err != nil {
			return "", err
		}
		return val, nil
	}

	client, err := dialRPC(targetNode.Addr)
	if err != nil {
		return "", fmt.Errorf("[%s] failed to connect to %s: %v", n.ID, targetNode.Addr, err)
	}
	defer client.Close()

	var val string
	if err := client.Call("KVStore.Get", key, &val); err != nil {
		return "", fmt.Errorf("[%s] forwarding GET('%s') → [%s] failed: %v",
			n.ID, key, targetNode.Addr, err)
	}

	log.Printf("[%s] forwarded GET('%s') → [%s]", n.ID, key, targetNode.Addr)
	return val, nil
}
