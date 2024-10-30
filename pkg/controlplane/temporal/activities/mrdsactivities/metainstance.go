package mrdsactivities

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/gen/api/mrdspb"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type MetaInstanceActivities struct {
	client mrdspb.MetaInstancesClient
}

// NewMetaInstanceActivities creates a new instance of MetaInstanceActivities.
func NewMetaInstanceActivities(client mrdspb.MetaInstancesClient, registry worker.Registry) *MetaInstanceActivities {
	a := &MetaInstanceActivities{client: client}
	registry.RegisterActivity(a.CreateMetaInstance)
	registry.RegisterActivity(a.GetMetaInstanceByMetadata)
	registry.RegisterActivity(a.GetMetaInstanceByName)
	registry.RegisterActivity(a.UpdateMetaInstanceStatus)
	registry.RegisterActivity(a.ListMetaInstance)
	registry.RegisterActivity(a.DeleteMetaInstance)
	registry.RegisterActivity(a.AddRuntimeInstance)
	registry.RegisterActivity(a.UpdateRuntimeStatus)
	registry.RegisterActivity(a.RemoveRuntimeInstance)
	registry.RegisterActivity(a.AddOperation)
	registry.RegisterActivity(a.UpdateOperationStatus)
	registry.RegisterActivity(a.RemoveOperation)
	registry.RegisterActivity(a.UpdateDeploymentID)
	return a
}

// CreateMetaInstance is an activity that interacts with the gRPC service to create a MetaInstance.
func (c *MetaInstanceActivities) CreateMetaInstance(ctx context.Context, req *mrdspb.CreateMetaInstanceRequest) (*mrdspb.CreateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Creating MetaInstance", "request", req)

	// Check if the context has a deadline to handle timeout.
	if deadline, ok := ctx.Deadline(); ok {
		activity.GetLogger(ctx).Info("Context has a deadline: %v", deadline)
	}

	// Call gRPC method with context for timeout
	resp, err := c.client.Create(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to create MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to create MetaInstance: %w", err)
	}

	return resp, nil
}

// GetMetaInstanceByMetadata fetches MetaInstance details based on metadata.
func (c *MetaInstanceActivities) GetMetaInstanceByMetadata(ctx context.Context, req *mrdspb.GetMetaInstanceByMetadataRequest) (*mrdspb.GetMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Fetching MetaInstance by metadata", "request", req)

	resp, err := c.client.GetByMetadata(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by metadata", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by metadata: %w", err)
	}

	return resp, nil
}

// GetMetaInstanceByName fetches MetaInstance details based on name.
func (c *MetaInstanceActivities) GetMetaInstanceByName(ctx context.Context, req *mrdspb.GetMetaInstanceByNameRequest) (*mrdspb.GetMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Fetching MetaInstance by name", "request", req)

	resp, err := c.client.GetByName(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by name", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by name: %w", err)
	}

	return resp, nil
}

// UpdateMetaInstanceState updates the state of a MetaInstance.
func (c *MetaInstanceActivities) UpdateMetaInstanceStatus(ctx context.Context, req *mrdspb.UpdateMetaInstanceStatusRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Updating MetaInstance state", "request", req)

	resp, err := c.client.UpdateStatus(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update MetaInstance state", "error", err)
		return nil, fmt.Errorf("failed to update MetaInstance state: %w", err)
	}

	return resp, nil
}

// ListMetaInstance lists all MetaInstances based on the request.
func (c *MetaInstanceActivities) ListMetaInstance(ctx context.Context, req *mrdspb.ListMetaInstanceRequest) (*mrdspb.ListMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Listing MetaInstances", "request", req)

	resp, err := c.client.List(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to list MetaInstances", "error", err)
		return nil, fmt.Errorf("failed to list MetaInstances: %w", err)
	}

	return resp, nil
}

func (c *MetaInstanceActivities) DeleteMetaInstance(ctx context.Context, req *mrdspb.DeleteMetaInstanceRequest) (*mrdspb.DeleteMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Deleting MetaInstance", "request", req)

	resp, err := c.client.Delete(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to delete MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to delete MetaInstance: %w", err)
	}

	return resp, nil
}

// AddRuntimeInstance is an activity that interacts with the gRPC service to add a RuntimeInstance to a MetaInstance.
func (c *MetaInstanceActivities) AddRuntimeInstance(ctx context.Context, req *mrdspb.AddRuntimeInstanceRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Adding RuntimeInstance to MetaInstance", "request", req)

	resp, err := c.client.AddRuntimeInstance(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to add RuntimeInstance to MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to add RuntimeInstance to MetaInstance: %w", err)
	}

	return resp, nil
}

// UpdateRuntimeStatus is an activity that interacts with the gRPC service to update the status of a RuntimeInstance.
func (c *MetaInstanceActivities) UpdateRuntimeStatus(ctx context.Context, req *mrdspb.UpdateRuntimeStatusRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Updating RuntimeInstance status", "request", req)

	resp, err := c.client.UpdateRuntimeStatus(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update RuntimeInstance status", "error", err)
		return nil, fmt.Errorf("failed to update RuntimeInstance status: %w", err)
	}

	return resp, nil
}

// RemoveRuntimeInstance is an activity that interacts with the gRPC service to remove a RuntimeInstance from a MetaInstance.
func (c *MetaInstanceActivities) RemoveRuntimeInstance(ctx context.Context, req *mrdspb.RemoveRuntimeInstanceRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Removing RuntimeInstance from MetaInstance", "request", req)

	resp, err := c.client.RemoveRuntimeInstance(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to remove RuntimeInstance from MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to remove RuntimeInstance from MetaInstance: %w", err)
	}

	return resp, nil
}

// AddOperation is an activity that interacts with the gRPC service to add an Operation to a MetaInstance.
func (c *MetaInstanceActivities) AddOperation(ctx context.Context, req *mrdspb.AddOperationRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Adding Operation to MetaInstance", "request", req)

	resp, err := c.client.AddOperation(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to add Operation to MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to add Operation to MetaInstance: %w", err)
	}

	return resp, nil
}

// UpdateOperationStatus is an activity that interacts with the gRPC service to update the status of an Operation.
func (c *MetaInstanceActivities) UpdateOperationStatus(ctx context.Context, req *mrdspb.UpdateOperationStatusRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Updating Operation status", "request", req)

	resp, err := c.client.UpdateOperationStatus(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update Operation status", "error", err)
		return nil, fmt.Errorf("failed to update Operation status: %w", err)
	}

	return resp, nil
}

// RemoveOperation is an activity that interacts with the gRPC service to remove an Operation from a MetaInstance.
func (c *MetaInstanceActivities) RemoveOperation(ctx context.Context, req *mrdspb.RemoveOperationRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Removing Operation from MetaInstance", "request", req)

	resp, err := c.client.RemoveOperation(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to remove Operation from MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to remove Operation from MetaInstance: %w", err)
	}

	return resp, nil
}

// UpdateDeploymentID is an activity that interacts with the gRPC service to update the DeploymentID of a MetaInstance.
func (c *MetaInstanceActivities) UpdateDeploymentID(ctx context.Context, req *mrdspb.UpdateDeploymentIDRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Updating DeploymentID of MetaInstance", "request", req)

	resp, err := c.client.UpdateDeploymentID(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update DeploymentID of MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to update DeploymentID of MetaInstance: %w", err)
	}

	return resp, nil
}
