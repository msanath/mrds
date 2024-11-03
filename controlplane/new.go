package controlplane

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/controlplane/operators"
	"github.com/msanath/mrds/controlplane/temporal/activities/runtime"
	"github.com/msanath/mrds/controlplane/temporal/workers"
	"github.com/msanath/mrds/gen/api/mrdspb"

	"github.com/msanath/gondolf/pkg/ctxslog"

	temporalclient "go.temporal.io/sdk/client"
	"google.golang.org/grpc"
)

type ControlPlane struct {
	mrdsConn          *grpc.ClientConn
	temporalClient    temporalclient.Client
	runtimeActivities runtime.RuntimeActivities
}

func NewControlPlane(
	mrdsConn *grpc.ClientConn,
	temporalClient temporalclient.Client,
	runtimeActivities runtime.RuntimeActivities,
) *ControlPlane {
	return &ControlPlane{
		mrdsConn:          mrdsConn,
		temporalClient:    temporalClient,
		runtimeActivities: runtimeActivities,
	}
}

func (c *ControlPlane) Start(ctx context.Context) error {
	log := ctxslog.FromContext(ctx)
	log.Info("Starting control plane")

	err := workers.NewWorker(ctx, c.mrdsConn, c.temporalClient, c.runtimeActivities)
	if err != nil {
		return fmt.Errorf("failed to start worker: %w", err)
	}

	deploymentOperator := operators.NewDeploymentOperator(c.temporalClient, mrdspb.NewDeploymentPlansClient(c.mrdsConn))
	go func() {
		err := deploymentOperator.RunBlocking(ctx)
		if err != nil {
			log.Error("failed to run deployment manager", "error", err)
		}
	}()

	operationsOperator := operators.NewOperationsOperator(c.temporalClient, mrdspb.NewMetaInstancesClient(c.mrdsConn))
	go func() {
		err := operationsOperator.RunBlocking(ctx)
		if err != nil {
			log.Error("failed to run operations manager", "error", err)
		}
	}()

	return nil
}
