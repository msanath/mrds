package mrds

import (
	"context"
	"fmt"
	"time"

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
	registry.RegisterActivity(a.GetMetaInstanceByID)
	registry.RegisterActivity(a.GetMetaInstanceByName)
	registry.RegisterActivity(a.WaitForOperationStatusApproved)
	registry.RegisterActivity(a.ListMetaInstance)
	registry.RegisterActivity(a.DeleteMetaInstance)
	registry.RegisterActivity(a.AddRuntimeInstance)
	registry.RegisterActivity(a.UpdateRuntimeStatus)
	registry.RegisterActivity(a.UpdateRuntimeActiveState)
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

// GetMetaInstanceByID fetches MetaInstance details based on ID.
func (c *MetaInstanceActivities) GetMetaInstanceByID(ctx context.Context, req *mrdspb.GetMetaInstanceByIDRequest) (*mrdspb.GetMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Fetching MetaInstance by ID", "request", req)

	resp, err := c.client.GetByID(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
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

type DeleteMetaInstanceRequest struct {
	MetaInstanceID string
}

type DeleteMetaInstanceResponse struct{}

func (c *MetaInstanceActivities) DeleteMetaInstance(ctx context.Context, req *DeleteMetaInstanceRequest) (*DeleteMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Deleting MetaInstance", "request", req)

	// Get the Meta Instance by ID
	metaInstance, err := c.GetMetaInstanceByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{Id: req.MetaInstanceID})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
	}

	_, err = c.client.Delete(ctx, &mrdspb.DeleteMetaInstanceRequest{
		Metadata: metaInstance.Record.Metadata,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to delete MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to delete MetaInstance: %w", err)
	}

	return &DeleteMetaInstanceResponse{}, nil
}

type AddRuntimeInstanceRequest struct {
	MetaInstanceID  string
	RuntimeInstance *mrdspb.RuntimeInstance
}

type AddRuntimeInstanceResponse struct {
	MetaInstance *mrdspb.MetaInstance
}

// AddRuntimeInstance is an activity that interacts with the gRPC service to add a RuntimeInstance to a MetaInstance.
func (c *MetaInstanceActivities) AddRuntimeInstance(ctx context.Context, req *AddRuntimeInstanceRequest) (*AddRuntimeInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Adding RuntimeInstance to MetaInstance", "request", req)

	// Get the Meta Instance by ID
	metaInstance, err := c.GetMetaInstanceByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{Id: req.MetaInstanceID})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
	}

	resp, err := c.client.AddRuntimeInstance(ctx, &mrdspb.AddRuntimeInstanceRequest{
		Metadata:        metaInstance.Record.Metadata,
		RuntimeInstance: req.RuntimeInstance,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to add RuntimeInstance to MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to add RuntimeInstance to MetaInstance: %w", err)
	}

	return &AddRuntimeInstanceResponse{
		MetaInstance: resp.Record,
	}, nil
}

type UpdateRuntimeStatusRequest struct {
	MetaInstanceID    string
	RuntimeInstanceID string
	Status            mrdspb.RuntimeInstanceStatus
}

type UpdateRuntimeStatusResponse struct {
	MetaInstance *mrdspb.MetaInstance
}

// UpdateRuntimeStatus is an activity that interacts with the gRPC service to update the status of a RuntimeInstance.
func (c *MetaInstanceActivities) UpdateRuntimeStatus(ctx context.Context, req *UpdateRuntimeStatusRequest) (*UpdateRuntimeStatusResponse, error) {
	activity.GetLogger(ctx).Info("Updating RuntimeInstance status", "request", req)

	// Get the Meta Instance by ID
	metaInstance, err := c.GetMetaInstanceByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{Id: req.MetaInstanceID})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
	}

	resp, err := c.client.UpdateRuntimeStatus(ctx, &mrdspb.UpdateRuntimeStatusRequest{
		Metadata:          metaInstance.Record.Metadata,
		RuntimeInstanceId: req.RuntimeInstanceID,
		Status:            &req.Status,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update RuntimeInstance status", "error", err)
		return nil, fmt.Errorf("failed to update RuntimeInstance status: %w", err)
	}

	return &UpdateRuntimeStatusResponse{
		MetaInstance: resp.Record,
	}, nil
}

type UpdateRuntimeActiveStateRequest struct {
	MetaInstanceID    string
	RuntimeInstanceID string
	IsActive          bool
}

type UpdateRuntimeActiveStateResponse struct {
	MetaInstance *mrdspb.MetaInstance
}

// UpdateRuntimeActiveState is an activity that interacts with the gRPC service to update the active state of a RuntimeInstance.
func (c *MetaInstanceActivities) UpdateRuntimeActiveState(ctx context.Context, req *UpdateRuntimeActiveStateRequest) (*UpdateRuntimeActiveStateResponse, error) {
	activity.GetLogger(ctx).Info("Updating RuntimeInstance active state", "request", req)

	// Get the Meta Instance by ID
	metaInstance, err := c.GetMetaInstanceByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{Id: req.MetaInstanceID})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
	}

	resp, err := c.client.UpdateRuntimeActiveState(ctx, &mrdspb.UpdateRuntimeActiveStateRequest{
		Metadata:          metaInstance.Record.Metadata,
		RuntimeInstanceId: req.RuntimeInstanceID,
		IsActive:          req.IsActive,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update RuntimeInstance active state", "error", err)
		return nil, fmt.Errorf("failed to update RuntimeInstance active state: %w", err)
	}

	return &UpdateRuntimeActiveStateResponse{
		MetaInstance: resp.Record,
	}, nil
}

type RemoveRuntimeInstanceRequest struct {
	MetaInstanceID    string
	RuntimeInstanceID string
}

type RemoveRuntimeInstanceResponse struct {
	MetaInstance *mrdspb.MetaInstance
}

// RemoveRuntimeInstance is an activity that interacts with the gRPC service to remove a RuntimeInstance from a MetaInstance.
func (c *MetaInstanceActivities) RemoveRuntimeInstance(ctx context.Context, req *RemoveRuntimeInstanceRequest) (*RemoveRuntimeInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Removing RuntimeInstance from MetaInstance", "request", req)

	// Get the Meta Instance by ID
	metaInstance, err := c.GetMetaInstanceByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{Id: req.MetaInstanceID})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
	}

	resp, err := c.client.RemoveRuntimeInstance(ctx, &mrdspb.RemoveRuntimeInstanceRequest{
		Metadata:          metaInstance.Record.Metadata,
		RuntimeInstanceId: req.RuntimeInstanceID,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to remove RuntimeInstance from MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to remove RuntimeInstance from MetaInstance: %w", err)
	}

	return &RemoveRuntimeInstanceResponse{
		MetaInstance: resp.Record,
	}, nil
}

type AddOperationRequest struct {
	MetaInstanceID string
	Operation      *mrdspb.Operation
}

// AddOperationResponse is the response from the AddOperation activity.
type AddOperationResponse struct {
	MetaInstance *mrdspb.MetaInstance
}

// AddOperation is an activity that interacts with the gRPC service to add an Operation to a MetaInstance.
func (c *MetaInstanceActivities) AddOperation(ctx context.Context, req *AddOperationRequest) (*AddOperationResponse, error) {
	activity.GetLogger(ctx).Info("Adding Operation to MetaInstance", "request", req)

	// Get the Meta Instance by ID
	metaInstance, err := c.GetMetaInstanceByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{Id: req.MetaInstanceID})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
	}

	resp, err := c.client.AddOperation(ctx, &mrdspb.AddOperationRequest{
		Metadata:  metaInstance.Record.Metadata,
		Operation: req.Operation,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to add Operation to MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to add Operation to MetaInstance: %w", err)
	}

	return &AddOperationResponse{
		MetaInstance: resp.Record,
	}, nil
}

type UpdateOperationStatusRequest struct {
	MetaInstanceID string
	OperationID    string
	State          mrdspb.OperationState
	Message        string
}

type UdpateOperationStatusResponse struct {
	MetaInstance *mrdspb.MetaInstance
}

func (c *MetaInstanceActivities) UpdateOperationStatus(ctx context.Context, req *UpdateOperationStatusRequest) (*UdpateOperationStatusResponse, error) {
	activity.GetLogger(ctx).Info("Updating operation state", "request", req)

	// Get the Meta Instance by ID
	metaInstance, err := c.GetMetaInstanceByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{Id: req.MetaInstanceID})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
	}

	resp, err := c.client.UpdateOperationStatus(ctx, &mrdspb.UpdateOperationStatusRequest{
		Metadata:    metaInstance.Record.Metadata,
		OperationId: req.OperationID,
		Status: &mrdspb.OperationStatus{
			State:   req.State,
			Message: req.Message,
		},
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to mark Operation status as completed", "error", err)
		return nil, fmt.Errorf("failed to mark Operation status as completed: %w", err)
	}

	return &UdpateOperationStatusResponse{
		MetaInstance: resp.Record,
	}, nil
}

type WaitForOperationStatusApprovedRequest struct {
	MetaInstanceID string
	OperationID    string
}

type WaitForOperationStatusApprovedResponse struct {
	MetaInstance *mrdspb.MetaInstance
}

func (c *MetaInstanceActivities) WaitForOperationStatusApproved(ctx context.Context, req *WaitForOperationStatusApprovedRequest) (*WaitForOperationStatusApprovedResponse, error) {
	activity.GetLogger(ctx).Info("Waiting for Operation status to be approved", "request", req)

	timeInterval := 5 * time.Second

	for {
		resp, err := c.client.GetByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{
			Id: req.MetaInstanceID,
		})
		if err != nil {
			activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
			return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
		}
		for _, operation := range resp.Record.Operations {
			if operation.Id == req.OperationID && operation.Status.State == mrdspb.OperationState_OperationState_APPROVED {
				return &WaitForOperationStatusApprovedResponse{MetaInstance: resp.Record}, nil
			}
		}

		activity.GetLogger(ctx).Info("Operation status is not approved yet, waiting...", "interval", timeInterval)
		time.Sleep(timeInterval)
	}
}

type RemoveOperationRequest struct {
	MetaInstanceID string
	OperationID    string
}

type RemoveOperationResponse struct {
	MetaInstance *mrdspb.MetaInstance
}

// RemoveOperation is an activity that interacts with the gRPC service to remove an Operation from a MetaInstance.
func (c *MetaInstanceActivities) RemoveOperation(ctx context.Context, req *RemoveOperationRequest) (*RemoveOperationResponse, error) {
	activity.GetLogger(ctx).Info("Removing Operation from MetaInstance", "request", req)

	// Get the Meta Instance by ID
	metaInstance, err := c.GetMetaInstanceByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{Id: req.MetaInstanceID})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
	}

	resp, err := c.client.RemoveOperation(ctx, &mrdspb.RemoveOperationRequest{
		Metadata:    metaInstance.Record.Metadata,
		OperationId: req.OperationID,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to remove Operation from MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to remove Operation from MetaInstance: %w", err)
	}

	return &RemoveOperationResponse{
		MetaInstance: resp.Record,
	}, nil
}

type UpdateDeploymentIDRequest struct {
	MetaInstanceID string
	DeploymentID   string
}

// UpdateDeploymentIDResponse is the response from the UpdateDeploymentID activity.
type UpdateDeploymentIDResponse struct {
	MetaInstance *mrdspb.MetaInstance
}

// UpdateDeploymentID is an activity that interacts with the gRPC service to update the DeploymentID of a MetaInstance.
func (c *MetaInstanceActivities) UpdateDeploymentID(ctx context.Context, req *UpdateDeploymentIDRequest) (*UpdateDeploymentIDResponse, error) {
	activity.GetLogger(ctx).Info("Updating DeploymentID of MetaInstance", "request", req)

	// Get the Meta Instance by ID
	metaInstance, err := c.GetMetaInstanceByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{Id: req.MetaInstanceID})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get MetaInstance by ID", "error", err)
		return nil, fmt.Errorf("failed to get MetaInstance by ID: %w", err)
	}

	resp, err := c.client.UpdateDeploymentID(ctx, &mrdspb.UpdateDeploymentIDRequest{
		Metadata:     metaInstance.Record.Metadata,
		DeploymentId: req.DeploymentID,
	})
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update DeploymentID of MetaInstance", "error", err)
		return nil, fmt.Errorf("failed to update DeploymentID of MetaInstance: %w", err)
	}

	return &UpdateDeploymentIDResponse{
		MetaInstance: resp.Record,
	}, nil
}
