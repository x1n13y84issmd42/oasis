# Oasis CLI

Oasis has an unconventional command line interface which syntax resembles natural English naguage, or MySQL. So it configures not with "arguments" in conventional sense, as flags or switches; it's more like a query.

Some clauses are straightforward and resemble flags (`from path/to/spec`), while others may be a bit more complex, like `expect` or `log`. Both of these have "child clauses" which may be used in arbitrary order. For example, if one would like to specify a logging level 1 and plain style, they'd expressed it either as `log at level 4 in plain style` or `log in plain style at level 1`. The `log` clause itself may be anywhere on command line, just like any other top-level clause (see [Arguments](#arguments)).

#### Arguments
Argument|Example|Description
-|-|-
`from [SPECFILE]`|`from spec/petstore.yml`|Specifies the OAS spec file to use
`test [OPLIST]`|`test op1,op_two_,op_iii`|Specifies a comma-separated list of operations you want to test. Both operation IDs & names work.
`use`|See below.|Specifies how you want your requests to be configured.
`use security [NAME]`|`use security "APIKey - Header"`|Specifies which security scheme you want to use. This allows you to choose a security scheme when there are multiple defined for an operation.
`expect`|See below.|Specifies you expectations as for the operation outcome, such as response status & content type. The values from here will be used to select a proper `Response` from the OAS spec.
`expect CT [CT_NAME]`|`expect CT application/json`<br/>`expect CT "*"`|Makes Oasis choose a spec `Response` with the specified Content-Type. Asterisk means "use the first one in the spec", and is default dehavior.
`expect status [STATUS_CODE]`|`expect status 201`|Makes Oasis choose a spec `Response` with the specified response status code.
log|See below|Logging control.
`log at level [LEVEL]`|`log at level 4`|Set the log verbosity level using values 0-5.
`log in [STYLE] style`|`log in plain style`|Set the log style. `plain` means plain text log, and `festive` is a colorized version.