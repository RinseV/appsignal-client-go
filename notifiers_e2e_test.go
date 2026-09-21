//go:build e2e

package appsignal_test

import (
	"context"
	"testing"
)

func TestGetAppNotifiers(t *testing.T) {
	client := testClient(t)

	notifiers, err := client.GetAppNotifiers(context.Background(), testAppID)
	if err != nil {
		t.Fatalf("GetAppNotifiers: %v", err)
	}

	for _, notifier := range notifiers {
		if notifier.ID == "" {
			t.Errorf("notifier %+v has an empty ID", notifier)
		}
		if notifier.Name == "" {
			t.Errorf("notifier %q has an empty Name", notifier.ID)
		}
	}
}
