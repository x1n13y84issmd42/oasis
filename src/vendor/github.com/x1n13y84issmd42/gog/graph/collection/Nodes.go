package collection

import (
	"sort"

	"github.com/x1n13y84issmd42/gog/graph/contract"
)

// Nodes is a set of INode instances.
type Nodes struct {
	Map map[contract.NodeID]contract.Node
}

// NewNodes creates a new node set instance.
func NewNodes(nodes ...contract.Node) *Nodes {
	res := &Nodes{
		Map: make(map[contract.NodeID]contract.Node),
	}

	for _, n := range nodes {
		res.Add(n)
	}

	return res
}

// Add adds a node to the set and returns true if it had been inserted for the first time.
func (set *Nodes) Add(n contract.Node) bool {
	had := set.Has(n.ID())
	set.Map[n.ID()] = n
	return had
}

// Remove removes a node from the set. Returns true if node was present in the set before removal.
func (set *Nodes) Remove(nID contract.NodeID) bool {
	had := set.Has(nID)
	delete(set.Map, nID)
	return had
}

// Has tells whether a node is present in the set.
func (set *Nodes) Has(nID contract.NodeID) bool {
	_, ok := set.Map[nID]
	return ok
}

// Get returns a node from the set.
func (set *Nodes) Get(nID contract.NodeID) contract.Node {
	return set.Map[nID]
}

// Count tells how many nodes are currently in the set.
func (set Nodes) Count() int {
	return len(set.Map)
}

// Clone creates a new set by copying the receiver set.
func (set Nodes) Clone() contract.Nodes {
	res := NewNodes()
	for _, n := range set.Map {
		res.Add(n)
	}
	return res
}

// Values creates a slice of values taken from the set.
func (set Nodes) Values() []contract.Node {
	res := []contract.Node{}
	for n := range set.Range() {
		res = append(res, n)
	}
	return res
}

// Range iterates over elements in lexicographic key order.
func (set Nodes) Range() contract.NChannel {
	ch := make(chan contract.Node)
	go func() {
		// Collecting keys.
		keys := []string{}
		for nID := range set.Map {
			keys = append(keys, string(nID))
		}
		// Sorting them.
		sort.Strings(keys)

		// Sending values to the channel.
		for _, k := range keys {
			ch <- set.Map[contract.NodeID(k)]
		}

		close(ch)
	}()
	return ch
}
