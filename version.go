package appsignal

// Version is the semantic version of this client library. The trailing comment
// lets release-please rewrite the literal when it cuts a release; keep it.
const Version = "0.1.1" // x-release-please-version

// userAgent identifies this library to the AppSignal API.
const userAgent = "appsignal-client-go/" + Version
