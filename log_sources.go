package appsignal

import "context"

// LogSourceFormat is the log line format of a log source. It maps onto the
// SourceFormatEnum GraphQL enum.
type LogSourceFormat string

const (
	LogSourceFormatPlaintext  LogSourceFormat = "PLAINTEXT"
	LogSourceFormatLogfmt     LogSourceFormat = "LOGFMT"
	LogSourceFormatJSON       LogSourceFormat = "JSON"
	LogSourceFormatAutodetect LogSourceFormat = "AUTODETECT"
)

type LogSource struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	Key  string          `json:"key"`
	Type string          `json:"type"`
	Fmt  LogSourceFormat `json:"fmt"`
}

const getLogSourceQuery = `
query GetAppLogSource($appId: String!, $logSourceId: String!) {
	app(id: $appId) {
		logs {
			source(id: $logSourceId) {
				id
				name
				key
				type
				fmt
			}
		}
	}
}
`

func (c *Client) GetAppLogSource(ctx context.Context, appID string, logSourceID string) (*LogSource, error) {
	var out struct {
		App struct {
			Logs struct {
				Source *LogSource `json:"source"`
			} `json:"logs"`
		} `json:"app"`
	}

	if err := c.Query(ctx, getLogSourceQuery, map[string]any{"appId": appID, "logSourceId": logSourceID}, &out); err != nil {
		return nil, err
	}

	return out.App.Logs.Source, nil
}

const createLogSourceMutation = `
mutation CreateLogSource($appId: String!, $fmt: SourceFormatEnum!, $name: String!, $type: String!) {
	createLogSource(appId: $appId, fmt: $fmt, name: $name, type: $type) {
		id
		name
		key
		type
		fmt
	}
}
`

type CreateAppLogSourceInput struct {
	AppID string          `json:"appId"`
	Fmt   LogSourceFormat `json:"fmt"`
	Name  string          `json:"name"`
	Type  string          `json:"type"`
}

func (c *Client) CreateAppLogSource(ctx context.Context, input CreateAppLogSourceInput) (*LogSource, error) {
	var out struct {
		LogSource *LogSource `json:"createLogSource"`
	}

	variables := map[string]any{
		"appId": input.AppID,
		"fmt":   input.Fmt,
		"name":  input.Name,
		"type":  input.Type,
	}

	if err := c.Mutate(ctx, createLogSourceMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.LogSource, nil
}

const updateLogSourceMutation = `
mutation UpdateLogSource($appId: String!, $fmt: SourceFormatEnum!, $logSourceId: String!, $name: String!) {
	updateLogSource(appId: $appId, fmt: $fmt, id: $logSourceId, name: $name) {
		id
		name
		key
		type
		fmt
	}
}
`

type UpdateAppLogSourceInput struct {
	AppID       string          `json:"appId"`
	Fmt         LogSourceFormat `json:"fmt"`
	LogSourceID string          `json:"logSourceId"`
	Name        string          `json:"name"`
}

func (c *Client) UpdateAppLogSource(ctx context.Context, input UpdateAppLogSourceInput) (*LogSource, error) {
	var out struct {
		LogSource *LogSource `json:"updateLogSource"`
	}

	variables := map[string]any{
		"appId":       input.AppID,
		"fmt":         input.Fmt,
		"logSourceId": input.LogSourceID,
		"name":        input.Name,
	}

	if err := c.Mutate(ctx, updateLogSourceMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.LogSource, nil
}

const deleteLogSourceMutation = `
mutation DeleteLogSource($appId: String!, $logSourceId: String!) {
	deleteLogSource(appId: $appId, id: $logSourceId) {
		id
		name
		key
		type
		fmt
	}
}
`

type DeleteAppLogSourceInput struct {
	AppID       string `json:"appId"`
	LogSourceID string `json:"logSourceId"`
}

func (c *Client) DeleteAppLogSource(ctx context.Context, input DeleteAppLogSourceInput) (*LogSource, error) {
	var out struct {
		LogSource *LogSource `json:"deleteLogSource"`
	}

	variables := map[string]any{
		"appId":       input.AppID,
		"logSourceId": input.LogSourceID,
	}

	if err := c.Mutate(ctx, deleteLogSourceMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.LogSource, nil
}
