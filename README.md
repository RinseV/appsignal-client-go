# AppSignal API Client (Go)

A Go client package that can be used to interact with AppSignal's public GraphQL API.

## Tests

To run the E2E tests, you need an AppSignal personal API token. Once acquired, copy the `.env.example` file to `.env` and add your API key. You can then run the tests with:

```bash
go test -tags=e2e -v ./...
```