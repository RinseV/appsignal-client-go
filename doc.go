// Package appsignal is a client for AppSignal's public GraphQL API.
//
// It wraps the API's organizations, apps, log sources, log views and log
// triggers in typed methods on [Client], and exposes [Client.Query] and
// [Client.Mutate] for anything the typed methods do not cover yet.
//
// # Authentication
//
// Every request is authenticated with an AppSignal personal API token, which
// you can create under Personal settings in the AppSignal dashboard. Pass it to
// [NewClient] and it is sent as a "token" query parameter on each request:
//
//	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))
//	viewer, err := client.GetViewer(context.Background())
//
// # Base URL
//
// An empty host means the public endpoint, [HostURL]:
//
//	client := appsignal.NewClient("", token)
//
// Pass a host to send requests somewhere else.
//
// # HTTP client
//
// Requests go through [http.DefaultClient] unless you replace
// [Client.HTTPClient], which is where you configure timeouts, proxies or a
// transport that retries:
//
//	client := appsignal.NewClient("", token)
//	client.HTTPClient = &http.Client{Timeout: 10 * time.Second}
//
// Every method takes a [context.Context] as well, so individual calls can carry
// their own deadline or be cancelled.
//
// # Errors
//
// The API answers with HTTP 200 and an "errors" array when a request is well
// formed but wrong, for example when an app ID does not exist. Those come back
// as [Errors]. Responses with any other status come back as [HTTPError]. Use
// [errors.As] to tell them apart:
//
//	_, err := client.GetApp(ctx, "does-not-exist")
//
//	var gqlErrs appsignal.Errors
//	var httpErr *appsignal.HTTPError
//
//	switch {
//	case errors.As(err, &gqlErrs):
//		fmt.Println("api error:", gqlErrs)
//	case errors.As(err, &httpErr):
//		fmt.Println("http status:", httpErr.StatusCode)
//	}
//
// Lookups of a single object return a nil pointer and a nil error when nothing
// matches, so check both.
package appsignal
