// Package transaction offers transaction manager implementations.
package transaction

import "context"

// NoopManager executes callbacks without transactional guarantees.
type NoopManager struct{}

// NewNoopManager constructs a new manager.
func NewNoopManager() NoopManager {
	return NoopManager{}
}

// WithinTransaction executes the callback immediately.
func (NoopManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
