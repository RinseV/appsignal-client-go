# AppSignal API Client (Go)

[![Go Reference](https://pkg.go.dev/badge/github.com/RinseV/appsignal-client-go.svg)](https://pkg.go.dev/github.com/RinseV/appsignal-client-go)

A Go client package that can be used to interact with AppSignal's public GraphQL API.

## Installation

```bash
go get github.com/RinseV/appsignal-client-go
```

## Usage

Create a client with your AppSignal personal API token. Pass an empty host to use the default (`https://appsignal.com/graphql`).

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func main() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))
	ctx := context.Background()

	viewer, err := client.GetViewer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(viewer.Name, viewer.Email)
}
```

From there the client covers organizations, apps, log sources, log views and log triggers, and lets you send your own queries and mutations for anything it does not wrap yet.

## Documentation

The full API reference, with a usage example for every method, lives on [pkg.go.dev](https://pkg.go.dev/github.com/RinseV/appsignal-client-go). The package documentation also covers authentication, pointing the client at a different host, swapping the HTTP client, and telling GraphQL errors apart from HTTP errors.

## Tests

To run the E2E tests, you need an AppSignal personal API token. Once acquired, copy the `.env.example` file to `.env` and add your API key. You can then run the tests with:

```bash
make test-e2e
```

## Makefile

A `Makefile` wraps the common commands. Run `make` (or `make help`) to list them:

| Command | Description |
| --- | --- |
| `make build` | Compile the package |
| `make test` | Run the unit tests |
| `make test-e2e` | Run the E2E tests, loading `.env` if it exists |
| `make cover` | Run the unit tests and report coverage |
| `make fmt` | Format the code |
| `make vet` | Run `go vet`, including the `e2e`-tagged files |
| `make tidy` | Tidy `go.mod` and `go.sum` |
| `make check` | Format, vet and run the unit tests |
