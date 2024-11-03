package operators

import "context"

type Operator interface {
	RunBlocking(ctx context.Context) error
}
