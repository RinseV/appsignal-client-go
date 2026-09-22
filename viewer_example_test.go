package appsignal_test

import (
	"context"
	"fmt"
	"log"
	"os"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func ExampleClient_GetViewer() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	// A quick way to check that a token is valid.
	viewer, err := client.GetViewer(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if viewer == nil {
		log.Fatal("token does not belong to a user")
	}

	fmt.Println(viewer.Name, viewer.Email)
}
