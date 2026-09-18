//go:build e2e

package appsignal_test

import (
	"context"
	"testing"
)

func TestGetViewer(t *testing.T) {
	client := testClient(t)

	viewer, err := client.GetViewer(context.Background())
	if err != nil {
		t.Fatalf("GetViewer: %v", err)
	}
	if viewer == nil {
		t.Fatal("GetViewer returned no viewer and no error")
	}

	if viewer.ID == "" {
		t.Error("GetViewer ID is empty")
	}
}
