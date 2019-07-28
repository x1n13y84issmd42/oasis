package test

import (
	"fmt"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
)

// Runner runs a test for a single operation
type Runner struct {
	Spec api.Spec
	Log  log.ILogger
}

// Test --
func (runner Runner) Test(hostName string, operationName string, requestContentType string, responseStatus int, responseContentType string) {
	runner.Log.TestingProject(runner.Spec.GetProjectInfo())

	host := runner.Spec.GetHost(hostName)
	runner.Log.UsingHost(host)

	operation := runner.Spec.GetOperation(operationName)

	if operation != nil {
		tOp := Operation{
			runner.Log,
			host,
			operation,
		}
		tOp.Run(requestContentType, responseStatus, responseContentType)
	} else {
		fmt.Printf("The operation \"%s\" isn't there.\n", operationName)
		runner.printOperations()
	}

	fmt.Println("")
}

func (runner Runner) printOperations() {
	fmt.Println("The following operations are available:")
	for _, op := range runner.Spec.GetOperations() {
		fmt.Printf("\t%s\n", op.Name)
		fmt.Printf("\t%s @ %s\n\n", op.Method, op.Path)
	}
}
