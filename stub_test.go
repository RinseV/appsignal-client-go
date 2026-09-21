package appsignal_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	appsignal "github.com/RinseV/appsignal-client-go"
)

// stubToken is not a real credential.
const stubToken = "stub-token"

// operationNameRE pulls "GetAppLogView" out of "query GetAppLogView(".
var operationNameRE = regexp.MustCompile(`(?m)^\s*(?:query|mutation)\s+(\w+)`)

// newStubClient returns a client pointed at a stub of the GraphQL API, so the
// decoding tests need no token and no network. Keys in responses are operation
// names, values are the response "data" field; json.RawMessage expresses a
// null or an absent field exactly.
func newStubClient(t *testing.T, responses map[string]any) *appsignal.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("token"); got != stubToken {
			t.Errorf("request token = %q, want %q", got, stubToken)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		match := operationNameRE.FindStringSubmatch(req.Query)
		if match == nil {
			t.Errorf("no operation name in query: %s", req.Query)
			http.Error(w, "missing operation name", http.StatusBadRequest)
			return
		}

		data, ok := responses[match[1]]
		if !ok {
			t.Errorf("unexpected GraphQL operation %q", match[1])
			http.Error(w, "unexpected operation", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"data": data}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return appsignal.NewClient(server.URL, stubToken)
}
