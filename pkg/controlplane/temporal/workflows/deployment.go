package workflows

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/mrds"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type DeploymentWorkflow struct {
	deploymentPlanActivities *mrds.DeploymentPlanActivities
	metaInstanceActivities   *mrds.MetaInstanceActivities
}

// DeploymentWorkflow is a Temporal workflow that deploys a new cluster.
func NewDeploymentWorkflow(
	deploymentPlan *mrds.DeploymentPlanActivities,
	metaInstance *mrds.MetaInstanceActivities,
	registry worker.Registry,
) *DeploymentWorkflow {

	d := &DeploymentWorkflow{
		deploymentPlanActivities: deploymentPlan,
		metaInstanceActivities:   metaInstance,
	}

	registry.RegisterWorkflow(d.RunDeployment)
	return d
}

const RunDeploymentWorkflowName = "RunDeployment"

type RunDeploymentWorkflowParams struct {
	DeploymentPlan *mrdspb.DeploymentPlanRecord
	Deployment     *mrdspb.Deployment
}

func (d *DeploymentWorkflow) RunDeployment(ctx workflow.Context, params RunDeploymentWorkflowParams) error {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 2 * time.Hour,
		StartToCloseTimeout:    2 * time.Hour,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// 1. Set the deployment state to InProgress
	var updateDeploymentPlanResponse mrds.UpdateDeploymentStatusResponse
	err := workflow.ExecuteActivity(ctx, d.deploymentPlanActivities.UpdateDeploymentStatus, mrds.UpdateDeploymentStatusRequest{
		DeploymentPlanID: params.DeploymentPlan.Metadata.Id,
		DeploymentID:     params.Deployment.Id,
		Status: &mrdspb.DeploymentStatus{
			State:   mrdspb.DeploymentState_DeploymentState_IN_PROGRESS,
			Message: "Deployment is running",
		},
	}).Get(ctx, &updateDeploymentPlanResponse)
	if err != nil {
		return err
	}

	// 2. Get a list of instances that are tagged to the deployment
	var listMetaInstancesResponse mrdspb.ListMetaInstanceResponse
	err = workflow.ExecuteActivity(ctx, d.metaInstanceActivities.ListMetaInstance, &mrdspb.ListMetaInstanceRequest{
		DeploymentPlanIdIn: []string{params.DeploymentPlan.Metadata.Id},
	}).Get(ctx, &listMetaInstancesResponse)
	if err != nil {
		return err
	}
	metaInstances := listMetaInstancesResponse.Records

	remainingInstancesToCreate := int(params.Deployment.InstanceCount) - len(metaInstances)

	if remainingInstancesToCreate > 0 {
		var futures []workflow.Future
		for i := 0; i < remainingInstancesToCreate; i++ {
			f := workflow.ExecuteActivity(ctx, d.metaInstanceActivities.CreateMetaInstance, &mrdspb.CreateMetaInstanceRequest{
				Name:             fmt.Sprintf("%s-instance-%s", params.DeploymentPlan.Name, shortUUID()),
				DeploymentPlanId: params.DeploymentPlan.Metadata.Id,
				DeploymentId:     params.Deployment.Id,
			})
			futures = append(futures, f)
		}

		for _, f := range futures {
			var createMetaInstanceResponse mrdspb.CreateMetaInstanceResponse
			err := f.Get(ctx, &createMetaInstanceResponse)
			if err != nil {
				return err
			}
			workflow.GetLogger(ctx).Info("Activity completed", "InstanceName", createMetaInstanceResponse.Record.Name)
		}
	}

	// Get list of all instances.
	err = workflow.ExecuteActivity(ctx, d.metaInstanceActivities.ListMetaInstance, &mrdspb.ListMetaInstanceRequest{
		DeploymentPlanIdIn: []string{params.DeploymentPlan.Metadata.Id},
	}).Get(ctx, &listMetaInstancesResponse)
	if err != nil {
		return err
	}
	metaInstances = listMetaInstancesResponse.Records

	var operationFutures []workflow.Future
	for _, instance := range metaInstances {
		if instance.DeploymentId != params.Deployment.Id {
			var updateDeploymentIDResponse mrds.UpdateDeploymentIDResponse
			err := workflow.ExecuteActivity(ctx, d.metaInstanceActivities.UpdateDeploymentID, &mrds.UpdateDeploymentIDRequest{
				MetaInstanceID: instance.Metadata.Id,
				DeploymentID:   params.Deployment.Id,
			}).Get(ctx, &updateDeploymentIDResponse)
			if err != nil {
				return err
			}
			instance = updateDeploymentIDResponse.MetaInstance
		}

		if len(instance.RuntimeInstances) == 0 {
			operationID := fmt.Sprintf("CREATE-%s", uuid.New().String())

			var addOperationResponse mrds.AddOperationResponse
			err := workflow.ExecuteActivity(ctx, d.metaInstanceActivities.AddOperation, &mrds.AddOperationRequest{
				MetaInstanceID: instance.Metadata.Id,
				Operation: &mrdspb.Operation{
					Id:       operationID,
					Type:     "CREATE",
					IntentId: instance.DeploymentId,
					Status: &mrdspb.OperationStatus{
						State:   mrdspb.OperationState_OperationState_PREPARING,
						Message: "Instance creation requested",
					},
				},
			}).Get(ctx, &addOperationResponse)
			if err != nil {
				return err
			}

			cwo := workflow.ChildWorkflowOptions{
				WorkflowID: fmt.Sprintf("%s-%s", params.DeploymentPlan.Name, operationID),
			}
			ctx = workflow.WithChildOptions(ctx, cwo)
			operationFutures = append(operationFutures, workflow.ExecuteChildWorkflow(
				ctx,
				OperationsWorkflowName,
				RunOperationWorkflowParams{
					MetaInstanceID: instance.Metadata.Id,
					OperationID:    operationID,
					OperationType:  "CREATE",
				},
			))

		} else {
			operationID := fmt.Sprintf("UPDATE-%s", uuid.New().String())
			var addOperationResponse mrds.AddOperationResponse
			err := workflow.ExecuteActivity(ctx, d.metaInstanceActivities.AddOperation, &mrds.AddOperationRequest{
				MetaInstanceID: instance.Metadata.Id,
				Operation: &mrdspb.Operation{
					Id:       operationID,
					Type:     "UPDATE",
					IntentId: instance.DeploymentId,
					Status: &mrdspb.OperationStatus{
						State:   mrdspb.OperationState_OperationState_PREPARING,
						Message: "Instance creation requested",
					},
				},
			}).Get(ctx, &addOperationResponse)
			if err != nil {
				return err
			}

			cwo := workflow.ChildWorkflowOptions{
				WorkflowID: fmt.Sprintf("%s-%s", params.DeploymentPlan.Name, operationID),
			}
			ctx = workflow.WithChildOptions(ctx, cwo)
			operationFutures = append(operationFutures, workflow.ExecuteChildWorkflow(
				ctx,
				OperationsWorkflowName,
				RunOperationWorkflowParams{
					MetaInstanceID: instance.Metadata.Id,
					OperationID:    operationID,
					OperationType:  "UPDATE",
				},
			))
		}
	}

	for _, f := range operationFutures {
		var updateMetaInstanceResponse mrdspb.UpdateMetaInstanceResponse
		err := f.Get(ctx, &updateMetaInstanceResponse)
		if err != nil {
			return err
		}
		workflow.GetLogger(ctx).Info("Operation completed", "MetaInstanceID", updateMetaInstanceResponse.Record.Metadata.Id)
	}

	// Mark the deployment as completed
	var updateDeploymentStatusResponse mrds.UpdateDeploymentStatusResponse
	err = workflow.ExecuteActivity(ctx, d.deploymentPlanActivities.UpdateDeploymentStatus, &mrds.UpdateDeploymentStatusRequest{
		DeploymentPlanID: params.DeploymentPlan.Metadata.Id,
		DeploymentID:     params.Deployment.Id,
		Status: &mrdspb.DeploymentStatus{
			State:   mrdspb.DeploymentState_DeploymentState_COMPLETED,
			Message: "Deployment completed successfully",
		},
	}).Get(ctx, &updateDeploymentStatusResponse)
	if err != nil {
		return err
	}

	return nil
}

func shortUUID() string {
	u := uuid.New()
	encoded := base64.URLEncoding.EncodeToString(u[:])
	return strings.TrimRight(encoded, "=") // Remove padding characters for a shorter string
}