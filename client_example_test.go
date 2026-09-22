package appsignal_test

// The examples in this package have no "// Output:" comment on purpose: Go
// compiles such examples but never runs them, so they can show real API calls
// without needing a token or a network connection during `go test`.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func ExampleNewClient() {
	// An empty host means the public endpoint, appsignal.HostURL. Pass a real
	// host to point the client at a test server instead.
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	viewer, err := client.GetViewer(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(viewer.Email)
}

func ExampleClient_Query() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	// Any query the typed helpers do not cover can be sent by hand.
	const query = `
query GetAppName($id: String!) {
	app(id: $id) {
		name
	}
}
`

	var out struct {
		App struct {
			Name string `json:"name"`
		} `json:"app"`
	}

	variables := map[string]any{"id": "6aad24d1ba6bc351255e7cb5"}
	if err := client.Query(context.Background(), query, variables, &out); err != nil {
		log.Fatal(err)
	}

	fmt.Println(out.App.Name)
}

func ExampleClient_Mutate() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	const mutation = `
mutation CreateApp($organizationSlug: String!, $name: String!, $environment: String!) {
	createApp(organizationSlug: $organizationSlug, name: $name, environment: $environment) {
		id
	}
}
`

	var out struct {
		App struct {
			ID string `json:"id"`
		} `json:"createApp"`
	}

	variables := map[string]any{
		"organizationSlug": "my-organization",
		"name":             "my-app",
		"environment":      "production",
	}
	if err := client.Mutate(context.Background(), mutation, variables, &out); err != nil {
		log.Fatal(err)
	}

	fmt.Println(out.App.ID)
}

func ExampleHTTPError_Error() {
	client := appsignal.NewClient("", "not-a-real-token")

	_, err := client.GetViewer(context.Background())

	// A bad token gets a non-200 response with no GraphQL errors in it.
	var httpErr *appsignal.HTTPError
	if errors.As(err, &httpErr) {
		// The message carries the status and the response body.
		fmt.Println(httpErr.Error())
		fmt.Println(httpErr.StatusCode)
	}
}

func ExampleErrors_Error() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	_, err := client.GetApp(context.Background(), "not-a-valid-id")

	// The API answers with 200 and a list of GraphQL errors.
	var gqlErrs appsignal.Errors
	if errors.As(err, &gqlErrs) {
		// Error() joins every message with "; ".
		fmt.Println(gqlErrs.Error())

		for _, gqlErr := range gqlErrs {
			fmt.Println(gqlErr.Message)
		}
	}
}
