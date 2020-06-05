package contract

// Nodes is a set of nodes.
type Nodes interface {
	Add(n Node) bool
	Remove(nID NodeID) bool
	Has(nID NodeID) bool
	Get(nID NodeID) Node
	Count() int
	Clone() Nodes
	Values() []Node
	Range() NChannel
}
