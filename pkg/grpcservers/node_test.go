package grpcservers_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/msanath/mrds/gen/api/mrdspb"
	servertest "github.com/msanath/mrds/pkg/grpcservers/test"

	"github.com/stretchr/testify/require"
)

func TestNodeServer(t *testing.T) {
	ts, err := servertest.NewTestServer()
	require.NoError(t, err)
	defer ts.Close()

	client := mrdspb.NewNodesClient(ts.Conn())
	ctx := context.Background()

	// create
	resp, err := client.Create(ctx, &mrdspb.CreateNodeRequest{
		Name:         "test-Node",
		UpdateDomain: "test-domain",
		TotalResources: &mrdspb.Resources{
			Cores:  64,
			Memory: 512,
		},
		SystemReservedResources: &mrdspb.Resources{
			Cores:  4,
			Memory: 32,
		},
		LocalVolumes: []*mrdspb.NodeLocalVolume{
			{
				MountPath:       "/var/lib/docker",
				StorageClass:    "SSD",
				StorageCapacity: 1024,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "test-Node", resp.Record.Name)
	require.Len(t, resp.Record.LocalVolumes, 1)

	// get by id
	getResp, err := client.GetByID(ctx, &mrdspb.GetNodeByIDRequest{Id: resp.Record.Metadata.Id})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, "test-Node", getResp.Record.Name)

	// get by name
	getByNameResp, err := client.GetByName(ctx, &mrdspb.GetNodeByNameRequest{Name: "test-Node"})
	require.NoError(t, err)
	require.NotNil(t, getByNameResp)
	require.Equal(t, "test-Node", getByNameResp.Record.Name)

	// update
	updateResp, err := client.UpdateStatus(ctx, &mrdspb.UpdateNodeStatusRequest{
		Metadata: resp.Record.Metadata,
		Status: &mrdspb.NodeStatus{
			State:   mrdspb.NodeState_NodeState_ALLOCATING,
			Message: "test-message",
		},
		ClusterId: "test-cluster",
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-Node", updateResp.Record.Name)
	require.Equal(t, mrdspb.NodeState_NodeState_ALLOCATING, updateResp.Record.Status.State)

	// Add disrurption
	updateResp, err = client.AddDisruption(ctx, &mrdspb.AddDisruptionRequest{
		Metadata: updateResp.Record.Metadata,
		Disruption: &mrdspb.NodeDisruption{
			Id:          uuid.New().String(),
			ShouldEvict: false,
			Status: &mrdspb.DisruptionStatus{
				State:   mrdspb.DisruptionState_DisruptionState_SCHEDULED,
				Message: "test-message",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-Node", updateResp.Record.Name)
	require.Len(t, updateResp.Record.Disruptions, 1)

	// Update disrurption status
	updateResp, err = client.UpdateDisruptionStatus(ctx, &mrdspb.UpdateDisruptionStatusRequest{
		Metadata:     updateResp.Record.Metadata,
		DisruptionId: updateResp.Record.Disruptions[0].Id,
		Status: &mrdspb.DisruptionStatus{
			State:   mrdspb.DisruptionState_DisruptionState_COMPLETED,
			Message: "test-message",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-Node", updateResp.Record.Name)
	require.Equal(t, mrdspb.DisruptionState_DisruptionState_COMPLETED, updateResp.Record.Disruptions[0].Status.State)

	// Remove disrurption
	updateResp, err = client.RemoveDisruption(ctx, &mrdspb.RemoveDisruptionRequest{
		Metadata:     updateResp.Record.Metadata,
		DisruptionId: updateResp.Record.Disruptions[0].Id,
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-Node", updateResp.Record.Name)
	require.Len(t, updateResp.Record.Disruptions, 0)

	// Add capability
	updateResp, err = client.AddCapability(ctx, &mrdspb.AddCapabilityRequest{
		Metadata:     updateResp.Record.Metadata,
		CapabilityId: "test-capability",
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-Node", updateResp.Record.Name)
	require.Len(t, updateResp.Record.CapabilityIds, 1)
	require.Equal(t, "test-capability", updateResp.Record.CapabilityIds[0])

	// Remove capability
	updateResp, err = client.RemoveCapability(ctx, &mrdspb.RemoveCapabilityRequest{
		Metadata:     updateResp.Record.Metadata,
		CapabilityId: "test-capability",
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-Node", updateResp.Record.Name)
	require.Len(t, updateResp.Record.CapabilityIds, 0)

	// Create another
	resp2, err := client.Create(ctx, &mrdspb.CreateNodeRequest{
		Name:         "test-Node-2",
		UpdateDomain: "test-domain",
		TotalResources: &mrdspb.Resources{
			Cores:  64,
			Memory: 512,
		},
		SystemReservedResources: &mrdspb.Resources{
			Cores:  4,
			Memory: 32,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp2)

	// list
	listResp, err := client.List(ctx, &mrdspb.ListNodeRequest{
		StateIn: []mrdspb.NodeState{mrdspb.NodeState_NodeState_ALLOCATING},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Records, 1)

	// Delete
	_, err = client.Delete(ctx, &mrdspb.DeleteNodeRequest{Metadata: resp2.Record.Metadata})
	require.NoError(t, err)

	// Get deleted by name
	_, err = client.GetByName(ctx, &mrdspb.GetNodeByNameRequest{Name: "test-Node-2"})
	require.Error(t, err)
}
