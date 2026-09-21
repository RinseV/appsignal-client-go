package appsignal

import (
	"context"
)

type LogTrigger struct {
	ID                       string                        `json:"id"`
	Name                     string                        `json:"name"`
	Query                    string                        `json:"query"`
	SourceIDs                []string                      `json:"sourceIds"`
	ActionType               LogTriggerActionType          `json:"actionType"`
	Description              *string                       `json:"description"`
	NotificationOptions      *LogTriggerNotificationOption `json:"notificationOptions"`
	NotificationTriggerValue *int32                        `json:"notificationTriggerValue"`
	Order                    int32                         `json:"order"`
	Severities               []LogSeverity                 `json:"severities"`
	Notifiers                []*Notifier                   `json:"notifiers"`
}

const getLogTriggersQuery = `
query GetAppLogTriggers($appId: String!) {
	app(id: $appId) {
		logs {
			triggers {
				id
				name
				query
				sourceIds
				actionType
				description
				notificationOptions
				notificationTriggerValue
				order
				severities
				notifiers {
					id
					name
					icon
				}
			}
		}
	}
}
`

func (c *Client) GetAppLogTriggers(ctx context.Context, appID string) ([]LogTrigger, error) {
	var out struct {
		App struct {
			Logs struct {
				Triggers []LogTrigger `json:"triggers"`
			} `json:"logs"`
		} `json:"app"`
	}

	if err := c.Query(ctx, getLogTriggersQuery, map[string]any{"appId": appID}, &out); err != nil {
		return nil, err
	}

	return out.App.Logs.Triggers, nil
}

const createLogTriggerMutation = `
mutation CreateLogTrigger(
	$appId: String!,
	$name: String!,
	$query: String!,
	$description: String,
	$notificationOptions: IncidentNotificationFrequencyEnum,
	$notificationTriggerValue: Int,
	$severities: [String!],
	$sourceIds: [String!],
	$notifierIds: [String!],
) {
  createLogTrigger(
    appId: $appId
    name: $name
    query: $query
    description: $description
    notificationOptions: $notificationOptions
    notificationTriggerValue: $notificationTriggerValue
    severities: $severities
    sourceIds: $sourceIds,
		notifierIds: $notifierIds,
  ) {
		id
		name
		query
		sourceIds
		actionType
		description
		notificationOptions
		notificationTriggerValue
		order
		severities
		notifiers {
			id
			name
      icon
		}
	}
}
`

type CreateAppLogTriggerInput struct {
	AppID                    string                        `json:"appId"`
	Name                     string                        `json:"name"`
	Query                    string                        `json:"query"`
	Description              *string                       `json:"description,omitempty"`
	NotificationOptions      *LogTriggerNotificationOption `json:"notificationOptions,omitempty"`
	NotificationTriggerValue *int32                        `json:"notificationTriggerValue,omitempty"`
	Severities               []LogSeverity                 `json:"severities,omitempty"`
	SourceIDs                []string                      `json:"sourceIds,omitempty"`
	NotifierIDs              []string                      `json:"notifierIds,omitempty"`
}

func (c *Client) CreateAppLogTrigger(ctx context.Context, input CreateAppLogTriggerInput) (*LogTrigger, error) {
	var out struct {
		LogTrigger *LogTrigger `json:"createLogTrigger"`
	}

	variables := map[string]any{
		"appId": input.AppID,
		"name":  input.Name,
		"query": input.Query,
	}

	if input.Description != nil {
		variables["description"] = *input.Description
	}
	if input.NotificationOptions != nil {
		variables["notificationOptions"] = *input.NotificationOptions
	}
	if input.NotificationTriggerValue != nil {
		variables["notificationTriggerValue"] = *input.NotificationTriggerValue
	}
	if input.Severities != nil {
		variables["severities"] = input.Severities
	}
	if input.SourceIDs != nil {
		variables["sourceIds"] = input.SourceIDs
	}
	if input.NotifierIDs != nil {
		variables["notifierIds"] = input.NotifierIDs
	}

	if err := c.Mutate(ctx, createLogTriggerMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.LogTrigger, nil
}

const updateLogTriggerMutation = `
mutation UpdateLogTrigger(
	$appId: String!,
	$logTriggerId: String!,
	$name: String,
	$query: String,
	$description: String,
	$notificationOptions: IncidentNotificationFrequencyEnum,
	$notificationTriggerValue: Int,
	$severities: [String!],
	$sourceIds: [String!],
	$notifierIds: [String!],
) {
  updateLogTrigger(
    appId: $appId
    id: $logTriggerId
    name: $name
    query: $query
    description: $description
    notificationOptions: $notificationOptions
    notificationTriggerValue: $notificationTriggerValue
    severities: $severities
    sourceIds: $sourceIds,
		notifierIds: $notifierIds,
  ) {
		id
		name
		query
		sourceIds
		actionType
		description
		notificationOptions
		notificationTriggerValue
		order
		severities
		notifiers {
			id
			name
			icon
		}
	}
}
`

type UpdateAppLogTriggerInput struct {
	AppID                    string                        `json:"appId"`
	LogTriggerID             string                        `json:"logTriggerId"`
	Name                     *string                       `json:"name,omitempty"`
	Query                    *string                       `json:"query,omitempty"`
	Description              *string                       `json:"description,omitempty"`
	NotificationOptions      *LogTriggerNotificationOption `json:"notificationOptions,omitempty"`
	NotificationTriggerValue *int32                        `json:"notificationTriggerValue,omitempty"`
	Severities               []LogSeverity                 `json:"severities,omitempty"`
	SourceIDs                []string                      `json:"sourceIds,omitempty"`
	NotifierIDs              []string                      `json:"notifierIds,omitempty"`
}

func (c *Client) UpdateAppLogTrigger(ctx context.Context, input UpdateAppLogTriggerInput) (*LogTrigger, error) {
	var out struct {
		LogTrigger *LogTrigger `json:"updateLogTrigger"`
	}

	variables := map[string]any{
		"appId":        input.AppID,
		"logTriggerId": input.LogTriggerID,
	}

	if input.Name != nil {
		variables["name"] = *input.Name
	}
	if input.Query != nil {
		variables["query"] = *input.Query
	}
	if input.Description != nil {
		variables["description"] = *input.Description
	}
	if input.NotificationOptions != nil {
		variables["notificationOptions"] = *input.NotificationOptions
	}
	if input.NotificationTriggerValue != nil {
		variables["notificationTriggerValue"] = *input.NotificationTriggerValue
	}
	if input.Severities != nil {
		variables["severities"] = input.Severities
	}
	if input.SourceIDs != nil {
		variables["sourceIds"] = input.SourceIDs
	}
	if input.NotifierIDs != nil {
		variables["notifierIds"] = input.NotifierIDs
	}

	if err := c.Mutate(ctx, updateLogTriggerMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.LogTrigger, nil
}

const deleteLogTriggerMutation = `
mutation DeleteLogTrigger(
	$appId: String!,
	$logTriggerId: String!,
) {
	deleteLogTrigger(appId: $appId, id: $logTriggerId) {
		id
		name
		query
		sourceIds
		actionType
		description
		notificationOptions
		notificationTriggerValue
		order
		severities
		notifiers {
			id
			name
			icon
		}
	}
}
`

type DeleteAppLogTriggerInput struct {
	AppID        string `json:"appId"`
	LogTriggerID string `json:"logTriggerId"`
}

func (c *Client) DeleteAppLogTrigger(ctx context.Context, input DeleteAppLogTriggerInput) (*LogTrigger, error) {
	var out struct {
		LogTrigger *LogTrigger `json:"deleteLogTrigger"`
	}

	variables := map[string]any{
		"appId":        input.AppID,
		"logTriggerId": input.LogTriggerID,
	}

	if err := c.Mutate(ctx, deleteLogTriggerMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.LogTrigger, nil
}
