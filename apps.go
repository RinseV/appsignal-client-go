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

// GetApp returns the app with the given ID.
func (c *Client) GetApp(ctx context.Context, id string) (*App, error) {
	var out struct {
		App *App `json:"app"`
	}

	if err := c.Query(ctx, getAppQuery, map[string]any{"id": id}, &out); err != nil {
		return nil, err
	}

	return out.App, nil
}
