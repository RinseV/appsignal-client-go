package appsignal_test

import (
	"context"
	"fmt"
	"log"
	"os"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func ExampleClient_GetAppNotifiers() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	notifiers, err := client.GetAppNotifiers(context.Background(), "6aad24d1ba6bc351255e7cb5")
	if err != nil {
		log.Fatal(err)
	}

	// These IDs are what you pass as NotifierIDs on a log trigger.
	for _, notifier := range notifiers {
		fmt.Println(notifier.ID, notifier.Name)
	}
}
