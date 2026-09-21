package appsignal

// Helpers that turn the nulls from a GraphQL response into zero values.

// deref returns the value behind p, or the zero value of T when p is nil.
func deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}

	return *p
}

// orEmpty returns s, or an empty slice when s is nil.
func orEmpty[T any](s []T) []T {
	if s == nil {
		return []T{}
	}

	return s
}

// derefSlice flattens a slice of pointers, skipping nils. Never returns nil.
func derefSlice[T any](s []*T) []T {
	out := make([]T, 0, len(s))
	for _, p := range s {
		if p != nil {
			out = append(out, *p)
		}
	}

	return out
}
