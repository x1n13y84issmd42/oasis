package test

import (
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Runner runs a test for a single operation
type Runner struct {
	Spec api.Spec
	Log  log.ILogger
	Host *api.Host
}

// UseHost selects a host to use for test.
// Empty hostName means using a first one in the list as default.
func (runner *Runner) UseHost(hostName string) bool {
	//	Figuring out the host name.
	//	Empty string means default host.
	var host *api.Host

	if hostName != "" {
		host, _ = runner.Spec.GetHost(hostName)
	} else {
		runner.Log.UsingDefaultHost()
		host, _ = runner.Spec.GetDefaultHost()
	}

	if host != nil {
		runner.Log.UsingHost(host)
	} else {
		runner.Log.HostNotFound(hostName)
		return false
	}

	runner.Host = host

	return true
}

// Test --
func (runner *Runner) Test(operationName string, requestContentType string, responseStatus int, responseContentType string) bool {
	/* //	Getting the tested operation from the spec.
	operation := runner.Spec.GetOperation(operationName)
	result := false

	if operation != nil {
		//	Operation test.
		tOp := Operation{
			runner.Log,
			runner.Host,
			operation,
		}
		//	Running the op test.
		result = tOp.Run(requestContentType, responseStatus, responseContentType)
	} else {
		runner.Log.OperationNotFound(operationName)
		runner.printOperations()
	}

	fmt.Println("") */
	return false
}

func (runner *Runner) printOperations() {
	/* fmt.Println("The following operations are available:")
	for _, op := range runner.Spec.GetOperations() {
		if op.ID != "" {
			fmt.Printf("\t%s [%s]\n", op.Name, op.ID)
		} else {
			fmt.Printf("\t%s\n", op.Name)
		}
		fmt.Printf("\t%s @ %s\n\n", op.Method, op.Path.Path)
	} */
}
