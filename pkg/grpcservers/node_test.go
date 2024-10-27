package grpcservers_test

import (
	"context"
	"testing"

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
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "test-Node", resp.Record.Name)

	// get by metadata
	getResp, err := client.GetByMetadata(ctx, &mrdspb.GetNodeByMetadataRequest{
		Metadata: resp.Record.Metadata,
	})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, "test-Node", getResp.Record.Name)

	// get by name
	getByNameResp, err := client.GetByName(ctx, &mrdspb.GetNodeByNameRequest{Name: "test-Node"})
	require.NoError(t, err)
	require.NotNil(t, getByNameResp)
	require.Equal(t, "test-Node", getByNameResp.Record.Name)

	// update
	updateResp, err := client.UpdateState(ctx, &mrdspb.UpdateNodeStateRequest{
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
