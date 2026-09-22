package appsignal_test

import (
	"context"
	"fmt"
	"log"
	"os"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func ExampleClient_GetApp() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	app, err := client.GetApp(context.Background(), "6aad24d1ba6bc351255e7cb5")
	if err != nil {
		log.Fatal(err)
	}
	if app == nil {
		log.Fatal("no app with that ID")
	}

	fmt.Println(app.Name, app.Environment)
}

func ExampleClient_GetOrganizationApps() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	apps, err := client.GetOrganizationApps(context.Background(), "my-organization")
	if err != nil {
		log.Fatal(err)
	}

	for _, app := range apps {
		fmt.Println(app.ID, app.Name, app.Environment)
	}
}

func ExampleClient_CreateApp() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	app, err := client.CreateApp(context.Background(), appsignal.CreateAppInput{
		OrganizationSlug: "my-organization",
		Name:             "my-app",
		Environment:      "production",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(app.ID)
}

func ExampleClient_DeleteApp() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	// The returned app is the one that was just deleted.
	app, err := client.DeleteApp(context.Background(), "6aad24d1ba6bc351255e7cb5")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("deleted", app.Name)
}
