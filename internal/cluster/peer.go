// Represents known peers and their state.
package cluster

import (
	"net/rpc"
)

type Peer struct {
	ID      string
	Address string
}

func (p *Peer) SendPut(key, value string) error {
	client, err := rpc.DialHTTP("tcp", p.Address)
	if err != nil {
		return err
	}

	defer client.Close()

	args := map[string]string{key: value}
	var reply bool
	return client.Call("KVStore.Put", args, &reply)
}
