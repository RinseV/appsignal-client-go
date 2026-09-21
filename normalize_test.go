package appsignal_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	appsignal "github.com/RinseV/appsignal-client-go"
)

// A missing value can arrive as a null, as an empty list, or not at all. All
// three have to decode to a zero value or an empty slice, never a nil.

func TestLogViewNormalization(t *testing.T) {
	tests := []struct {
		name string
		view string
		want appsignal.LogView
	}{
		{
			name: "values pass through untouched",
			view: `{"id":"1","name":"view","query":"hostname=web","columns":["timestamp","message"],
				"lineHeight":"0","severities":["ERROR"],"sourceIds":["s1"]}`,
			want: appsignal.LogView{
				ID:         "1",
				Name:       "view",
				Query:      "hostname=web",
				Columns:    []string{"timestamp", "message"},
				LineHeight: "0",
				Severities: []appsignal.LogSeverity{appsignal.SeverityError},
				SourceIDs:  []string{"s1"},
			},
		},
		{
			name: "null scalars read as empty strings",
			view: `{"id":"1","name":"view","query":null,"lineHeight":null,
				"columns":[],"severities":[],"sourceIds":[]}`,
			want: emptyLogView(),
		},
		{
			name: "null collections read as empty slices",
			view: `{"id":"1","name":"view","query":"","lineHeight":"",
				"columns":null,"severities":null,"sourceIds":null}`,
			want: emptyLogView(),
		},
		{
			name: "absent fields read as zero values",
			view: `{"id":"1","name":"view"}`,
			want: emptyLogView(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newStubClient(t, map[string]any{
				"GetAppLogView": json.RawMessage(`{"app":{"logView":` + tt.view + `}}`),
			})

			got, err := client.GetAppLogView(context.Background(), "app-id", "view-id")
			if err != nil {
				t.Fatalf("GetAppLogView: %v", err)
			}
			if got == nil {
				t.Fatal("GetAppLogView returned no log view and no error")
			}

			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("GetAppLogView() =\n\t%#v\nwant\n\t%#v", *got, tt.want)
			}
		})
	}
}

// emptyLogView is what every "nothing was sent" case has to produce.
func emptyLogView() appsignal.LogView {
	return appsignal.LogView{
		ID:         "1",
		Name:       "view",
		Query:      "",
		Columns:    []string{},
		LineHeight: "",
		Severities: []appsignal.LogSeverity{},
		SourceIDs:  []string{},
	}
}

// Views inside a list get the same treatment as one fetched on its own.
func TestLogViewsNormalization(t *testing.T) {
	client := newStubClient(t, map[string]any{
		"GetAppLogViews": json.RawMessage(`{"app":{"logViews":[{"id":"1","name":"view"}]}}`),
	})

	views, err := client.GetAppLogViews(context.Background(), "app-id")
	if err != nil {
		t.Fatalf("GetAppLogViews: %v", err)
	}
	if len(views) != 1 {
		t.Fatalf("GetAppLogViews returned %d views, want 1", len(views))
	}

	if !reflect.DeepEqual(views[0], emptyLogView()) {
		t.Errorf("GetAppLogViews()[0] =\n\t%#v\nwant\n\t%#v", views[0], emptyLogView())
	}
}

// The one nil the client still returns: a view the API does not know about.
func TestGetAppLogViewMissing(t *testing.T) {
	client := newStubClient(t, map[string]any{
		"GetAppLogView": json.RawMessage(`{"app":{"logView":null}}`),
	})

	view, err := client.GetAppLogView(context.Background(), "app-id", "view-id")
	if err != nil {
		t.Fatalf("GetAppLogView: %v", err)
	}
	if view != nil {
		t.Errorf("GetAppLogView() = %#v, want nil", view)
	}
}

func TestLogTriggerNormalization(t *testing.T) {
	tests := []struct {
		name    string
		trigger string
		want    appsignal.LogTrigger
	}{
		{
			name: "values pass through untouched",
			trigger: `{"id":"1","name":"trigger","query":"hostname=web","sourceIds":["s1"],
				"actionType":"TRIGGER","description":"a trigger","notificationOptions":"NEVER",
				"notificationTriggerValue":3,"order":2,"severities":["ERROR"],
				"notifiers":[{"id":"n1","name":"Slack","icon":"slack"}]}`,
			want: appsignal.LogTrigger{
				ID:                       "1",
				Name:                     "trigger",
				Query:                    "hostname=web",
				SourceIDs:                []string{"s1"},
				ActionType:               appsignal.ActionTypeTrigger,
				Description:              "a trigger",
				NotificationOptions:      appsignal.NotificationOptionNever,
				NotificationTriggerValue: 3,
				Order:                    2,
				Severities:               []appsignal.LogSeverity{appsignal.SeverityError},
				Notifiers:                []appsignal.Notifier{{ID: "n1", Name: "Slack", Icon: "slack"}},
			},
		},
		{
			name: "null scalars read as zero values",
			trigger: `{"id":"1","name":"trigger","description":null,"notificationOptions":null,
				"notificationTriggerValue":null,"sourceIds":[],"severities":[],"notifiers":[]}`,
			want: emptyLogTrigger(),
		},
		{
			name: "null collections read as empty slices",
			trigger: `{"id":"1","name":"trigger","sourceIds":null,"severities":null,
				"notifiers":null}`,
			want: emptyLogTrigger(),
		},
		{
			name:    "absent fields read as zero values",
			trigger: `{"id":"1","name":"trigger"}`,
			want:    emptyLogTrigger(),
		},
		{
			name:    "a null notifier is dropped rather than carried",
			trigger: `{"id":"1","name":"trigger","notifiers":[null]}`,
			want:    emptyLogTrigger(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newStubClient(t, map[string]any{
				"GetAppLogTriggers": json.RawMessage(`{"app":{"logs":{"triggers":[` + tt.trigger + `]}}}`),
			})

			triggers, err := client.GetAppLogTriggers(context.Background(), "app-id")
			if err != nil {
				t.Fatalf("GetAppLogTriggers: %v", err)
			}
			if len(triggers) != 1 {
				t.Fatalf("GetAppLogTriggers returned %d triggers, want 1", len(triggers))
			}

			if !reflect.DeepEqual(triggers[0], tt.want) {
				t.Errorf("GetAppLogTriggers()[0] =\n\t%#v\nwant\n\t%#v", triggers[0], tt.want)
			}
		})
	}
}

// emptyLogTrigger is what every "nothing was sent" case has to produce.
func emptyLogTrigger() appsignal.LogTrigger {
	return appsignal.LogTrigger{
		ID:         "1",
		Name:       "trigger",
		SourceIDs:  []string{},
		Severities: []appsignal.LogSeverity{},
		Notifiers:  []appsignal.Notifier{},
	}
}
