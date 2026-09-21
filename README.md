# AppSignal API Client (Go)

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

### Organizations and apps

```go
org, err := client.GetOrganization(ctx, "my-org")

apps, err := client.GetOrganizationApps(ctx, "my-org")
for _, app := range apps {
	fmt.Println(app.ID, app.Name, app.Environment)
}

app, err := client.GetApp(ctx, "app-id")

app, err = client.CreateApp(ctx, appsignal.CreateAppInput{
	OrganizationSlug: "my-org",
	Name:             "my-app",
	Environment:      "production",
})

app, err = client.DeleteApp(ctx, app.ID)
```

### Log sources

```go
sources, err := client.GetAppLogSources(ctx, appID)

source, err := client.GetAppLogSource(ctx, appID, "source-id")

source, err = client.CreateAppLogSource(ctx, appsignal.CreateAppLogSourceInput{
	AppID: appID,
	Name:  "my-source",
	Type:  "custom",
	Fmt:   appsignal.LogSourceFormatJSON,
})

source, err = client.UpdateAppLogSource(ctx, appsignal.UpdateAppLogSourceInput{
	AppID:       appID,
	LogSourceID: source.ID,
	Name:        "renamed-source",
	Fmt:         appsignal.LogSourceFormatLogfmt,
})

source, err = client.DeleteAppLogSource(ctx, appsignal.DeleteAppLogSourceInput{
	AppID:       appID,
	LogSourceID: source.ID,
})
```

### Log views

Optional fields are pointers, so leaving one out means "don't change it".

```go
query := "error"

views, err := client.GetAppLogViews(ctx, appID)

view, err := client.GetAppLogView(ctx, appID, "view-id")

view, err = client.CreateAppLogView(ctx, appsignal.CreateAppLogViewInput{
	AppID:      appID,
	Name:       "Errors",
	Query:      &query,
	Severities: []appsignal.LogSeverity{appsignal.SeverityError, appsignal.SeverityFatal},
})

name := "Errors and warnings"
view, err = client.UpdateAppLogView(ctx, appsignal.UpdateAppLogViewInput{
	AppID:     appID,
	LogViewID: view.ID,
	Name:      &name,
})

view, err = client.DeleteAppLogView(ctx, appsignal.DeleteAppLogViewInput{
	AppID:     appID,
	LogViewID: view.ID,
})
```

### Log triggers

```go
notifiers, err := client.GetAppNotifiers(ctx, appID)

triggers, err := client.GetAppLogTriggers(ctx, appID)

options := appsignal.NotificationOptionAlways
trigger, err := client.CreateAppLogTrigger(ctx, appsignal.CreateAppLogTriggerInput{
	AppID:               appID,
	Name:                "Too many errors",
	Query:               "error",
	Severities:          []appsignal.LogSeverity{appsignal.SeverityError},
	NotificationOptions: &options,
	NotifierIDs:         []string{notifiers[0].ID},
})

triggerName := "Way too many errors"
trigger, err = client.UpdateAppLogTrigger(ctx, appsignal.UpdateAppLogTriggerInput{
	AppID:        appID,
	LogTriggerID: trigger.ID,
	Name:         &triggerName,
})

trigger, err = client.DeleteAppLogTrigger(ctx, appsignal.DeleteAppLogTriggerInput{
	AppID:        appID,
	LogTriggerID: trigger.ID,
})
```

### Errors

The API returns GraphQL errors as `appsignal.Errors` and non-200 responses as `*appsignal.HTTPError`.

```go
_, err := client.GetApp(ctx, "does-not-exist")

var gqlErrs appsignal.Errors
var httpErr *appsignal.HTTPError

switch {
case errors.As(err, &gqlErrs):
	fmt.Println("api error:", gqlErrs)
case errors.As(err, &httpErr):
	fmt.Println("http status:", httpErr.StatusCode)
}
```

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
