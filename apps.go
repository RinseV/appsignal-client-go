package appsignal

import (
	"context"
	"time"
)

type App struct {
	ID          string    `json:"id"`
	Environment string    `json:"environment"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const getAppQuery = `
query GetApp($id: String!) {
	app(id: $id) {
		id
		environment
		name
		createdAt
		updatedAt
	}
}
`

// GetApp looks up a single app by its ID. It returns a nil app when no app with
// that ID exists.
func (c *Client) GetApp(ctx context.Context, id string) (*App, error) {
	var out struct {
		App *App `json:"app"`
	}

	if err := c.Query(ctx, getAppQuery, map[string]any{"id": id}, &out); err != nil {
		return nil, err
	}

	return out.App, nil
}

const getOrganizationAppsQuery = `
query GetOrganizationApps($slug: String!) {
	organization(slug: $slug) {
		apps {
			id
			environment
			name
			createdAt
			updatedAt
		}
	}
}
`

// GetOrganizationApps lists every app belonging to the organization with the
// given slug.
func (c *Client) GetOrganizationApps(ctx context.Context, slug string) ([]App, error) {
	var out struct {
		Organization struct {
			Apps []App `json:"apps"`
		} `json:"organization"`
	}

	if err := c.Query(ctx, getOrganizationAppsQuery, map[string]any{"slug": slug}, &out); err != nil {
		return nil, err
	}

	return out.Organization.Apps, nil
}

const createAppMutation = `
mutation CreateApp($organizationSlug: String!, $name: String!, $environment: String!) {
	createApp(organizationSlug: $organizationSlug, name: $name, environment: $environment) {
		id
		environment
		name
		createdAt
		updatedAt
	}
}
`

type CreateAppInput struct {
	OrganizationSlug string `json:"organizationSlug"`
	Name             string `json:"name"`
	Environment      string `json:"environment"`
}

// CreateApp adds a new app to an organization and returns the created app.
func (c *Client) CreateApp(ctx context.Context, input CreateAppInput) (*App, error) {
	var out struct {
		App *App `json:"createApp"`
	}

	variables := map[string]any{
		"organizationSlug": input.OrganizationSlug,
		"name":             input.Name,
		"environment":      input.Environment,
	}

	if err := c.Mutate(ctx, createAppMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.App, nil
}

const deleteAppMutation = `
mutation DeleteApp($appId: String!) {
	deleteApp(appId: $appId) {
		id
		environment
		name
		createdAt
		updatedAt
	}
}
`

// DeleteApp removes an app and returns the app as it was just before deletion.
func (c *Client) DeleteApp(ctx context.Context, appID string) (*App, error) {
	var out struct {
		App *App `json:"deleteApp"`
	}

	variables := map[string]any{
		"appId": appID,
	}

	if err := c.Mutate(ctx, deleteAppMutation, variables, &out); err != nil {
		return nil, err
	}

	return out.App, nil
}
