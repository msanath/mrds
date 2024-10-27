package grpcservers_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/gen/api/mrdspb"
	servertest "github.com/msanath/mrds/pkg/grpcservers/test"

	"github.com/stretchr/testify/require"
)

func TestDeploymentServer(t *testing.T) {
	ts, err := servertest.NewTestServer()
	require.NoError(t, err)
	defer ts.Close()

	client := mrdspb.NewDeploymentsClient(ts.Conn())
	ctx := context.Background()

	// create
	resp, err := client.Create(ctx, &mrdspb.CreateDeploymentRequest{Name: "test-Deployment"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "test-Deployment", resp.Record.Name)

	// get by metadata
	getResp, err := client.GetByMetadata(ctx, &mrdspb.GetDeploymentByMetadataRequest{
		Metadata: resp.Record.Metadata,
	})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, "test-Deployment", getResp.Record.Name)

	// get by name
	getByNameResp, err := client.GetByName(ctx, &mrdspb.GetDeploymentByNameRequest{Name: "test-Deployment"})
	require.NoError(t, err)
	require.NotNil(t, getByNameResp)
	require.Equal(t, "test-Deployment", getByNameResp.Record.Name)

	// update
	updateResp, err := client.UpdateStatus(ctx, &mrdspb.UpdateDeploymentStatusRequest{
		Metadata: resp.Record.Metadata,
		Status: &mrdspb.DeploymentStatus{
			State:   mrdspb.DeploymentState_DeploymentState_ACTIVE,
			Message: "test-message",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-Deployment", updateResp.Record.Name)
	require.Equal(t, mrdspb.DeploymentState_DeploymentState_ACTIVE, updateResp.Record.Status.State)

	// Create another
	resp2, err := client.Create(ctx, &mrdspb.CreateDeploymentRequest{Name: "test-Deployment-2"})
	require.NoError(t, err)
	require.NotNil(t, resp2)

	// list
	listResp, err := client.List(ctx, &mrdspb.ListDeploymentRequest{
		StateIn: []mrdspb.DeploymentState{mrdspb.DeploymentState_DeploymentState_ACTIVE},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Records, 1)

	// Delete
	_, err = client.Delete(ctx, &mrdspb.DeleteDeploymentRequest{Metadata: resp2.Record.Metadata})
	require.NoError(t, err)

	// Get deleted by name
	_, err = client.GetByName(ctx, &mrdspb.GetDeploymentByNameRequest{Name: "test-Deployment-2"})
	require.Error(t, err)
}
