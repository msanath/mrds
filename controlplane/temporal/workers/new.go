package workers

import (
	"context"

	"github.com/msanath/mrds/controlplane/temporal/activities/mrds"
	"github.com/msanath/mrds/controlplane/temporal/activities/runtime"
	"github.com/msanath/mrds/controlplane/temporal/activities/scheduler"
	"github.com/msanath/mrds/controlplane/temporal/workflows"
	"github.com/msanath/mrds/gen/api/mrdspb"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
)

const (
	DeploymentTaskQueue = "continuos-deployment"
)

func NewWorker(
	ctx context.Context,
	mrdsConn *grpc.ClientConn,
	client client.Client,
	runtimeActivities runtime.RuntimeActivities,
) error {
	w := worker.New(client, DeploymentTaskQueue, worker.Options{})

	// Initialize and Register all the activities
	deploymentPlanActivities := mrds.NewDeploymentPlanActivities(mrdspb.NewDeploymentPlansClient(mrdsConn), w)
	metaInstanceActivities := mrds.NewMetaInstanceActivities(mrdspb.NewMetaInstancesClient(mrdsConn), w)
	schedulerActivities := scheduler.NewSchedulerActivities(
		mrdspb.NewMetaInstancesClient(mrdsConn),
		mrdspb.NewNodesClient(mrdsConn),
		mrdspb.NewDeploymentPlansClient(mrdsConn),
		w,
	)
	// Initialize and Register all the workflows
	_ = workflows.NewDeploymentWorkflow(
		deploymentPlanActivities,
		metaInstanceActivities,
		w,
	)

	// Register the runtime activities
	runtimeActivities.Register(w)

	_ = workflows.NewOperationsWorkflow(
		metaInstanceActivities,
		schedulerActivities,
		runtimeActivities,
		w,
	)

	return w.Start()
}
