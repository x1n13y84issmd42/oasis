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
- [How it works](#how-it-works)
  - [Operation request data](#operation-request-data)
  - [Operation security](#operation-security)
  - [Operation response validation](#operation-response-validation)
    - [HTTP response status code](#http-response-status-code)
    - [HTTP response headers](#http-response-headers)
    - [HTTP response body](#http-response-body)
    - [Schema properties](#schema-properties)
- [Schema extensions](#schema-extensions)
- [Security Object Schema](#security-object-schema)

## Usage
Oasis can be used either in manual mode, which allows to test few operations at a time, or in script mode, which is designed to test complex interaction scenarios, involving multiple endpoints and reusing data across them.

### CLI
Oasis has an unconventional command line interface which syntax resembles natural English naguage, or MySQL. So it configures not with "arguments" in conventional sense, as flags or switches; it's more like a query.

Some clauses are straightforward and resemble flags (`from path/to/spec`), while others may be a bit more complex, like `expect` or `log`. Both of these have "child clauses" which may be used in arbitrary order. For example, if one would like to specify a logging level 1 and plain style, they'd expressed it either as `log at level 4 in plain style` or `log in plain style at level 1`. The `log` clause itself may be anywhere on command line, just like any other top-level clause (see [Arguments](#arguments)).

### Manual mode
`oasis from spec/oasis.yaml test secure-apikey-cookie,meta-query-echo-body,meta-query-echo-headers,meta-headers-echo-body`

`oasis from spec/oasis.yaml test meta-number-fail expect status 200 CT "*"`

`run/oasis from spec/oasis.yaml test meta-bool-fail,meta-bool,meta-number-fail,meta-number log at level 4 in festive style`

```
run/oasis from spec/petstore.yaml \
    test findPetsByStatus \
    log at level 5 \
    use \
      query status=pending \
    expect \
      status 200 \
      CT application/json
```

```
run/oasis from spec/petstore.yaml \
    test getPetById \
    log at level 5 \
    use \
      path parameters petId=9 \
    expect \
      status 200 \
      CT application/json
```

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

### Script mode
Coming soon.

## How it works
It reads an OAS specification YML file, collects information about endpoints, request & response properties (path & query parameters, headers & request bodies), data & security schemas, makes requests to the API and validates responses.

### Operation request data
In order to make make valid requests, Oasis uses example data where available for path & query parameters, request headers & request bodies.

Some components of the OAS spec have been extended with additional Oasis-specific example fields to gain more control over requests. See the [Schema extensions](#schema-extensions) part.

### Operation security
At the moment Oasis supports the following OAS security types:
* API Key (OAS: `type: apiKey`)
* HTTP Basic (OAS: `type: http` & `scheme: basic`) (See [extensions](#security-object-schema))
* HTTP Digest (OAS: `type: http` & `scheme: digest`)(See [extensions](#security-object-schema))

OAuth2 & OpenIdConnect are not supported because they require user interaction.

### Operation response validation
Oasis uses the [OAS Responses](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#responses-object) as a definition of an operation response: a status code, headers & content schema where available.

#### HTTP response status code
The status code from the [OAS Responses](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#responses-object) object is used.

#### HTTP response headers
The [OAS Header](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#header-object) object is used.

If a header has a `required` property, then it must be present in the operation response.

If it has a schema, then the header values will be validated against it. At least one value must be valid in order for test to succeed.

#### HTTP response body
Oasis uses the [OAS Schema](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#schema-object) definition to validate structured response data.

#### Schema properties
Properties' types are checked first.

If a property has a `required` set to `true`, it must be present in response data.

If a property has some of the following validation rules or value formats, those are used to further validate the value.

Supported validation rules:
* multipleOf
* maximum
* exclusiveMaximum
* minimum
* exclusiveMinimum
* maxLength
* minLength
* pattern
* maxItems
* minItems
* uniqueItems
* maxProperties
* minProperties
* required
* enum

Supported value formats:

Format|Meaning|Example
-|-|-
date|A date value|
date-time|A date & time value|
hostname|A host name.|
email|An e-mail address.|
ipv4|A IP V4 address.|127.0.0.1
ipv6|A IP V6 address.|fe80::10db:7611:fbff:2b3d

## Schema extensions
In order to gain more control over the data used in requests, besides the standard and limited `example` field, Oasis introduces few extensions to the OAS schema.

### Security Object Schema
Field Name|Applies To|Description
-|-|-
`x-oasis-username`|HTTP Basic & Digest security|See below.
`x-oasis-password`|HTTP Basic & Digest security|A username & password pair to use for authentication instead of an encoded value from `example`. These fields have priority over the `example` field when present.

## Resources
[OpenAPI Spec](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#securitySchemeObject)

[HTTP Auth Registry](http://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml)

[Stoplight OAS Editor](https://stoplight.io/p/studio/sl/_/6c5faofy)

[Reyesoft API Playground](https://jsonapiplayground.reyesoft.com/)

[Nuxeo API Playground](https://nuxeo.github.io/api-playground/#/resources)

[Oasis V2](https://github.com/x1n13y84issmd42/oasis/tree/c88c9a15e0a05abbf732f7fd95aa30f7cf4947fd)