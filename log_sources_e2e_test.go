//go:build e2e

package appsignal_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func TestGetAppLogSources(t *testing.T) {
	client := testClient(t)

	sources, err := client.GetAppLogSources(context.Background(), testAppID)
	if err != nil {
		t.Fatalf("GetAppLogSources: %v", err)
	}
	if sources == nil {
		t.Fatal("GetAppLogSources returned no log sources and no error")
	}

	for _, source := range sources {
		if source.ID == "" {
			t.Errorf("source %+v has an empty ID", source)
		}
		if source.Name == "" {
			t.Errorf("source %q has an empty Name", source.ID)
		}
	}
}

func TestGetAppLogSource(t *testing.T) {
	client := testClient(t)

	source, err := client.GetAppLogSource(context.Background(), testAppID, testAppLogSourceID)
	if err != nil {
		t.Fatalf("GetAppLogSource: %v", err)
	}
	if source == nil {
		t.Fatal("GetAppLogSource returned no log source and no error")
	}

	if source.ID != testAppLogSourceID {
		t.Errorf("ID = %q, want %q", source.ID, testAppLogSourceID)
	}
}

func TestCreateUpdateAndDeleteAppLogSource(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	name := fmt.Sprintf("client-go-e2e-%d", time.Now().UnixNano())

	created, err := client.CreateAppLogSource(ctx, appsignal.CreateAppLogSourceInput{
		AppID: testAppID,
		Name:  name,
		Type:  "http",
		Fmt:   appsignal.LogSourceFormatPlaintext,
	})
	if err != nil {
		t.Fatalf("CreateAppLogSource: %v", err)
	}
	if created == nil {
		t.Fatal("CreateAppLogSource returned no log source and no error")
	}
	if created.ID == "" {
		t.Fatal("CreateAppLogSource returned a log source without an ID, cannot clean up")
	}

	deleted := false
	t.Cleanup(func() {
		if deleted {
			return
		}
		input := appsignal.DeleteAppLogSourceInput{AppID: testAppID, LogSourceID: created.ID}
		if _, err := client.DeleteAppLogSource(context.Background(), input); err != nil {
			t.Errorf("cleanup DeleteAppLogSource(%q): %v", created.ID, err)
		}
	})

	if created.Name != name {
		t.Errorf("Name = %q, want %q", created.Name, name)
	}
	if created.Fmt != appsignal.LogSourceFormatPlaintext {
		t.Errorf("Fmt = %q, want %q", created.Fmt, appsignal.LogSourceFormatPlaintext)
	}
	if created.Key == "" {
		t.Error("Key is empty")
	}

	got, err := client.GetAppLogSource(ctx, testAppID, created.ID)
	if err != nil {
		t.Errorf("GetAppLogSource after create: %v", err)
	} else if got == nil {
		t.Error("GetAppLogSource after create returned no log source and no error")
	} else if got.ID != created.ID {
		t.Errorf("GetAppLogSource after create: ID = %q, want %q", got.ID, created.ID)
	}

	updatedName := name + "-updated"
	updated, err := client.UpdateAppLogSource(ctx, appsignal.UpdateAppLogSourceInput{
		AppID:       testAppID,
		LogSourceID: created.ID,
		Name:        updatedName,
		Fmt:         appsignal.LogSourceFormatJSON,
	})
	if err != nil {
		t.Fatalf("UpdateAppLogSource: %v", err)
	}
	if updated == nil {
		t.Fatal("UpdateAppLogSource returned no log source and no error")
	}

	if updated.ID != created.ID {
		t.Errorf("UpdateAppLogSource: ID = %q, want %q", updated.ID, created.ID)
	}
	if updated.Name != updatedName {
		t.Errorf("UpdateAppLogSource: Name = %q, want %q", updated.Name, updatedName)
	}
	if updated.Fmt != appsignal.LogSourceFormatJSON {
		t.Errorf("UpdateAppLogSource: Fmt = %q, want %q", updated.Fmt, appsignal.LogSourceFormatJSON)
	}

	got, err = client.GetAppLogSource(ctx, testAppID, created.ID)
	if err != nil {
		t.Errorf("GetAppLogSource after update: %v", err)
	} else if got == nil {
		t.Error("GetAppLogSource after update returned no log source and no error")
	} else if got.Name != updatedName {
		t.Errorf("GetAppLogSource after update: Name = %q, want %q", got.Name, updatedName)
	}

	del, err := client.DeleteAppLogSource(ctx, appsignal.DeleteAppLogSourceInput{
		AppID:       testAppID,
		LogSourceID: created.ID,
	})
	if err != nil {
		t.Fatalf("DeleteAppLogSource: %v", err)
	}
	deleted = true

	if del == nil {
		t.Error("DeleteAppLogSource returned no log source and no error")
	} else if del.ID != created.ID {
		t.Errorf("DeleteAppLogSource: ID = %q, want %q", del.ID, created.ID)
	}

	source, err := client.GetAppLogSource(ctx, testAppID, created.ID)
	t.Logf("GetAppLogSource after delete: source = %+v, err = %v", source, err)
}
