package api

import "github.com/x1n13y84issmd42/oasis/src/contract"

// OperationCache proxies the GetOperation method
// and caches returned operations.
type OperationCache struct {
	spec  contract.Spec
	stash map[string]contract.Operation
}

// NewOperationCache creates a new OperationCache instance.
func NewOperationCache(spec contract.Spec) OperationCache {
	return OperationCache{
		spec:  spec,
		stash: make(map[string]contract.Operation),
	}
}

// GetOperation returns a previously cache operation instance, when present,
// otherwise gets it from the spec, caches and returns.
func (cache OperationCache) GetOperation(id string) contract.Operation {
	if op, ok := cache.stash[id]; ok {
		return op
	}

	cache.stash[id] = cache.spec.GetOperation(id)
	return cache.stash[id]
}
