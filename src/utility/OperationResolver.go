package utility

import (
	"regexp"
	"sort"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// OperationResolver is used to resolve operations from provided masks.
// Masks come as "prefix-*-postfix", asterisk means (.*?) in regex terms.
type OperationResolver struct {
	Spec   contract.Spec
	Logger contract.Logger
}

// NewOperationResolver creates a new OperationResolver instance.
func NewOperationResolver(spec contract.Spec, logger contract.Logger) *OperationResolver {
	return &OperationResolver{
		Spec:   spec,
		Logger: logger,
	}
}

// OpMap is a map of unique operation names.
type OpMap = map[string]contract.Operation

// Resolve returns a list of operation IDs, matching the provided mask.
func (resolver *OperationResolver) Resolve(masks []string) (ops []contract.Operation) {
	uniqueOps := OpMap{}

	for specOp := range resolver.Spec.Operations() {
		for _, mask := range masks {
			maskRX := resolver.MaskRegex(mask)
			if maskRX.Match([]byte(specOp.ID())) {
				uniqueOps[specOp.ID()] = specOp
			}
		}
	}

	// Sorting the uniqueOps keys lexicographically.
	keys := []string{}

	for oasPath := range uniqueOps {
		keys = append(keys, oasPath)
	}

	sort.Strings(keys)

	for _, key := range keys {
		ops = append(ops, uniqueOps[key])
	}

	return
}

// MaskRegex converts a string mask to a regular expression to match operations.
func (resolver *OperationResolver) MaskRegex(mask string) *regexp.Regexp {
	return regexp.MustCompile("^" + strings.ReplaceAll(mask, "*", "(.*?)") + "$")
}
