# OASIS
Open API Specification Intelligence Services. Or in a less fancy way, a tool to test APIs which uses OAS/Swagger spec files as a test suite.

[![Build Status](https://travis-ci.com/x1n13y84issmd42/oasis.svg?branch=master)](https://travis-ci.com/x1n13y84issmd42/oasis)
<a href="https://codeclimate.com/github/x1n13y84issmd42/oasis/maintainability"><img src="https://api.codeclimate.com/v1/badges/a6348253063e179ba44f/maintainability" /></a>
<a href="https://codeclimate.com/github/x1n13y84issmd42/oasis/test_coverage"><img src="https://api.codeclimate.com/v1/badges/a6348253063e179ba44f/test_coverage" /></a>
[![Go Report Card](https://goreportcard.com/badge/github.com/x1n13y84issmd42/oasis)](https://goreportcard.com/report/github.com/x1n13y84issmd42/oasis)

Work in progress.

- [Usage](#usage)
  - [Manual mode](#manual-mode)
  - [Script mode](#script-mode)

## Usage
Oasis can be used either in manual mode, which allows to test single operations, or in script mode, which is designed to test complex interaction scenarios, involving multiple endpoints and reusing data across them.

### ☑️ Manual mode
`run/oasis from spec/petstore.yaml test findPetsByStatus`

Oasis uses the example value for the `status` query parameter for the `findPetsByStatus` operation defined in the spec file.

You can override any parameter from CLI:

`run/oasis from spec/petstore.yaml test getPetById use path parameters petId=10`

Increase logging verbosity to see how parameters are used.

`run/oasis from spec/petstore.yaml test getPetById log at level 6 use path parameters petId=10`

For an example of error reporting use the malformed spec file which doesn't correspond to the actual API responses:

`run/oasis from spec/errors.yaml test taskList log at level 6`

`run/oasis from spec/errors.yaml test configTypes log at level 6`

📖 [Learn more about CLI](doc/CLI.md)

📖 [Learn more about operation parameters](doc/Parameters.md)

### 🔀 Script mode
For complex scenarios involving multiple endpoints and data reuse across them there is a script mode:

`run/oasis execute script/petstore.yaml`
`run/oasis execute script/nuxeo.yaml`

Script is a graph of dependent operations. Cycles are not allowed:

`run/oasis execute script/cycle.yaml`

📖 [Learn more about scripts](doc/Script.md)

## Resources
[OpenAPI Spec](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#securitySchemeObject)

[HTTP Auth Registry](http://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml)

[Stoplight OAS Editor](https://stoplight.io/p/studio/sl/_/6c5faofy)

[Reyesoft API Playground](https://jsonapiplayground.reyesoft.com/)

[Nuxeo API Playground](https://nuxeo.github.io/api-playground/#/resources)

[Booker API Playground](https://restful-booker.herokuapp.com/apidoc/index.html)

[Oasis V2](https://github.com/x1n13y84issmd42/oasis/tree/c88c9a15e0a05abbf732f7fd95aa30f7cf4947fd)
