package workers

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/mrds"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/runtime"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/scheduler"
	"github.com/msanath/mrds/pkg/controlplane/temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
)

const (
	DeploymentTaskQueue = "continuos-deployment"
)

func NewWorker(ctx context.Context, mrdsConn *grpc.ClientConn, client client.Client) error {
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
	runtimeActivities := runtime.NewKindRuntime(
		mrdspb.NewMetaInstancesClient(mrdsConn),
		mrdspb.NewDeploymentPlansClient(mrdsConn),
		mrdspb.NewNodesClient(mrdsConn),
		w,
	)

	// Initialize and Register all the workflows
	_ = workflows.NewDeploymentWorkflow(
		deploymentPlanActivities,
		metaInstanceActivities,
		w,
	)

	_ = workflows.NewOperationsWorkflow(
		metaInstanceActivities,
		schedulerActivities,
		runtimeActivities,
		w,
	)

	return w.Start()
}
