# Operation Parameters

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