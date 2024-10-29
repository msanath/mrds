package mrdsactivities

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/gen/api/mrdspb"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type DeploymentPlanActivities struct {
	client mrdspb.DeploymentPlansClient
}

// NewDeploymentPlanActivities creates a new instance of DeploymentPlanActivities.
func NewDeploymentPlanActivities(client mrdspb.DeploymentPlansClient, registry worker.Registry) *DeploymentPlanActivities {
	a := &DeploymentPlanActivities{client: client}
	registry.RegisterActivity(a.CreateDeploymentPlan)
	registry.RegisterActivity(a.GetDeploymentPlanByMetadata)
	registry.RegisterActivity(a.GetDeploymentPlanByName)
	registry.RegisterActivity(a.UpdateDeploymentPlanStatus)
	registry.RegisterActivity(a.ListDeploymentPlan)
	registry.RegisterActivity(a.DeleteDeploymentPlan)
	return a
}

// CreateDeploymentPlan is an activity that interacts with the gRPC service to create a DeploymentPlan.
func (c *DeploymentPlanActivities) CreateDeploymentPlan(ctx context.Context, req *mrdspb.CreateDeploymentPlanRequest) (*mrdspb.CreateDeploymentPlanResponse, error) {
	activity.GetLogger(ctx).Info("Creating DeploymentPlan", "request", req)

	// Check if the context has a deadline to handle timeout.
	if deadline, ok := ctx.Deadline(); ok {
		activity.GetLogger(ctx).Info("Context has a deadline: %v", deadline)
	}

	// Call gRPC method with context for timeout
	resp, err := c.client.Create(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to create DeploymentPlan", "error", err)
		return nil, fmt.Errorf("failed to create DeploymentPlan: %w", err)
	}

	return resp, nil
}

// GetDeploymentPlanByMetadata fetches DeploymentPlan details based on metadata.
func (c *DeploymentPlanActivities) GetDeploymentPlanByMetadata(ctx context.Context, req *mrdspb.GetDeploymentPlanByMetadataRequest) (*mrdspb.GetDeploymentPlanResponse, error) {
	activity.GetLogger(ctx).Info("Fetching DeploymentPlan by metadata", "request", req)

	resp, err := c.client.GetByMetadata(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get DeploymentPlan by metadata", "error", err)
		return nil, fmt.Errorf("failed to get DeploymentPlan by metadata: %w", err)
	}

	return resp, nil
}

// GetDeploymentPlanByName fetches DeploymentPlan details based on name.
func (c *DeploymentPlanActivities) GetDeploymentPlanByName(ctx context.Context, req *mrdspb.GetDeploymentPlanByNameRequest) (*mrdspb.GetDeploymentPlanResponse, error) {
	activity.GetLogger(ctx).Info("Fetching DeploymentPlan by name", "request", req)

	resp, err := c.client.GetByName(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get DeploymentPlan by name", "error", err)
		return nil, fmt.Errorf("failed to get DeploymentPlan by name: %w", err)
	}

	return resp, nil
}

// UpdateDeploymentPlanState updates the state of a DeploymentPlan.
func (c *DeploymentPlanActivities) UpdateDeploymentPlanStatus(ctx context.Context, req *mrdspb.UpdateDeploymentPlanStatusRequest) (*mrdspb.UpdateDeploymentPlanResponse, error) {
	activity.GetLogger(ctx).Info("Updating DeploymentPlan state", "request", req)

	resp, err := c.client.UpdateStatus(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update DeploymentPlan state", "error", err)
		return nil, fmt.Errorf("failed to update DeploymentPlan state: %w", err)
	}

	return resp, nil
}

// ListDeploymentPlan lists all DeploymentPlans based on the request.
func (c *DeploymentPlanActivities) ListDeploymentPlan(ctx context.Context, req *mrdspb.ListDeploymentPlanRequest) (*mrdspb.ListDeploymentPlanResponse, error) {
	activity.GetLogger(ctx).Info("Listing DeploymentPlans", "request", req)

	resp, err := c.client.List(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to list DeploymentPlans", "error", err)
		return nil, fmt.Errorf("failed to list DeploymentPlans: %w", err)
	}

	return resp, nil
}

func (c *DeploymentPlanActivities) DeleteDeploymentPlan(ctx context.Context, req *mrdspb.DeleteDeploymentPlanRequest) (*mrdspb.DeleteDeploymentPlanResponse, error) {
	activity.GetLogger(ctx).Info("Deleting DeploymentPlan", "request", req)

	resp, err := c.client.Delete(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to delete DeploymentPlan", "error", err)
		return nil, fmt.Errorf("failed to delete DeploymentPlan: %w", err)
	}

	return resp, nil
}

// AddDeployment adds a deployment to a DeploymentPlan.
func (c *DeploymentPlanActivities) AddDeployment(ctx context.Context, req *mrdspb.AddDeploymentRequest) (*mrdspb.UpdateDeploymentPlanResponse, error) {
	activity.GetLogger(ctx).Info("Adding Deployment to DeploymentPlan", "request", req)

	resp, err := c.client.AddDeployment(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to add Deployment to DeploymentPlan", "error", err)
		return nil, fmt.Errorf("failed to add Deployment to DeploymentPlan: %w", err)
	}

	return resp, nil
}

// UpdateDeploymentStatus updates the status of a deployment in a DeploymentPlan.
func (c *DeploymentPlanActivities) UpdateDeploymentStatus(ctx context.Context, req *mrdspb.UpdateDeploymentStatusRequest) (*mrdspb.UpdateDeploymentPlanResponse, error) {
	activity.GetLogger(ctx).Info("Updating Deployment status in DeploymentPlan", "request", req)

	resp, err := c.client.UpdateDeploymentStatus(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update Deployment status in DeploymentPlan", "error", err)
		return nil, fmt.Errorf("failed to update Deployment status in DeploymentPlan: %w", err)
	}

	return resp, nil
}
