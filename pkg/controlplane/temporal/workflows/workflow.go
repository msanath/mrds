package workflows

import (
	"time"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/mrdsactivities"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type DeploymentWorkflow struct {
	activities DeploymentActivities
}

// DeploymentWorkflow is a Temporal workflow that deploys a new cluster.
func NewDeploymentWorkflow(activities DeploymentActivities, registry worker.Registry) *DeploymentWorkflow {
	d := &DeploymentWorkflow{
		activities: activities,
	}

	registry.RegisterWorkflow(d.CreateDeploymentWorkflow)
	return d
}

type DeploymentActivities struct {
	Cluster *mrdsactivities.ClusterActivities
}

func (d *DeploymentWorkflow) CreateDeploymentWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var createResponse mrdspb.CreateClusterResponse
	err := workflow.ExecuteActivity(ctx, d.activities.Cluster.CreateCluster, &mrdspb.CreateClusterRequest{
		Name: "my-cluster",
	}).Get(ctx, &createResponse)
	if err != nil {
		return err
	}

	var updateResponse mrdspb.UpdateClusterResponse
	err = workflow.ExecuteActivity(ctx, d.activities.Cluster.UpdateClusterStatus, &mrdspb.UpdateClusterStatusRequest{
		Metadata: createResponse.Record.Metadata,
		Status: &mrdspb.ClusterStatus{
			State:   mrdspb.ClusterState_ClusterState_ACTIVE,
			Message: "Cluster is running",
		},
	}).Get(ctx, &updateResponse)
	if err != nil {
		return err
	}

	var deleteResponse mrdspb.DeleteClusterResponse
	err = workflow.ExecuteActivity(ctx, d.activities.Cluster.DeleteCluster, &mrdspb.DeleteClusterRequest{
		Metadata: updateResponse.Record.Metadata,
	}).Get(ctx, &deleteResponse)
	if err != nil {
		return err
	}

	return nil
}
