package appsignal

import "context"

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

const getOrganizationQuery = `
query GetOrganization($slug: String!) {
	organization(slug: $slug) {
		id
		name
		slug
	}
}
`

func (c *Client) GetOrganization(ctx context.Context, slug string) (*Organization, error) {
	var out struct {
		Organization *Organization `json:"organization"`
	}

	if err := c.Query(ctx, getOrganizationQuery, map[string]any{"slug": slug}, &out); err != nil {
		return nil, err
	}

	return out.Organization, nil
}
