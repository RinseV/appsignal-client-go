//go:build e2e

package appsignal_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func TestGetAppLogViews(t *testing.T) {
	client := testClient(t)

	views, err := client.GetAppLogViews(context.Background(), testAppID)
	if err != nil {
		t.Fatalf("GetAppLogViews: %v", err)
	}

	for _, view := range views {
		if view.ID == "" {
			t.Errorf("view %+v has an empty ID", view)
		}
		if view.Name == "" {
			t.Errorf("view %q has an empty Name", view.ID)
		}
	}
}

func TestCreateUpdateAndDeleteAppLogView(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	name := fmt.Sprintf("client-go-e2e-%d", time.Now().UnixNano())
	query := "test"
	lineHeight := "compact"

	created, err := client.CreateAppLogView(ctx, appsignal.CreateAppLogViewInput{
		AppID:      testAppID,
		Name:       name,
		Query:      &query,
		Columns:    []string{"timestamp", "message"},
		LineHeight: &lineHeight,
		Severities: []appsignal.LogSeverity{appsignal.SeverityError},
		SourceIDs:  []string{testAppLogSourceID},
	})
	if err != nil {
		t.Fatalf("CreateAppLogView: %v", err)
	}
	if created == nil {
		t.Fatal("CreateAppLogView returned no log view and no error")
	}
	if created.ID == "" {
		t.Fatal("CreateAppLogView returned a log view without an ID, cannot clean up")
	}

	deleted := false
	t.Cleanup(func() {
		if deleted {
			return
		}
		input := appsignal.DeleteAppLogViewInput{AppID: testAppID, LogViewID: created.ID}
		if _, err := client.DeleteAppLogView(context.Background(), input); err != nil {
			t.Errorf("cleanup DeleteAppLogView(%q): %v", created.ID, err)
		}
	})

	if created.Name != name {
		t.Errorf("Name = %q, want %q", created.Name, name)
	}
	if created.Query != query {
		t.Errorf("Query = %q, want %q", created.Query, query)
	}
	if created.LineHeight != lineHeight {
		t.Errorf("LineHeight = %q, want %q", created.LineHeight, lineHeight)
	}
	if len(created.Columns) != 2 {
		t.Errorf("Columns = %v, want 2 entries", created.Columns)
	}
	if len(created.Severities) != 1 || created.Severities[0] != appsignal.SeverityError {
		t.Errorf("Severities = %v, want %v", created.Severities, []appsignal.LogSeverity{appsignal.SeverityError})
	}
	if len(created.SourceIDs) != 1 || created.SourceIDs[0] != testAppLogSourceID {
		t.Errorf("SourceIDs = %v, want %v", created.SourceIDs, []string{testAppLogSourceID})
	}

	got, err := client.GetAppLogView(ctx, testAppID, created.ID)
	if err != nil {
		t.Errorf("GetAppLogView after create: %v", err)
	} else if got == nil {
		t.Error("GetAppLogView after create returned no log view and no error")
	} else if got.Name != name {
		t.Errorf("GetAppLogView after create: Name = %q, want %q", got.Name, name)
	}

	updatedName := name + "-updated"
	updatedQuery := "test-updated"
	updatedLineHeight := "relaxed"

	updated, err := client.UpdateAppLogView(ctx, appsignal.UpdateAppLogViewInput{
		AppID:      testAppID,
		LogViewID:  created.ID,
		Name:       &updatedName,
		Query:      &updatedQuery,
		Columns:    []string{"timestamp", "severity", "message"},
		LineHeight: &updatedLineHeight,
		Severities: []appsignal.LogSeverity{appsignal.SeverityError, appsignal.SeverityWarn},
		SourceIDs:  []string{testAppLogSourceID},
	})
	if err != nil {
		t.Fatalf("UpdateAppLogView: %v", err)
	}
	if updated == nil {
		t.Fatal("UpdateAppLogView returned no log view and no error")
	}

	if updated.ID != created.ID {
		t.Errorf("UpdateAppLogView: ID = %q, want %q", updated.ID, created.ID)
	}
	if updated.Name != updatedName {
		t.Errorf("UpdateAppLogView: Name = %q, want %q", updated.Name, updatedName)
	}
	if updated.Query != updatedQuery {
		t.Errorf("UpdateAppLogView: Query = %q, want %q", updated.Query, updatedQuery)
	}
	if updated.LineHeight != updatedLineHeight {
		t.Errorf("UpdateAppLogView: LineHeight = %q, want %q", updated.LineHeight, updatedLineHeight)
	}
	if len(updated.Columns) != 3 {
		t.Errorf("UpdateAppLogView: Columns = %v, want 3 entries", updated.Columns)
	}
	if len(updated.Severities) != 2 {
		t.Errorf("UpdateAppLogView: Severities = %v, want 2 entries", updated.Severities)
	}

	got, err = client.GetAppLogView(ctx, testAppID, created.ID)
	if err != nil {
		t.Errorf("GetAppLogView after update: %v", err)
	} else if got == nil {
		t.Error("GetAppLogView after update returned no log view and no error")
	} else if got.Name != updatedName {
		t.Errorf("GetAppLogView after update: Name = %q, want %q", got.Name, updatedName)
	}

	del, err := client.DeleteAppLogView(ctx, appsignal.DeleteAppLogViewInput{
		AppID:     testAppID,
		LogViewID: created.ID,
	})
	if err != nil {
		t.Fatalf("DeleteAppLogView: %v", err)
	}
	deleted = true

	if del == nil {
		t.Error("DeleteAppLogView returned no log view and no error")
	} else if del.ID != created.ID {
		t.Errorf("DeleteAppLogView: ID = %q, want %q", del.ID, created.ID)
	}

	view, err := client.GetAppLogView(ctx, testAppID, created.ID)
	t.Logf("GetAppLogView after delete: view = %+v, err = %v", view, err)
}
