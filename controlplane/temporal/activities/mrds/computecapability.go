package mrds

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/gen/api/mrdspb"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type ComputeCapabilityActivities struct {
	client mrdspb.ComputeCapabilitiesClient
}

// NewComputeCapabilityActivities creates a new instance of ComputeCapabilityActivities.
func NewComputeCapabilityActivities(client mrdspb.ComputeCapabilitiesClient, registry worker.Registry) *ComputeCapabilityActivities {
	a := &ComputeCapabilityActivities{client: client}
	registry.RegisterActivity(a.CreateComputeCapability)
	registry.RegisterActivity(a.GetComputeCapabilityByID)
	registry.RegisterActivity(a.GetComputeCapabilityByName)
	registry.RegisterActivity(a.UpdateComputeCapabilityStatus)
	registry.RegisterActivity(a.ListComputeCapability)
	registry.RegisterActivity(a.DeleteComputeCapability)
	return a
}

// CreateComputeCapability is an activity that interacts with the gRPC service to create a ComputeCapability.
func (c *ComputeCapabilityActivities) CreateComputeCapability(ctx context.Context, req *mrdspb.CreateComputeCapabilityRequest) (*mrdspb.CreateComputeCapabilityResponse, error) {
	activity.GetLogger(ctx).Info("Creating ComputeCapability", "request", req)

	// Check if the context has a deadline to handle timeout.
	if deadline, ok := ctx.Deadline(); ok {
		activity.GetLogger(ctx).Info("Context has a deadline: %v", deadline)
	}

	// Call gRPC method with context for timeout
	resp, err := c.client.Create(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to create ComputeCapability", "error", err)
		return nil, fmt.Errorf("failed to create ComputeCapability: %w", err)
	}

	return resp, nil
}

// GetComputeCapabilityByID fetches ComputeCapability details based on ID.
func (c *ComputeCapabilityActivities) GetComputeCapabilityByID(ctx context.Context, req *mrdspb.GetComputeCapabilityByIDRequest) (*mrdspb.GetComputeCapabilityResponse, error) {
	activity.GetLogger(ctx).Info("Fetching ComputeCapability by ID", "request", req)

	resp, err := c.client.GetByID(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get ComputeCapability by ID", "error", err)
		return nil, fmt.Errorf("failed to get ComputeCapability by ID: %w", err)
	}

	return resp, nil
}

// GetComputeCapabilityByName fetches ComputeCapability details based on name.
func (c *ComputeCapabilityActivities) GetComputeCapabilityByName(ctx context.Context, req *mrdspb.GetComputeCapabilityByNameRequest) (*mrdspb.GetComputeCapabilityResponse, error) {
	activity.GetLogger(ctx).Info("Fetching ComputeCapability by name", "request", req)

	resp, err := c.client.GetByName(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get ComputeCapability by name", "error", err)
		return nil, fmt.Errorf("failed to get ComputeCapability by name: %w", err)
	}

	return resp, nil
}

// UpdateComputeCapabilityState updates the state of a ComputeCapability.
func (c *ComputeCapabilityActivities) UpdateComputeCapabilityStatus(ctx context.Context, req *mrdspb.UpdateComputeCapabilityStatusRequest) (*mrdspb.UpdateComputeCapabilityResponse, error) {
	activity.GetLogger(ctx).Info("Updating ComputeCapability state", "request", req)

	resp, err := c.client.UpdateStatus(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update ComputeCapability state", "error", err)
		return nil, fmt.Errorf("failed to update ComputeCapability state: %w", err)
	}

	return resp, nil
}

// ListComputeCapability lists all ComputeCapabilitys based on the request.
func (c *ComputeCapabilityActivities) ListComputeCapability(ctx context.Context, req *mrdspb.ListComputeCapabilityRequest) (*mrdspb.ListComputeCapabilityResponse, error) {
	activity.GetLogger(ctx).Info("Listing ComputeCapabilitys", "request", req)

	resp, err := c.client.List(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to list ComputeCapabilitys", "error", err)
		return nil, fmt.Errorf("failed to list ComputeCapabilitys: %w", err)
	}

	return resp, nil
}

func (c *ComputeCapabilityActivities) DeleteComputeCapability(ctx context.Context, req *mrdspb.DeleteComputeCapabilityRequest) (*mrdspb.DeleteComputeCapabilityResponse, error) {
	activity.GetLogger(ctx).Info("Deleting ComputeCapability", "request", req)

	resp, err := c.client.Delete(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to delete ComputeCapability", "error", err)
		return nil, fmt.Errorf("failed to delete ComputeCapability: %w", err)
	}

	return resp, nil
}
