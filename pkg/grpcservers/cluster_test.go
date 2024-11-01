package grpcservers_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/gen/api/mrdspb"
	servertest "github.com/msanath/mrds/pkg/grpcservers/test"

	"github.com/stretchr/testify/require"
)

func TestClusterServer(t *testing.T) {
	ts, err := servertest.NewTestServer()
	require.NoError(t, err)
	defer ts.Close()

	client := mrdspb.NewClustersClient(ts.Conn())
	ctx := context.Background()

	// create
	resp, err := client.Create(ctx, &mrdspb.CreateClusterRequest{Name: "test-Cluster"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "test-Cluster", resp.Record.Name)

	// get by id
	getResp, err := client.GetByID(ctx, &mrdspb.GetClusterByIDRequest{Id: resp.Record.Metadata.Id})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, "test-Cluster", getResp.Record.Name)

	// get by name
	getByNameResp, err := client.GetByName(ctx, &mrdspb.GetClusterByNameRequest{Name: "test-Cluster"})
	require.NoError(t, err)
	require.NotNil(t, getByNameResp)
	require.Equal(t, "test-Cluster", getByNameResp.Record.Name)

	// update
	updateResp, err := client.UpdateStatus(ctx, &mrdspb.UpdateClusterStatusRequest{
		Metadata: resp.Record.Metadata,
		Status: &mrdspb.ClusterStatus{
			State:   mrdspb.ClusterState_ClusterState_ACTIVE,
			Message: "test-message",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-Cluster", updateResp.Record.Name)
	require.Equal(t, mrdspb.ClusterState_ClusterState_ACTIVE, updateResp.Record.Status.State)

	// Create another
	resp2, err := client.Create(ctx, &mrdspb.CreateClusterRequest{Name: "test-Cluster-2"})
	require.NoError(t, err)
	require.NotNil(t, resp2)

	// list
	listResp, err := client.List(ctx, &mrdspb.ListClusterRequest{
		StateIn: []mrdspb.ClusterState{mrdspb.ClusterState_ClusterState_ACTIVE},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Records, 1)

	// Delete
	_, err = client.Delete(ctx, &mrdspb.DeleteClusterRequest{Metadata: resp2.Record.Metadata})
	require.NoError(t, err)

	// Get deleted by name
	_, err = client.GetByName(ctx, &mrdspb.GetClusterByNameRequest{Name: "test-Cluster-2"})
	require.Error(t, err)
}
