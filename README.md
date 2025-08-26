# Go Library for draft-poirier-rats-eat-da

This Go library provides types and utilities for working with Device Attestation (DA) tokens as specified in [draft-poirier-rats-eat-da](https://datatracker.ietf.org/doc/draft-poirier-rats-eat-da/).

## Features

* DA token type and helpers
* CBOR serialization/deserialization
* Utilities for parsing and validating DA tokens
* Unit tests with example vectors

## Installation

```sh
go get github.com/veraison/da
```

## Usage

```go
import "github.com/veraison/da"

// Unmarshal a DA token from CBOR
var token da.DAToken
err := cbor.Unmarshal(cborBytes, &token)

// Marshal a DA token to CBOR
cborBytes, err := cbor.Marshal(token)
```

## Running Tests

```sh
go test -v ./...
```

## References

* [draft-poirier-rats-eat-da](https://datatracker.ietf.org/doc/draft-poirier-rats-eat-da/).
* [EAT Profiles](https://www.rfc-editor.org/rfc/rfc9711.html#section-6)


## License

Apache-2.0. See [LICENSE](LICENSE) for details.