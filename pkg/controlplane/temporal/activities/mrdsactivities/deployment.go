package mrdsactivities

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/gen/api/mrdspb"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type DeploymentActivities struct {
	client mrdspb.DeploymentsClient
}

// NewDeploymentActivities creates a new instance of DeploymentActivities.
func NewDeploymentActivities(client mrdspb.DeploymentsClient, registry worker.Registry) *DeploymentActivities {
	a := &DeploymentActivities{client: client}
	registry.RegisterActivity(a.CreateDeployment)
	registry.RegisterActivity(a.GetDeploymentByMetadata)
	registry.RegisterActivity(a.GetDeploymentByName)
	registry.RegisterActivity(a.UpdateDeploymentStatus)
	registry.RegisterActivity(a.ListDeployment)
	registry.RegisterActivity(a.DeleteDeployment)
	return a
}

// CreateDeployment is an activity that interacts with the gRPC service to create a Deployment.
func (c *DeploymentActivities) CreateDeployment(ctx context.Context, req *mrdspb.CreateDeploymentRequest) (*mrdspb.CreateDeploymentResponse, error) {
	activity.GetLogger(ctx).Info("Creating Deployment", "request", req)

	// Check if the context has a deadline to handle timeout.
	if deadline, ok := ctx.Deadline(); ok {
		activity.GetLogger(ctx).Info("Context has a deadline: %v", deadline)
	}

	// Call gRPC method with context for timeout
	resp, err := c.client.Create(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to create Deployment", "error", err)
		return nil, fmt.Errorf("failed to create Deployment: %w", err)
	}

	return resp, nil
}

// GetDeploymentByMetadata fetches Deployment details based on metadata.
func (c *DeploymentActivities) GetDeploymentByMetadata(ctx context.Context, req *mrdspb.GetDeploymentByMetadataRequest) (*mrdspb.GetDeploymentResponse, error) {
	activity.GetLogger(ctx).Info("Fetching Deployment by metadata", "request", req)

	resp, err := c.client.GetByMetadata(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get Deployment by metadata", "error", err)
		return nil, fmt.Errorf("failed to get Deployment by metadata: %w", err)
	}

	return resp, nil
}

// GetDeploymentByName fetches Deployment details based on name.
func (c *DeploymentActivities) GetDeploymentByName(ctx context.Context, req *mrdspb.GetDeploymentByNameRequest) (*mrdspb.GetDeploymentResponse, error) {
	activity.GetLogger(ctx).Info("Fetching Deployment by name", "request", req)

	resp, err := c.client.GetByName(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get Deployment by name", "error", err)
		return nil, fmt.Errorf("failed to get Deployment by name: %w", err)
	}

	return resp, nil
}

// UpdateDeploymentState updates the state of a Deployment.
func (c *DeploymentActivities) UpdateDeploymentStatus(ctx context.Context, req *mrdspb.UpdateDeploymentStatusRequest) (*mrdspb.UpdateDeploymentResponse, error) {
	activity.GetLogger(ctx).Info("Updating Deployment state", "request", req)

	resp, err := c.client.UpdateStatus(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update Deployment state", "error", err)
		return nil, fmt.Errorf("failed to update Deployment state: %w", err)
	}

	return resp, nil
}

// ListDeployment lists all Deployments based on the request.
func (c *DeploymentActivities) ListDeployment(ctx context.Context, req *mrdspb.ListDeploymentRequest) (*mrdspb.ListDeploymentResponse, error) {
	activity.GetLogger(ctx).Info("Listing Deployments", "request", req)

	resp, err := c.client.List(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to list Deployments", "error", err)
		return nil, fmt.Errorf("failed to list Deployments: %w", err)
	}

	return resp, nil
}

func (c *DeploymentActivities) DeleteDeployment(ctx context.Context, req *mrdspb.DeleteDeploymentRequest) (*mrdspb.DeleteDeploymentResponse, error) {
	activity.GetLogger(ctx).Info("Deleting Deployment", "request", req)

	resp, err := c.client.Delete(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to delete Deployment", "error", err)
		return nil, fmt.Errorf("failed to delete Deployment: %w", err)
	}

	return resp, nil
}
