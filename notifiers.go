package appsignal

import "context"

type Notifier struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

const getNotifiersQuery = `
query GetAppNotifiers($appId: String!) {
  app(id: $appId) {
    notifiers {
      id
      name
			icon
    }
  }
}
`

func (c *Client) GetAppNotifiers(ctx context.Context, appID string) ([]Notifier, error) {
	var out struct {
		App struct {
			Notifiers []Notifier `json:"notifiers"`
		} `json:"app"`
	}

	if err := c.Query(ctx, getNotifiersQuery, map[string]any{"appId": appID}, &out); err != nil {
		return nil, err
	}

	return out.App.Notifiers, nil
}
