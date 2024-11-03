package runtime

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"go.temporal.io/sdk/worker"
)

// The runtime interface that any runtime should implement.
type RuntimeActivities interface {
	Register(registry worker.Registry) // Register the activities with the worker
	StartInstance(ctx context.Context, req *RuntimeActivityRequest) (*RuntimeActivityResponse, error)
	StopInstance(ctx context.Context, req *RuntimeActivityRequest) (*RuntimeActivityResponse, error)
}

type RuntimeActivityRequest struct {
	MetaInstanceID    string
	RuntimeInstanceID string
}

type RuntimeActivityResponse struct {
	MetaInstance *mrdspb.MetaInstance
}
