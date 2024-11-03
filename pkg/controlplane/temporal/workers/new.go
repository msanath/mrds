package workers

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/mrds"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/runtime"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/scheduler"
	"github.com/msanath/mrds/pkg/controlplane/temporal/workflows"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
)

const (
	DeploymentTaskQueue = "continuos-deployment"
)

func NewWorker(ctx context.Context, mrdsConn *grpc.ClientConn, client client.Client) error {
	w := worker.New(client, DeploymentTaskQueue, worker.Options{})

	kubeconfig := clientcmd.RecommendedHomeFile
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	// Create a clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

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
		clientset,
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
