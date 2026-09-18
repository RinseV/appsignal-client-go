//go:build e2e

package appsignal_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func TestGetApp(t *testing.T) {
	client := testClient(t)

	app, err := client.GetApp(context.Background(), testAppID)
	if err != nil {
		t.Fatalf("GetApp: %v", err)
	}
	if app == nil {
		t.Fatal("GetApp returned no app and no error")
	}

	if app.ID != testAppID {
		t.Errorf("ID = %q, want %q", app.ID, testAppID)
	}
	if app.Name != "test-app" {
		t.Errorf("Name = %q, want %q", app.Name, "test-app")
	}
	if app.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}
}

func TestGetAppUnknownID(t *testing.T) {
	client := testClient(t)

	app, err := client.GetApp(context.Background(), "does-not-exist")
	t.Logf("app = %+v, err = %v", app, err)
}

func TestGetOrganizationApps(t *testing.T) {
	client := testClient(t)

	apps, err := client.GetOrganizationApps(context.Background(), testOrganizationSlug)
	if err != nil {
		t.Fatalf("GetOrganizationApps: %v", err)
	}
	if len(apps) == 0 {
		t.Fatal("GetOrganizationApps returned no apps and no error")
	}

	found := false
	for _, app := range apps {
		if app.ID == testAppID {
			found = true
		}
	}
	if !found {
		t.Errorf("App with ID %q not found", testAppID)
	}
}

func TestCreateAndDeleteApp(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	name := fmt.Sprintf("client-go-e2e-%d", time.Now().UnixNano())

	created, err := client.CreateApp(ctx, appsignal.CreateAppInput{
		OrganizationSlug: testOrganizationSlug,
		Name:             name,
		Environment:      "test",
	})
	if err != nil {
		t.Fatalf("CreateApp: %v", err)
	}
	if created == nil {
		t.Fatal("CreateApp returned no app and no error")
	}
	if created.ID == "" {
		t.Fatal("CreateApp returned an app without an ID, cannot clean up")
	}

	deleted := false
	t.Cleanup(func() {
		if deleted {
			return
		}
		if _, err := client.DeleteApp(context.Background(), created.ID); err != nil {
			t.Errorf("cleanup DeleteApp(%q): %v", created.ID, err)
		}
	})

	if created.Name != name {
		t.Errorf("Name = %q, want %q", created.Name, name)
	}
	if created.Environment != "test" {
		t.Errorf("Environment = %q, want %q", created.Environment, "test")
	}
	if created.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}

	got, err := client.GetApp(ctx, created.ID)
	if err != nil {
		t.Errorf("GetApp after create: %v", err)
	} else if got == nil {
		t.Error("GetApp after create returned no app and no error")
	} else if got.ID != created.ID {
		t.Errorf("GetApp after create: ID = %q, want %q", got.ID, created.ID)
	}

	del, err := client.DeleteApp(ctx, created.ID)
	if err != nil {
		t.Fatalf("DeleteApp: %v", err)
	}
	deleted = true

	if del == nil {
		t.Error("DeleteApp returned no app and no error")
	} else if del.ID != created.ID {
		t.Errorf("DeleteApp: ID = %q, want %q", del.ID, created.ID)
	}

	app, err := client.GetApp(ctx, created.ID)
	t.Logf("GetApp after delete: app = %+v, err = %v", app, err)
}
