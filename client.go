package appsignal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const HostURL string = "https://appsignal.com/graphql"

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(host, token string) *Client {
	if host == "" {
		host = HostURL
	}

	return &Client{
		HostURL:    host,
		HTTPClient: http.DefaultClient,
		Token:      token,
	}
}

type request struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type response struct {
	Data   json.RawMessage `json:"data"`
	Errors Errors          `json:"errors,omitempty"`
}

type HTTPError struct {
	StatusCode int
	Status     string
	Header     http.Header
	Body       []byte
}

func (e *HTTPError) Error() string {
	if len(e.Body) > 0 {
		return fmt.Sprintf("appsignal: unexpected status %s: %s", e.Status, e.Body)
	}
	return fmt.Sprintf("appsignal: unexpected status %s", e.Status)
}

type Error struct {
	Message string `json:"message"`
}

type Errors []Error

func (e Errors) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Message
	}
	return strings.Join(msgs, "; ")
}

func (c *Client) Query(ctx context.Context, query string, variables map[string]any, out any) error {
	body, err := json.Marshal(request{Query: query, Variables: variables})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	u, err := url.Parse(c.HostURL)
	if err != nil {
		return fmt.Errorf("parse host url: %w", err)
	}
	q := u.Query()
	q.Set("token", c.Token)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", c.redact(err))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", c.redact(err))
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil && res.StatusCode == http.StatusOK {
		return fmt.Errorf("read response: %w", err)
	}

	var parsed response
	decodeErr := json.Unmarshal(respBody, &parsed)

	if res.StatusCode != http.StatusOK && (decodeErr != nil || len(parsed.Errors) == 0) {
		return &HTTPError{
			StatusCode: res.StatusCode,
			Status:     res.Status,
			Header:     res.Header,
			Body:       respBody,
		}
	}

	if decodeErr != nil {
		return fmt.Errorf("decode response: %w", decodeErr)
	}

	if out != nil && len(parsed.Data) > 0 {
		if err := json.Unmarshal(parsed.Data, out); err != nil {
			return fmt.Errorf("unmarshal data: %w", err)
		}
	}

	if len(parsed.Errors) > 0 {
		return parsed.Errors
	}

	return nil
}

func (c *Client) Mutate(ctx context.Context, mutation string, variables map[string]any, out any) error {
	return c.Query(ctx, mutation, variables, out)
}

func (c *Client) redact(err error) error {
	if err == nil || c.Token == "" {
		return err
	}
	return &redactedError{err: err, secret: c.Token}
}

type redactedError struct {
	err    error
	secret string
}

func (e *redactedError) Error() string {
	return strings.ReplaceAll(e.err.Error(), e.secret, "REDACTED")
}

func (e *redactedError) Unwrap() error { return e.err }
