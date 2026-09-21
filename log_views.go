package appsignal

import (
	"context"
	"encoding/json"
)

type LogView struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Query      string        `json:"query"`
	Columns    []string      `json:"columns"`
	LineHeight string        `json:"lineHeight"`
	Severities []LogSeverity `json:"severities"`
	SourceIDs  []string      `json:"sourceIds"`
}

func (v *LogView) UnmarshalJSON(data []byte) error {
	type alias LogView
	var raw struct {
		alias
		Query      *string `json:"query"`
		LineHeight *string `json:"lineHeight"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*v = LogView(raw.alias)
	v.Query = deref(raw.Query)
	v.LineHeight = deref(raw.LineHeight)
	v.Columns = orEmpty(v.Columns)
	v.Severities = orEmpty(v.Severities)
	v.SourceIDs = orEmpty(v.SourceIDs)

	return nil
}

const getLogViewsQuery = `
query GetAppLogViews($appId: String!) {
	app(id: $appId) {
		logViews {
			id
			name
			query
			columns
			lineHeight
			severities
			sourceIds
		}
	}
}
`

func (c *Client) GetAppLogViews(ctx context.Context, appID string) ([]LogView, error) {
	var out struct {
		App struct {
			LogViews []LogView `json:"logViews"`
		} `json:"app"`
	}

	if err := c.Query(ctx, getLogViewsQuery, map[string]any{"appId": appID}, &out); err != nil {
		return nil, err
	}

	return out.App.LogViews, nil
}

const getLogViewQuery = `
query GetAppLogView($appId: String!, $logViewId: String!) {
	app(id: $appId) {
		logView(id: $logViewId) {
			id
			name
			query
			columns
			lineHeight
			severities
			sourceIds
		}
	}
}
`

func (c *Client) GetAppLogView(ctx context.Context, appID string, logViewID string) (*LogView, error) {
	var out struct {
		App struct {
			LogView *LogView `json:"logView"`
		} `json:"app"`
	}

	if err := c.Query(ctx, getLogViewQuery, map[string]any{"appId": appID, "logViewId": logViewID}, &out); err != nil {
		return nil, err
	}

	return out.App.LogView, nil
}

const createLogViewMutation = `
mutation CreateLogView(
	$appId: String!,
	$name: String!,
  $query: String,
  $columns: [String!],
  $lineHeight: String,
  $severities: [String!],
	$sourceIds: [String!]
) {
  createLogView(
		appId: $appId
		name: $name
		query: $query
		columns: $columns
		lineHeight: $lineHeight
		severities: $severities
		sourceIds: $sourceIds
	) {
		id
		name
		query
		columns
		lineHeight
		severities
		sourceIds
	}
}
`

type CreateAppLogViewInput struct {
	AppID      string        `json:"appId"`
	Name       string        `json:"name"`
	Query      *string       `json:"query,omitempty"`
	Columns    []string      `json:"columns,omitempty"`
	LineHeight *string       `json:"lineHeight,omitempty"`
	Severities []LogSeverity `json:"severities,omitempty"`
	SourceIDs  []string      `json:"sourceIds,omitempty"`
}

func (c *Client) CreateAppLogView(ctx context.Context, input CreateAppLogViewInput) (*LogView, error) {
	var out struct {
		LogView *LogView `json:"createLogView"`
	}

	variables := map[string]any{
		"appId": input.AppID,
		"name":  input.Name,
	}

	if input.Query != nil {
		variables["query"] = *input.Query
	}
	if input.Columns != nil {
		variables["columns"] = input.Columns
	}
	if input.LineHeight != nil {
		variables["lineHeight"] = *input.LineHeight
	}
	if input.Severities != nil {
		variables["severities"] = input.Severities
	}
	if input.SourceIDs != nil {
		variables["sourceIds"] = input.SourceIDs
	}

	if err := c.Mutate(ctx, createLogViewMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.LogView, nil
}

const updateLogViewMutation = `
mutation UpdateLogView(
	$appId: String!,
  $logViewId: String!,
	$name: String,
  $query: String,
  $columns: [String!],
  $lineHeight: String,
  $severities: [String!],
	$sourceIds: [String!]
) {
  updateLogView(
		appId: $appId
    id: $logViewId
		name: $name
		query: $query
		columns: $columns
		lineHeight: $lineHeight
		severities: $severities
		sourceIds: $sourceIds
	) {
		id
		name
		query
		columns
		lineHeight
		severities
		sourceIds
	}
}
`

type UpdateAppLogViewInput struct {
	AppID      string        `json:"appId"`
	LogViewID  string        `json:"logViewId"`
	Name       *string       `json:"name,omitempty"`
	Query      *string       `json:"query,omitempty"`
	Columns    []string      `json:"columns,omitempty"`
	LineHeight *string       `json:"lineHeight,omitempty"`
	Severities []LogSeverity `json:"severities,omitempty"`
	SourceIDs  []string      `json:"sourceIds,omitempty"`
}

func (c *Client) UpdateAppLogView(ctx context.Context, input UpdateAppLogViewInput) (*LogView, error) {
	var out struct {
		LogView *LogView `json:"updateLogView"`
	}

	variables := map[string]any{
		"appId":     input.AppID,
		"logViewId": input.LogViewID,
	}

	if input.Name != nil {
		variables["name"] = *input.Name
	}
	if input.Query != nil {
		variables["query"] = *input.Query
	}
	if input.Columns != nil {
		variables["columns"] = input.Columns
	}
	if input.LineHeight != nil {
		variables["lineHeight"] = *input.LineHeight
	}
	if input.Severities != nil {
		variables["severities"] = input.Severities
	}
	if input.SourceIDs != nil {
		variables["sourceIds"] = input.SourceIDs
	}

	if err := c.Mutate(ctx, updateLogViewMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.LogView, nil
}

const deleteLogViewMutation = `
mutation DeleteLogView(
	$appId: String!,
	$logViewId: String!
) {
	deleteLogView(
		appId: $appId
		id: $logViewId
	) {
		id
		name
		query
		columns
		lineHeight
		severities
		sourceIds
	}
}
`

type DeleteAppLogViewInput struct {
	AppID     string `json:"appId"`
	LogViewID string `json:"logViewId"`
}

func (c *Client) DeleteAppLogView(ctx context.Context, input DeleteAppLogViewInput) (*LogView, error) {
	var out struct {
		LogView *LogView `json:"deleteLogView"`
	}

	variables := map[string]any{
		"appId":     input.AppID,
		"logViewId": input.LogViewID,
	}

	if err := c.Mutate(ctx, deleteLogViewMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.LogView, nil
}
