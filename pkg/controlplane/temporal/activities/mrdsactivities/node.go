package mrdsactivities

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/gen/api/mrdspb"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type NodeActivities struct {
	client mrdspb.NodesClient
}

// NewNodeActivities creates a new instance of NodeActivities.
func NewNodeActivities(client mrdspb.NodesClient, registry worker.Registry) *NodeActivities {
	a := &NodeActivities{client: client}
	registry.RegisterActivity(a.CreateNode)
	registry.RegisterActivity(a.GetNodeByMetadata)
	registry.RegisterActivity(a.GetNodeByName)
	registry.RegisterActivity(a.UpdateNodeState)
	registry.RegisterActivity(a.ListNode)
	registry.RegisterActivity(a.DeleteNode)
	return a
}

// CreateNode is an activity that interacts with the gRPC service to create a Node.
func (c *NodeActivities) CreateNode(ctx context.Context, req *mrdspb.CreateNodeRequest) (*mrdspb.CreateNodeResponse, error) {
	activity.GetLogger(ctx).Info("Creating Node", "request", req)

	// Check if the context has a deadline to handle timeout.
	if deadline, ok := ctx.Deadline(); ok {
		activity.GetLogger(ctx).Info("Context has a deadline: %v", deadline)
	}

	// Call gRPC method with context for timeout
	resp, err := c.client.Create(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to create Node", "error", err)
		return nil, fmt.Errorf("failed to create Node: %w", err)
	}

	return resp, nil
}

// GetNodeByMetadata fetches Node details based on metadata.
func (c *NodeActivities) GetNodeByMetadata(ctx context.Context, req *mrdspb.GetNodeByMetadataRequest) (*mrdspb.GetNodeResponse, error) {
	activity.GetLogger(ctx).Info("Fetching Node by metadata", "request", req)

	resp, err := c.client.GetByMetadata(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get Node by metadata", "error", err)
		return nil, fmt.Errorf("failed to get Node by metadata: %w", err)
	}

	return resp, nil
}

// GetNodeByName fetches Node details based on name.
func (c *NodeActivities) GetNodeByName(ctx context.Context, req *mrdspb.GetNodeByNameRequest) (*mrdspb.GetNodeResponse, error) {
	activity.GetLogger(ctx).Info("Fetching Node by name", "request", req)

	resp, err := c.client.GetByName(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get Node by name", "error", err)
		return nil, fmt.Errorf("failed to get Node by name: %w", err)
	}

	return resp, nil
}

// UpdateNodeState updates the state of a Node.
func (c *NodeActivities) UpdateNodeState(ctx context.Context, req *mrdspb.UpdateNodeStateRequest) (*mrdspb.UpdateNodeResponse, error) {
	activity.GetLogger(ctx).Info("Updating Node state", "request", req)

	resp, err := c.client.UpdateState(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update Node state", "error", err)
		return nil, fmt.Errorf("failed to update Node state: %w", err)
	}

	return resp, nil
}

// ListNode lists all Nodes based on the request.
func (c *NodeActivities) ListNode(ctx context.Context, req *mrdspb.ListNodeRequest) (*mrdspb.ListNodeResponse, error) {
	activity.GetLogger(ctx).Info("Listing Nodes", "request", req)

	resp, err := c.client.List(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to list Nodes", "error", err)
		return nil, fmt.Errorf("failed to list Nodes: %w", err)
	}

	return resp, nil
}

func (c *NodeActivities) DeleteNode(ctx context.Context, req *mrdspb.DeleteNodeRequest) (*mrdspb.DeleteNodeResponse, error) {
	activity.GetLogger(ctx).Info("Deleting Node", "request", req)

	resp, err := c.client.Delete(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to delete Node", "error", err)
		return nil, fmt.Errorf("failed to delete Node: %w", err)
	}

	return resp, nil
}
