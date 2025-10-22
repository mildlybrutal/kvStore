// I need ring structure - a sorted list (or map) of hash points to node identifiers
// hashing func - murmur3
// I need to add nodes to the ring

package cluster

import (
	"errors"
	"slices"
	"sync"

	"github.com/spaolacci/murmur3"
)

type NodeInfo struct {
	ID   string
	Addr string
}

type ConsistentHashingRing struct {
	mu     sync.RWMutex
	nodes  map[uint32]*NodeInfo
	hashes []uint32
}

func NewConsistentHashingRing() *ConsistentHashingRing {
	return &ConsistentHashingRing{
		nodes:  make(map[uint32]*NodeInfo),
		hashes: []uint32{},
	}
}

func hashFunction(key string) uint32 {
	return murmur3.Sum32([]byte(key))
}

func (chr *ConsistentHashingRing) AddNode(id string, addr string) {
	chr.mu.Lock()
	defer chr.mu.Unlock()

	hash := hashFunction(id)

	chr.nodes[hash] = &NodeInfo{
		ID:   id,
		Addr: addr,
	}

	chr.hashes = append(chr.hashes, hash)
	slices.Sort(chr.hashes)
}

func (chr *ConsistentHashingRing) RemoveNode(id string) {
	chr.mu.Lock()
	defer chr.mu.Unlock()

	hash := hashFunction(id)
	delete(chr.nodes, hash)

	for i, h := range chr.hashes {
		if h == hash {
			chr.hashes = slices.Delete(chr.hashes, i, i+1)
			break
		}
	}
}

func (chr *ConsistentHashingRing) FindNearbyNodeIndex(key string) (*NodeInfo, error) {
	chr.mu.RLock()
	defer chr.mu.RUnlock()

	if len(chr.hashes) == 0 {
		return nil, errors.New("no nodes in the ring")
	}

	h := hashFunction(key)

	for _, nodeHash := range chr.hashes {
		if nodeHash > h {
			return chr.nodes[nodeHash], nil
		}
	}

	return chr.nodes[chr.hashes[0]], nil
}
