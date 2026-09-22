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

// GetViewer returns the user that the client's API token belongs to. It is a
// handy way to check that a token works.
func (c *Client) GetViewer(ctx context.Context) (*Viewer, error) {
	var out struct {
		Viewer *Viewer `json:"viewer"`
	}

	if err := c.Query(ctx, getViewerQuery, nil, &out); err != nil {
		return nil, err
	}

	return out.Viewer, nil
}
