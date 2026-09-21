package appsignal_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func TestQuerySendsUserAgent(t *testing.T) {
	var got string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"data": nil}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client := appsignal.NewClient(server.URL, stubToken)
	if err := client.Query(context.Background(), "query Viewer { viewer { id } }", nil, nil); err != nil {
		t.Fatalf("Query() error = %v", err)
	}

	want := "appsignal-client-go/" + appsignal.Version
	if got != want {
		t.Errorf("User-Agent = %q, want %q", got, want)
	}
}

// versionRE guards the literal release-please rewrites: a stray "v" prefix or a
// bare "1.2" would end up in the User-Agent of every request.
var versionRE = regexp.MustCompile(`^\d+\.\d+\.\d+`)

func TestVersionIsSemver(t *testing.T) {
	if !versionRE.MatchString(appsignal.Version) {
		t.Errorf("Version = %q, want a semver like 1.2.3 with no leading v", appsignal.Version)
	}
}
