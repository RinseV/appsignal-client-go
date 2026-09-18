//go:build e2e

package appsignal_test

import (
	"os"
	"strings"
	"testing"

	appsignal "github.com/RinseV/appsignal-client-go"
)

const testAppID = "69808b0ce5250a3229a9f634"
const testOrganizationSlug = "drieam"

func testClient(t *testing.T) *appsignal.Client {
	t.Helper()

	token := os.Getenv("APPSIGNAL_API_TOKEN")
	if token == "" {
		t.Skip("APPSIGNAL_API_TOKEN not set, skipping end-to-end test")
	}

	return appsignal.NewClient(appsignal.HostURL, token)
}

func TestMain(m *testing.M) {
	loadEnvFile(".env")
	os.Exit(m.Run())
}

// loadEnvFile reads a .env file into the process environment so the tests can
// be run without any shell or IDE setup. Variables already present in the
// environment are left alone, so the shell and CI keep the final say. A
// missing or unreadable file is ignored: the tests skip on a missing token
// anyway.
func loadEnvFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		value = strings.TrimSpace(value)
		if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
			value = value[1 : len(value)-1]
		}

		if _, set := os.LookupEnv(key); set {
			continue
		}
		os.Setenv(key, value)
	}
}
