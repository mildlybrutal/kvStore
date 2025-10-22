package client

import (
	"fmt"
	"net/rpc"
)

type Client struct {
	nodeAddress string
	rpcClient   *rpc.Client
}

func NewClient(nodeAddress string) (*Client, error) {
	client, err := rpc.DialHTTP("tcp", nodeAddress)

	if err != nil {
		return nil, fmt.Errorf("error dialing RPC server at %s: %v", nodeAddress, err)
	}

	return &Client{
		nodeAddress: nodeAddress,
		rpcClient:   client,
	}, nil
}

// calls the remote Get method on the node.
func (c *Client) Get(key string) (string, error) {
	var reply string
	err := c.rpcClient.Call("KVStore.Get", key, &reply)

	if err != nil {
		return "", fmt.Errorf("error calling Get for key '%s': %v", key, err)
	}
	return reply, nil
}

// calls the remote Puts method on the node.
func (c *Client) Put(key, value string) error {
	args := map[string]string{key: value}
	var reply bool

	err := c.rpcClient.Call("KVStore.Put", args, &reply)
	if err != nil {
		return fmt.Errorf("error calling Put for key '%s': %v", key, err)
	}
	if !reply {
		return fmt.Errorf("put operation failed for key '%s'", key)
	}
	return nil
}

func (c *Client) Close() error {
	return c.rpcClient.Close()
}
