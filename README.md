# AppSignal API Client (Go)

A Go client package that can be used to interact with AppSignal's public GraphQL API.

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
