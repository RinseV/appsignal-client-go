//go:build e2e

package appsignal_test

import (
	"context"
	"testing"
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
	if app.Name != "Zandbak EU" {
		t.Errorf("Name = %q, want %q", app.Name, "Zandbak EU")
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
