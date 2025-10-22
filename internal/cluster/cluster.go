// clusters lets every node know other nodes’ addresses and contact them.
package cluster

import "fmt"

type Cluster struct {
	ring *ConsistentHashingRing
	self *NodeInfo
}

func NewCluster(selfAddr string, peerAddrs []string) *Cluster {
	chr := NewConsistentHashingRing()
	for i, addr := range append(peerAddrs, selfAddr) {
		chr.AddNode(fmt.Sprintf("node-%d", i+1), addr)
	}

	return &Cluster{
		ring: chr,
		self: &NodeInfo{Addr: selfAddr},
	}
}

func (c *Cluster) GetNodeForKey(key string) *NodeInfo {
	node, _ := c.ring.FindNearbyNodeIndex(key)
	return node
}
