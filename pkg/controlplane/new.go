package controlplane

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/gen/api/mrdspb"
	continuousdeployment "github.com/msanath/mrds/pkg/controlplane/continuous_deployment"
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
		Logger:    log,
	})
	if err != nil {
		return err
	}

	err = workers.NewWorker(ctx, c.opts.MRDSConn, tc)
	if err != nil {
		return fmt.Errorf("failed to start worker: %w", err)
	}

	deploymentManager := continuousdeployment.NewManager(tc, mrdspb.NewDeploymentPlansClient(c.opts.MRDSConn))

	go func() {
		err := deploymentManager.RunBlocking(ctx)
		if err != nil {
			log.Error("failed to run deployment manager", "error", err)
		}
	}()

	return nil
}
