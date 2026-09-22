package appsignal_test

import (
	"context"
	"fmt"
	"log"
	"os"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func ExampleClient_GetAppLogViews() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	views, err := client.GetAppLogViews(context.Background(), "6aad24d1ba6bc351255e7cb5")
	if err != nil {
		log.Fatal(err)
	}

	for _, view := range views {
		fmt.Println(view.ID, view.Name, view.Query)
	}
}

func ExampleClient_GetAppLogView() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	view, err := client.GetAppLogView(context.Background(),
		"6aad24d1ba6bc351255e7cb5", "6aad24d1ba6bc351255e7cbb")
	if err != nil {
		log.Fatal(err)
	}
	if view == nil {
		log.Fatal("no log view with that ID")
	}

	fmt.Println(view.Name, view.Severities)
}

func ExampleClient_CreateAppLogView() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	// The optional fields are pointers and slices; leave them nil to let the
	// API pick its own defaults.
	query := "level:error"

	view, err := client.CreateAppLogView(context.Background(), appsignal.CreateAppLogViewInput{
		AppID:      "6aad24d1ba6bc351255e7cb5",
		Name:       "Errors",
		Query:      &query,
		Severities: []appsignal.LogSeverity{appsignal.SeverityError, appsignal.SeverityFatal},
		SourceIDs:  []string{"6aad24d1ba6bc351255e7cba"},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(view.ID)
}

func ExampleClient_UpdateAppLogView() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	// Only the fields that are set get sent, so this renames the view and
	// leaves its query, columns and severities alone.
	name := "Errors and worse"

	view, err := client.UpdateAppLogView(context.Background(), appsignal.UpdateAppLogViewInput{
		AppID:     "6aad24d1ba6bc351255e7cb5",
		LogViewID: "6aad24d1ba6bc351255e7cbb",
		Name:      &name,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(view.Name)
}

func ExampleClient_DeleteAppLogView() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	view, err := client.DeleteAppLogView(context.Background(), appsignal.DeleteAppLogViewInput{
		AppID:     "6aad24d1ba6bc351255e7cb5",
		LogViewID: "6aad24d1ba6bc351255e7cbb",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("deleted", view.Name)
}
