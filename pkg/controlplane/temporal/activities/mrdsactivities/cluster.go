package mrdsactivities

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/gen/api/mrdspb"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type ClusterActivities struct {
	client mrdspb.ClustersClient
}

// NewClusterActivities creates a new instance of ClusterActivities.
func NewClusterActivities(client mrdspb.ClustersClient, registry worker.Registry) *ClusterActivities {
	a := &ClusterActivities{client: client}
	registry.RegisterActivity(a.CreateCluster)
	registry.RegisterActivity(a.GetClusterByMetadata)
	registry.RegisterActivity(a.GetClusterByName)
	registry.RegisterActivity(a.UpdateClusterState)
	registry.RegisterActivity(a.ListCluster)
	registry.RegisterActivity(a.DeleteCluster)
	return a
}

// CreateCluster is an activity that interacts with the gRPC service to create a Cluster.
func (c *ClusterActivities) CreateCluster(ctx context.Context, req *mrdspb.CreateClusterRequest) (*mrdspb.CreateClusterResponse, error) {
	activity.GetLogger(ctx).Info("Creating Cluster", "request", req)

	// Check if the context has a deadline to handle timeout.
	if deadline, ok := ctx.Deadline(); ok {
		activity.GetLogger(ctx).Info("Context has a deadline: %v", deadline)
	}

	// Call gRPC method with context for timeout
	resp, err := c.client.Create(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to create Cluster", "error", err)
		return nil, fmt.Errorf("failed to create Cluster: %w", err)
	}

	return resp, nil
}

// GetClusterByMetadata fetches Cluster details based on metadata.
func (c *ClusterActivities) GetClusterByMetadata(ctx context.Context, req *mrdspb.GetClusterByMetadataRequest) (*mrdspb.GetClusterResponse, error) {
	activity.GetLogger(ctx).Info("Fetching Cluster by metadata", "request", req)

	resp, err := c.client.GetByMetadata(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get Cluster by metadata", "error", err)
		return nil, fmt.Errorf("failed to get Cluster by metadata: %w", err)
	}

	return resp, nil
}

// GetClusterByName fetches Cluster details based on name.
func (c *ClusterActivities) GetClusterByName(ctx context.Context, req *mrdspb.GetClusterByNameRequest) (*mrdspb.GetClusterResponse, error) {
	activity.GetLogger(ctx).Info("Fetching Cluster by name", "request", req)

	resp, err := c.client.GetByName(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get Cluster by name", "error", err)
		return nil, fmt.Errorf("failed to get Cluster by name: %w", err)
	}

	return resp, nil
}

// UpdateClusterState updates the state of a Cluster.
func (c *ClusterActivities) UpdateClusterState(ctx context.Context, req *mrdspb.UpdateClusterStateRequest) (*mrdspb.UpdateClusterResponse, error) {
	activity.GetLogger(ctx).Info("Updating Cluster state", "request", req)

	resp, err := c.client.UpdateState(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update Cluster state", "error", err)
		return nil, fmt.Errorf("failed to update Cluster state: %w", err)
	}

	return resp, nil
}

// ListCluster lists all Clusters based on the request.
func (c *ClusterActivities) ListCluster(ctx context.Context, req *mrdspb.ListClusterRequest) (*mrdspb.ListClusterResponse, error) {
	activity.GetLogger(ctx).Info("Listing Clusters", "request", req)

	resp, err := c.client.List(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to list Clusters", "error", err)
		return nil, fmt.Errorf("failed to list Clusters: %w", err)
	}

	return resp, nil
}

func (c *ClusterActivities) DeleteCluster(ctx context.Context, req *mrdspb.DeleteClusterRequest) (*mrdspb.DeleteClusterResponse, error) {
	activity.GetLogger(ctx).Info("Deleting Cluster", "request", req)

	resp, err := c.client.Delete(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to delete Cluster", "error", err)
		return nil, fmt.Errorf("failed to delete Cluster: %w", err)
	}

	return resp, nil
}
