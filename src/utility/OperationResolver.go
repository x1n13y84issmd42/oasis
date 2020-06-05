package utility

import (
	"regexp"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// OperationResolver is used to resolve operations from provided masks.
// Masks come as "prefix-*-postfix", asterisk means (.*?) in regex terms.
type OperationResolver struct {
	Spec   api.Spec
	Logger log.ILogger
}

// NewOperationResolver creates a new OperationResolver instance.
func NewOperationResolver(spec api.Spec, logger log.ILogger) *OperationResolver {
	return &OperationResolver{
		Spec:   spec,
		Logger: logger,
	}
}

// OpMap is a map of unique operation names.
type OpMap = map[string]bool

// Resolve returns a list of operation IDs, matching the provided mask.
func (resolver *OperationResolver) Resolve(masks []string, params *api.OperationParameters) (ops []*api.Operation) {
	specOps := resolver.Spec.ListOperations()
	uniqueOps := OpMap{}

	for _, mask := range masks {
		resolver.ResolveMask(mask, specOps, &uniqueOps)
	}

	for opID := range uniqueOps {
		specOp, specOpErr := resolver.Spec.GetOperation(opID, params)

		if specOpErr == nil {
			ops = append(ops, specOp)
		} else {
			resolver.Logger.Error(specOpErr)
		}
	}

	return
}

// ResolveMask matches a single mask to the provided list of operation IDs.
func (resolver *OperationResolver) ResolveMask(mask string, specOps *[]string, uniqueOps *OpMap) {
	maskRX := resolver.MaskRegex(mask)

	for _, specOpID := range *specOps {
		if maskRX.Match([]byte(specOpID)) {
			(*uniqueOps)[specOpID] = true
		}
	}
}

// MaskRegex converts a string mask to a regular expression to match operations.
func (resolver *OperationResolver) MaskRegex(mask string) *regexp.Regexp {
	return regexp.MustCompile(strings.ReplaceAll(mask, "*", "(.*?)"))
}
