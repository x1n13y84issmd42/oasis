package contract

import gcontract "github.com/x1n13y84issmd42/gog/graph/contract"

// Script is an interface to scenario scripts.
type Script interface {
	GetGraph() gcontract.Graph
}
