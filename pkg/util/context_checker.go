package util

import "context"

// isContextDone returns true if ctx is canceled.
func IsContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
