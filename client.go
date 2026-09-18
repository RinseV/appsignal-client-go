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

func NewClient(token string) *Client {
	return &Client{
		HostURL:    HostURL,
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
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return &HTTPError{
			StatusCode: res.StatusCode,
			Status:     res.Status,
			Header:     res.Header,
			Body:       body,
		}
	}

	var parsed response
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return fmt.Errorf("decode response: %w", err)
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
