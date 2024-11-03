package workflows

import (
	"fmt"
	"time"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/mrds"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/runtime"
	"github.com/msanath/mrds/pkg/controlplane/temporal/activities/scheduler"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type OperationsWorkflow struct {
	metaInstanceActivities *mrds.MetaInstanceActivities
	schedulerActivities    *scheduler.SchedulerActivities
	runtimeActivities      runtime.RuntimeActivity
}

func NewOperationsWorkflow(
	metaInstanceActivities *mrds.MetaInstanceActivities,
	schedulerActivities *scheduler.SchedulerActivities,
	runtimeActivities runtime.RuntimeActivity,
	registry worker.Registry,
) *OperationsWorkflow {
	d := &OperationsWorkflow{
		metaInstanceActivities: metaInstanceActivities,
		schedulerActivities:    schedulerActivities,
		runtimeActivities:      runtimeActivities,
	}

	registry.RegisterWorkflow(d.RunOperation)
	return d
}

const OperationsWorkflowName = "RunOperation"

type RunOperationWorkflowParams struct {
	MetaInstanceID string
	OperationID    string
	OperationType  string
}

func (d *OperationsWorkflow) RunOperation(ctx workflow.Context, params RunOperationWorkflowParams) (*mrdspb.UpdateMetaInstanceResponse, error) {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 1 * time.Hour,
		StartToCloseTimeout:    1 * time.Hour,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Get the meta instance
	var getMetaInstanceResponse mrdspb.GetMetaInstanceResponse
	log.Info("Getting meta instance", "metaInstanceID", params.MetaInstanceID)
	err := workflow.ExecuteActivity(ctx, d.metaInstanceActivities.GetMetaInstanceByID, &mrdspb.GetMetaInstanceByIDRequest{
		Id: params.MetaInstanceID,
	}).Get(ctx, &getMetaInstanceResponse)
	if err != nil {
		return nil, err
	}
	log.Info("Got meta instance", "metaInstance", getMetaInstanceResponse.Record)

	var operation *mrdspb.Operation
	for _, op := range getMetaInstanceResponse.Record.Operations {
		if op.Id == params.OperationID {
			operation = op
			break
		}
	}
	if operation == nil {
		return nil, fmt.Errorf("operation with ID %s not found", params.OperationID)
	}

	// If the operation is of type CREATE or relocate, then allocate a new runtime instance.
	var updateMetaInstanceResponse mrdspb.UpdateMetaInstanceResponse
	switch params.OperationType {
	case "CREATE":
		log.Info("Creating a new runtime instance")
		err := workflow.ExecuteActivity(ctx, d.schedulerActivities.AllocateRuntimeInstance, scheduler.AllocateRuntimeInstanceParams{
			MetaInstanceID: params.MetaInstanceID,
			IsActive:       true,
		}).Get(ctx, &updateMetaInstanceResponse)
		if err != nil {
			return nil, err
		}
		log.Info("Allocated runtime instance", "metaInstance", updateMetaInstanceResponse.Record)
	case "RELOCATE":
		log.Info("Creaing a new runtime instance to relocate to")
		err := workflow.ExecuteActivity(ctx, d.schedulerActivities.AllocateRuntimeInstance, scheduler.AllocateRuntimeInstanceParams{
			MetaInstanceID: params.MetaInstanceID,
			IsActive:       false,
		}).Get(ctx, &updateMetaInstanceResponse)
		if err != nil {
			return nil, err
		}
		log.Info("Allocated runtime instance", "metaInstance", updateMetaInstanceResponse.Record)
	}

	// Update the operation status to PENDING_APPROVAL.
	log.Info("Updating operation status to PENDING_APPROVAL")
	var updateOperationStatusResponse mrds.UdpateOperationStatusResponse
	err = workflow.ExecuteActivity(ctx, d.metaInstanceActivities.UpdateOperationStatus, mrds.UpdateOperationStatusRequest{
		MetaInstanceID: params.MetaInstanceID,
		OperationID:    params.OperationID,
		State:          mrdspb.OperationState_OperationState_PENDING_APPROVAL,
		Message:        "Operation is pending approval",
	}).Get(ctx, &updateOperationStatusResponse)
	if err != nil {
		return nil, err
	}
	log.Info("Updated operation status to PENDING_APPROVAL", "metaInstance", updateOperationStatusResponse.MetaInstance)

	// Wait for operation to be approved.
	log.Info("Waiting for operation to be approved")
	var waitResponse mrds.WaitForOperationStatusApprovedResponse
	err = workflow.ExecuteActivity(ctx, d.metaInstanceActivities.WaitForOperationStatusApproved, mrds.WaitForOperationStatusApprovedRequest{
		MetaInstanceID: params.MetaInstanceID,
		OperationID:    params.OperationID,
	}).Get(ctx, &waitResponse)
	if err != nil {
		return nil, err
	}
	log.Info("Operation approved", "metaInstance", waitResponse.MetaInstance)

	// After approval, run the post approval steps based on the type of operation.
	// The operations act on the active instance unless it is a relocate operation.
	log.Info("Running post approval steps")
	var activityErr error

	for _, ri := range waitResponse.MetaInstance.RuntimeInstances {
		if ri.IsActive {
			switch params.OperationType {
			// If the operation type is create or update - start the instance.
			case "CREATE":
				fallthrough
			case "UPDATE":
				log.Info("Starting runtime instance", "metaInstance", waitResponse.MetaInstance, "runtimeInstance", ri)
				err := workflow.ExecuteActivity(ctx, d.runtimeActivities.StartInstance, runtime.RuntimeActivityParams{
					MetaInstanceID:    params.MetaInstanceID,
					RuntimeInstanceID: ri.Id,
				}).Get(ctx, &updateMetaInstanceResponse)
				if err != nil {
					activityErr = err
					break
				}
				log.Info("Started runtime instance", "metaInstance", updateMetaInstanceResponse.Record, "runtimeInstance", ri)

			// If the operation type is relocate - stop the instance and remove it.
			case "RELOCATE":
				log.Info("Stopping runtime instance", "metaInstance", waitResponse.MetaInstance, "runtimeInstance", ri)
				err := workflow.ExecuteActivity(ctx, d.runtimeActivities.StopInstance, runtime.RuntimeActivityParams{
					MetaInstanceID:    params.MetaInstanceID,
					RuntimeInstanceID: ri.Id,
				}).Get(ctx, &updateMetaInstanceResponse)
				if err != nil {
					activityErr = err
					break
				}
				log.Info("Stopped runtime instance", "metaInstance", waitResponse.MetaInstance, "runtimeInstance", ri)

				// Remove the runtime instance.
				var removeRuntimeInstanceResponse mrds.RemoveRuntimeInstanceResponse
				err = workflow.ExecuteActivity(ctx, d.metaInstanceActivities.RemoveRuntimeInstance, &mrds.RemoveRuntimeInstanceRequest{
					MetaInstanceID:    params.MetaInstanceID,
					RuntimeInstanceID: ri.Id,
				}).Get(ctx, &removeRuntimeInstanceResponse)
				if err != nil {
					activityErr = err
					break
				}
				log.Info("Removed runtime instance", "metaInstance", removeRuntimeInstanceResponse.MetaInstance, "runtimeInstance", ri)
			}
		} else {
			switch params.OperationType {
			// If the operation type is relocate - start the passive instance as the relocate operation has been approved.
			case "RELOCATE":
				log.Info("Starting runtime instance", "metaInstance", waitResponse.MetaInstance, "runtimeInstance", ri)
				err := workflow.ExecuteActivity(ctx, d.runtimeActivities.StartInstance, runtime.RuntimeActivityParams{
					MetaInstanceID:    params.MetaInstanceID,
					RuntimeInstanceID: ri.Id,
				}).Get(ctx, &updateMetaInstanceResponse)
				if err != nil {
					activityErr = err
					break
				}
				log.Info("Started runtime instance", "metaInstance", updateMetaInstanceResponse.Record, "runtimeInstance", ri)
			}
		}
	}
	log.Info("Post approval steps completed")

	if activityErr != nil {
		log.Info("Activity failed. Updating operation status to FAILED", "error", activityErr)
		// Update the operation status to FAILED
		var updateOperationStatusResponse mrds.UdpateOperationStatusResponse
		err = workflow.ExecuteActivity(ctx, d.metaInstanceActivities.UpdateOperationStatus, mrds.UpdateOperationStatusRequest{
			MetaInstanceID: params.MetaInstanceID,
			OperationID:    params.OperationID,
			State:          mrdspb.OperationState_OperationState_FAILED,
			Message:        activityErr.Error(),
		}).Get(ctx, &updateOperationStatusResponse)
		if err != nil {
			return nil, err
		}
		return nil, activityErr
	}

	// Update the operation status to SUCCESS
	err = workflow.ExecuteActivity(ctx, d.metaInstanceActivities.UpdateOperationStatus, mrds.UpdateOperationStatusRequest{
		MetaInstanceID: params.MetaInstanceID,
		OperationID:    params.OperationID,
		State:          mrdspb.OperationState_OperationState_SUCCEEDED,
		Message:        "",
	}).Get(ctx, &updateOperationStatusResponse)
	if err != nil {
		return nil, err
	}

	return &updateMetaInstanceResponse, nil
}