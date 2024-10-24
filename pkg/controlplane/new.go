package controlplane

import (
	"context"
	"fmt"
	"time"

	"github.com/msanath/mrds/pkg/controlplane/temporal/workers"

	"github.com/msanath/gondolf/pkg/ctxslog"

	temporalclient "go.temporal.io/sdk/client"
	"google.golang.org/grpc"
)

type ControlPlane struct {
	opts ControlPlaneOptions
}

type ControlPlaneOptions struct {
	MRDSConn *grpc.ClientConn
}

func NewTemporalControlPlane(opts ControlPlaneOptions) *ControlPlane {
	return &ControlPlane{
		opts: opts,
	}
}

func (c *ControlPlane) Start(ctx context.Context) error {
	log := ctxslog.FromContext(ctx)
	log.Info("Starting control plane")

	tc, err := temporalclient.Dial(temporalclient.Options{
		HostPort:  "localhost:7233",
		Namespace: "mrds",
	})
	if err != nil {
		return err
	}

	workers.NewWorker(c.opts.MRDSConn, tc)

	we, err := tc.ExecuteWorkflow(ctx, temporalclient.StartWorkflowOptions{
		ID:        fmt.Sprintf("create-deployment-%s", time.Now()),
		TaskQueue: workers.DeploymentTaskQueue,
	}, "CreateDeploymentWorkflow")
	if err != nil {
		return fmt.Errorf("failed to start workflow: %w", err)
	}
	log.Info("Started workflow", "workflowID", we.GetID())

	return nil
}
