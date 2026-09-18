//go:build e2e

package appsignal_test

import (
	"context"
	"testing"
)

func TestGetOrganization(t *testing.T) {
	client := testClient(t)

	org, err := client.GetOrganization(context.Background(), testOrganizationSlug)
	if err != nil {
		t.Fatalf("GetOrganization: %v", err)
	}
	if org == nil {
		t.Fatal("GetOrganization returned no organization and no error")
	}

	if org.Slug != testOrganizationSlug {
		t.Errorf("Slug = %s; want %s", org.Slug, testOrganizationSlug)
	}
	if org.Name != "Terraform-Test" {
		t.Errorf("Name = %s; want %s", org.Name, "Drieam")
	}
}

func TestGetOrganizationUnknownSlug(t *testing.T) {
	client := testClient(t)

	org, err := client.GetOrganization(context.Background(), "does-not-exist")
	t.Logf("org = %+v, err = %v", org, err)
}
