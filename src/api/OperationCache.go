package api

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/strings"
)

// OperationCache proxies the GetOperation method
// and caches returned operations.
type OperationCache struct {
	specs map[string]contract.OperationAccess
	stash map[string]contract.Operation
}

// NewOperationCache creates a new OperationCache instance.
func NewOperationCache(specs map[string]contract.OperationAccess) OperationCache {
	return OperationCache{
		specs: specs,
		stash: make(map[string]contract.Operation),
	}
}

// GetOperation returns a previously cache operation instance, when present,
// otherwise gets it from the spec, caches and returns.
func (cache OperationCache) GetOperation(id string) contract.Operation {
	if op, ok := cache.stash[id]; ok {
		return op
	}

	parsedID := strings.Split(id, ".")
	specID := parsedID[0]
	opID := strings.Cut(id, len(specID+"."), len(id))

	cache.stash[id] = cache.specs[specID].GetOperation(opID)

	return cache.stash[id]
}
