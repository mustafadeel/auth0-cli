package logger0

import "context"

type ListDrainsOpts struct {
	Tenant    string
	PageSize  int32
	PageToken string
}

type DrainRepository interface {
	// Create creates a new drain.
	Create(ctx context.Context, drain *Drain) error

	// Get retrieves a drain by id.
	Get(ctx context.Context, id string) (*Drain, error)

	// List returns a list of drains according to the provided opts, as well as the next token for paging.
	List(ctx context.Context, opts ListDrainsOpts) ([]*Drain, string, error)

	// Update updates a drain.
	Update(ctx context.Context, drain *Drain) error

	// Delete destroys a drains.
	Delete(ctx context.Context, id string) error

	// TODO: Add a Consume() method for stateful streams of drain changes.
}
