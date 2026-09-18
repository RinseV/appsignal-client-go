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
