package node

import (
	"fmt"
	"sync"
)

// reps the local key-value storage on a node
type KVStore struct {
	mu    sync.RWMutex
	store map[string]string
}

// new instance of KVStore
func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

// Retrieves a value for a given key
func (kv *KVStore) Get(key string, reply *string) error {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	if val, ok := kv.store[key]; ok {
		*reply = val
		return nil
	}
	return fmt.Errorf("key '%s' not found", key)
}

// stores a key-value pair

func (kv *KVStore) Put(pair map[string]string, reply *bool) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	for key, value := range pair {
		kv.store[key] = value
	}
	*reply = true
	return nil
}
