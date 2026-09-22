package appsignal_test

import (
	"context"
	"fmt"
	"log"
	"os"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func ExampleClient_GetAppLogSources() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	sources, err := client.GetAppLogSources(context.Background(), "6aad24d1ba6bc351255e7cb5")
	if err != nil {
		log.Fatal(err)
	}

	for _, source := range sources {
		fmt.Println(source.ID, source.Name, source.Fmt)
	}
}

func ExampleClient_GetAppLogSource() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	source, err := client.GetAppLogSource(context.Background(),
		"6aad24d1ba6bc351255e7cb5", "6aad24d1ba6bc351255e7cba")
	if err != nil {
		log.Fatal(err)
	}
	if source == nil {
		log.Fatal("no log source with that ID")
	}

	fmt.Println(source.Name, source.Fmt)
}

func ExampleClient_CreateAppLogSource() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	source, err := client.CreateAppLogSource(context.Background(), appsignal.CreateAppLogSourceInput{
		AppID: "6aad24d1ba6bc351255e7cb5",
		Name:  "my-log-source",
		Type:  "custom",
		Fmt:   appsignal.LogSourceFormatJSON,
	})
	if err != nil {
		log.Fatal(err)
	}

	// The key is what you ship log lines with.
	fmt.Println(source.ID, source.Key)
}

func ExampleClient_UpdateAppLogSource() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	// Name and Fmt are both required, so pass the current value for whichever
	// one you are not changing.
	source, err := client.UpdateAppLogSource(context.Background(), appsignal.UpdateAppLogSourceInput{
		AppID:       "6aad24d1ba6bc351255e7cb5",
		LogSourceID: "6aad24d1ba6bc351255e7cba",
		Name:        "my-renamed-log-source",
		Fmt:         appsignal.LogSourceFormatLogfmt,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(source.Name, source.Fmt)
}

func ExampleClient_DeleteAppLogSource() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	source, err := client.DeleteAppLogSource(context.Background(), appsignal.DeleteAppLogSourceInput{
		AppID:       "6aad24d1ba6bc351255e7cb5",
		LogSourceID: "6aad24d1ba6bc351255e7cba",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("deleted", source.Name)
}
