package appsignal

import "context"

type Viewer struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

const getViewerQuery = `
query GetViewer {
	viewer {
		id
		email
		name
	}
}
`

func (c *Client) GetViewer(ctx context.Context) (*Viewer, error) {
	var out struct {
		Viewer *Viewer `json:"viewer"`
	}

	if err := c.Query(ctx, getViewerQuery, nil, &out); err != nil {
		return nil, err
	}

	return out.Viewer, nil
}
