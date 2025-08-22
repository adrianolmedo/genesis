// Package testhelper provides shared utilities for testing.
package testhelper

import (
	"context"
	"testing"
	"time"
)

// Ctx creates a context with a timeout for testing purposes.
func Ctx(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)
	return ctx
}
