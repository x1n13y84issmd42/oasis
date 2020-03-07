# GOASIS
A golang port of [Oasis](https://github.com/x1n13y84issmd42/oasis).

Work in progress.

## Resources
[OpenAPI Spec](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#securitySchemeObject)
[HTTP Auth Registry](http://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml)

## TODO
* [x] Fix logging (remove fmt)
* [x] Test response headers
* [x] Test redirects (should be achieved automatically with testing of headers & status code)
* [ ] apiKey security
* [ ] HTTP basic security
* [ ] HTTP digest security
* [ ] HTTP bearer security
* [ ] OAuth2 security
* [ ] Test ops with request bodies (POST/PATCH/PUT/etc)
* [x] CLI configuration
* [ ] Unit test it
* [ ] Scripts

### Testing
`go run src/main.go --spec=spec/oasis.yaml  --op="Get number"`

`go run src/main.go --spec=spec/oasis.yaml  --op="List visits"`

`go run src/main.go --spec=spec/oasis.yaml  --op="Meta Number Fail"`
