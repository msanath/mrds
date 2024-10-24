package workers

import (
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/mrdsactivities"
	"github.com/msanath/mrds/pkg/controlplane/temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
)

const (
	DeploymentTaskQueue = "continuos-deployment"
)

func NewWorker(mrdsConn *grpc.ClientConn, client client.Client) error {
	w := worker.New(client, DeploymentTaskQueue, worker.Options{})

	clustersClient := mrdspb.NewClustersClient(mrdsConn)
	clusterActivities := mrdsactivities.NewClusterActivities(clustersClient, w)
	deploymentActivities := workflows.DeploymentActivities{
		Cluster: clusterActivities,
	}
	_ = workflows.NewDeploymentWorkflow(deploymentActivities, w)

	return w.Start()
}
