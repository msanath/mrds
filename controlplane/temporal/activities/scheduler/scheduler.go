package scheduler

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/msanath/mrds/gen/api/mrdspb"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type SchedulerActivities struct {
	metaInstancesClient   mrdspb.MetaInstancesClient
	nodesClient           mrdspb.NodesClient
	deploymentPlansClient mrdspb.DeploymentPlansClient
}

// NewSchedulerActivities creates a new instance of ClusterActivities.
func NewSchedulerActivities(
	metaInstancesClient mrdspb.MetaInstancesClient,
	nodesClient mrdspb.NodesClient,
	deploymentPlansClient mrdspb.DeploymentPlansClient,
	registry worker.Registry,
) *SchedulerActivities {
	a := &SchedulerActivities{
		metaInstancesClient:   metaInstancesClient,
		nodesClient:           nodesClient,
		deploymentPlansClient: deploymentPlansClient,
	}
	registry.RegisterActivity(a.AllocateRuntimeInstance)
	return a
}

type AllocateRuntimeInstanceParams struct {
	MetaInstanceID string
	IsActive       bool // Indicates if the instance is expected to be active or not.
}

type AllocateRuntimeInstanceResponse struct {
	MetaInstance    *mrdspb.MetaInstance
	RuntimeInstance *mrdspb.RuntimeInstance
}

// CreateCluster is an activity that interacts with the gRPC service to create a Cluster.
func (c *SchedulerActivities) AllocateRuntimeInstance(ctx context.Context, req *AllocateRuntimeInstanceParams) (*AllocateRuntimeInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Creating Cluster", "request", req)

	// Get the meta instance
	metaInstanceGetResp, err := c.metaInstancesClient.GetByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{
		Id: req.MetaInstanceID,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance: %w", err)
	}
	metaInstance := metaInstanceGetResp.Record

	for _, instance := range metaInstance.RuntimeInstances {
		if instance.IsActive == req.IsActive {
			activity.GetLogger(ctx).Error("An instance with the same IsActive value already exists")
			return &AllocateRuntimeInstanceResponse{
				MetaInstance:    metaInstance,
				RuntimeInstance: instance,
			}, nil
		}
	}

	// Get the coresponding deployment Plan
	deploymentPlanGetResp, err := c.deploymentPlansClient.GetByID(ctx, &mrdspb.GetDeploymentPlanByIDRequest{
		Id: metaInstance.DeploymentPlanId,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get Deployment Plan", "error", err)
		return nil, fmt.Errorf("failed to get Deployment Plan: %w", err)
	}
	dp := deploymentPlanGetResp.Record

	payloadNames := make([]string, 0)
	requestedCores := uint32(0)
	requestedMemory := uint32(0)

	for _, app := range dp.Applications {
		payloadNames = append(payloadNames, app.PayloadName)
		requestedCores += app.Resources.Cores
		requestedMemory += app.Resources.Memory
	}

	// Find a node that can accomodate the requested resources and does not have
	// existing instances of the same payload.
	nodeListResp, err := c.nodesClient.List(ctx, &mrdspb.ListNodeRequest{
		StateIn:            []mrdspb.NodeState{mrdspb.NodeState_NodeState_ALLOCATED},
		RemainingCoresGte:  requestedCores,
		RemainingMemoryGte: requestedMemory,
		PayloadNameNotIn:   payloadNames,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to list nodes", "error", err)
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}
	if len(nodeListResp.Records) == 0 {
		activity.GetLogger(ctx).Error("No nodes available to allocate", "error", err)
		return nil, fmt.Errorf("no nodes available to allocate")
	}

	// For now, pick the first node that matches the criteria.
	chosenNode := nodeListResp.Records[0]

	// The metaInstance could've been updated, so get the latest version.
	metaInstanceGetResp, err = c.metaInstancesClient.GetByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{
		Id: metaInstance.Metadata.Id,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance: %w", err)
	}

	// Create the runtime instance
	runtimeInstance := &mrdspb.RuntimeInstance{
		Id:       uuid.New().String(),
		NodeId:   chosenNode.Metadata.Id,
		IsActive: req.IsActive,
		Status: &mrdspb.RuntimeInstanceStatus{
			State:   mrdspb.RuntimeInstanceState_RuntimeState_PENDING,
			Message: "",
		},
	}

	updateResp, err := c.metaInstancesClient.AddRuntimeInstance(ctx, &mrdspb.AddRuntimeInstanceRequest{
		Metadata:        metaInstanceGetResp.Record.Metadata,
		RuntimeInstance: runtimeInstance,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to add Runtime Instance", "error", err)
		return nil, fmt.Errorf("failed to add Runtime Instance: %w", err)
	}

	return &AllocateRuntimeInstanceResponse{
		MetaInstance:    updateResp.Record,
		RuntimeInstance: runtimeInstance,
	}, nil
}
