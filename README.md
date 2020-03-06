# GOASIS
A golang port of [Oasis](https://github.com/x1n13y84issmd42/oasis).

Work in progress.

## TODO
* [x] Fix logging (remove fmt)
* [x] Test response headers
* [ ] Test redirects (should be achieved automatically with testing of headers & status code)
* [ ] Use security in requests
* [x] CLI configuration
* [ ] Unit test it
* [ ] Scripts

### Testing
`go run src/main.go --spec=spec/oasis.yaml  --op="Get number"`

`go run src/main.go --spec=spec/oasis.yaml  --op="List visits"`
