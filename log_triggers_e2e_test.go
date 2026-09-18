//go:build e2e

package appsignal_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	appsignal "github.com/RinseV/appsignal-client-go"
)

// findLogTrigger returns the trigger with the given ID from the app's trigger
// list, or nil when it is not there. There is no single-trigger query on the
// API, so the tests read back through the list instead.
func findLogTrigger(t *testing.T, client *appsignal.Client, appID, triggerID string) *appsignal.LogTrigger {
	t.Helper()

	triggers, err := client.GetAppLogTriggers(context.Background(), appID)
	if err != nil {
		t.Errorf("GetAppLogTriggers: %v", err)
		return nil
	}

	for i := range triggers {
		if triggers[i].ID == triggerID {
			return &triggers[i]
		}
	}

	return nil
}

func TestGetAppLogTriggers(t *testing.T) {
	client := testClient(t)

	triggers, err := client.GetAppLogTriggers(context.Background(), testAppID)
	if err != nil {
		t.Fatalf("GetAppLogTriggers: %v", err)
	}

	for _, trigger := range triggers {
		if trigger.ID == "" {
			t.Errorf("trigger %+v has an empty ID", trigger)
		}
		if trigger.Name == "" {
			t.Errorf("trigger %q has an empty Name", trigger.ID)
		}
	}
}

func TestCreateUpdateAndDeleteAppLogTrigger(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	name := fmt.Sprintf("client-go-e2e-%d", time.Now().UnixNano())
	description := "created by the appsignal-client-go end-to-end tests"
	notificationOptions := appsignal.NotificationOptionAlways

	created, err := client.CreateAppLogTrigger(ctx, appsignal.CreateAppLogTriggerInput{
		AppID:               testAppID,
		Name:                name,
		Query:               "test",
		Description:         &description,
		NotificationOptions: &notificationOptions,
		Severities:          []appsignal.LogSeverity{appsignal.SeverityError},
		SourceIDs:           []string{testAppLogSourceID},
	})
	if err != nil {
		t.Fatalf("CreateAppLogTrigger: %v", err)
	}
	if created == nil {
		t.Fatal("CreateAppLogTrigger returned no log trigger and no error")
	}
	if created.ID == "" {
		t.Fatal("CreateAppLogTrigger returned a log trigger without an ID, cannot clean up")
	}

	deleted := false
	t.Cleanup(func() {
		if deleted {
			return
		}
		input := appsignal.DeleteAppLogTriggerInput{AppID: testAppID, LogTriggerID: created.ID}
		if _, err := client.DeleteAppLogTrigger(context.Background(), input); err != nil {
			t.Errorf("cleanup DeleteAppLogTrigger(%q): %v", created.ID, err)
		}
	})

	if created.Name != name {
		t.Errorf("Name = %q, want %q", created.Name, name)
	}
	if created.Query != "test" {
		t.Errorf("Query = %q, want %q", created.Query, "test")
	}
	if created.Description == nil {
		t.Error("Description is nil")
	} else if *created.Description != description {
		t.Errorf("Description = %q, want %q", *created.Description, description)
	}
	if created.NotificationOptions == nil {
		t.Error("NotificationOptions is nil")
	} else if *created.NotificationOptions != notificationOptions {
		t.Errorf("NotificationOptions = %q, want %q", *created.NotificationOptions, notificationOptions)
	}
	if len(created.Severities) != 1 || created.Severities[0] != appsignal.SeverityError {
		t.Errorf("Severities = %v, want %v", created.Severities, []appsignal.LogSeverity{appsignal.SeverityError})
	}
	if len(created.SourceIDs) != 1 || created.SourceIDs[0] != testAppLogSourceID {
		t.Errorf("SourceIDs = %v, want %v", created.SourceIDs, []string{testAppLogSourceID})
	}

	if got := findLogTrigger(t, client, testAppID, created.ID); got == nil {
		t.Errorf("GetAppLogTriggers after create: trigger %q is not in the list", created.ID)
	} else if got.Name != name {
		t.Errorf("GetAppLogTriggers after create: Name = %q, want %q", got.Name, name)
	}

	updatedName := name + "-updated"
	updatedQuery := "test-updated"
	updatedNotificationOptions := appsignal.NotificationOptionNever

	updated, err := client.UpdateAppLogTrigger(ctx, appsignal.UpdateAppLogTriggerInput{
		AppID:               testAppID,
		LogTriggerID:        created.ID,
		Name:                &updatedName,
		Query:               &updatedQuery,
		NotificationOptions: &updatedNotificationOptions,
		Severities:          []appsignal.LogSeverity{appsignal.SeverityError, appsignal.SeverityWarn},
		SourceIDs:           []string{testAppLogSourceID},
	})
	if err != nil {
		t.Fatalf("UpdateAppLogTrigger: %v", err)
	}
	if updated == nil {
		t.Fatal("UpdateAppLogTrigger returned no log trigger and no error")
	}

	if updated.ID != created.ID {
		t.Errorf("UpdateAppLogTrigger: ID = %q, want %q", updated.ID, created.ID)
	}
	if updated.Name != updatedName {
		t.Errorf("UpdateAppLogTrigger: Name = %q, want %q", updated.Name, updatedName)
	}
	if updated.Query != updatedQuery {
		t.Errorf("UpdateAppLogTrigger: Query = %q, want %q", updated.Query, updatedQuery)
	}
	if updated.NotificationOptions == nil {
		t.Error("UpdateAppLogTrigger: NotificationOptions is nil")
	} else if *updated.NotificationOptions != updatedNotificationOptions {
		t.Errorf("UpdateAppLogTrigger: NotificationOptions = %q, want %q", *updated.NotificationOptions, updatedNotificationOptions)
	}
	if len(updated.Severities) != 2 {
		t.Errorf("UpdateAppLogTrigger: Severities = %v, want 2 entries", updated.Severities)
	}

	if got := findLogTrigger(t, client, testAppID, created.ID); got == nil {
		t.Errorf("GetAppLogTriggers after update: trigger %q is not in the list", created.ID)
	} else if got.Name != updatedName {
		t.Errorf("GetAppLogTriggers after update: Name = %q, want %q", got.Name, updatedName)
	}

	del, err := client.DeleteAppLogTrigger(ctx, appsignal.DeleteAppLogTriggerInput{
		AppID:        testAppID,
		LogTriggerID: created.ID,
	})
	if err != nil {
		t.Fatalf("DeleteAppLogTrigger: %v", err)
	}
	deleted = true

	if del == nil {
		t.Error("DeleteAppLogTrigger returned no log trigger and no error")
	} else if del.ID != created.ID {
		t.Errorf("DeleteAppLogTrigger: ID = %q, want %q", del.ID, created.ID)
	}

	if got := findLogTrigger(t, client, testAppID, created.ID); got != nil {
		t.Errorf("GetAppLogTriggers after delete: trigger %q is still in the list: %+v", created.ID, got)
	}
}
