package contract

import gcontract "github.com/x1n13y84issmd42/gog/graph/contract"

// ScriptSecurity contains security values such as usernames & password and tokens.
type ScriptSecurity struct {
	Value    string `yaml:"value"`
	Token    string `yaml:"token"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// SecurityAccess contains parameter access functions
// for security values such as usernames & password and tokens.
// This is what is used in execution after script is loaded.
type SecurityAccess struct {
	Value    ParameterAccess
	Token    ParameterAccess
	Username ParameterAccess
	Password ParameterAccess
}

// Script is an interface to scenario scripts.
type Script interface {
	GetExecutionGraph() gcontract.Graph
	GetSecurity(name string) *SecurityAccess
}
