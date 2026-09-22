package appsignal_test

import (
	"context"
	"fmt"
	"log"
	"os"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func ExampleClient_GetOrganization() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	organization, err := client.GetOrganization(context.Background(), "my-organization")
	if err != nil {
		log.Fatal(err)
	}
	if organization == nil {
		log.Fatal("no organization with that slug")
	}

	fmt.Println(organization.ID, organization.Name)
}
