
package grpcservers_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/gen/api/mrdspb"
	servertest "github.com/msanath/mrds/pkg/grpcservers/test"

	"github.com/stretchr/testify/require"
)

func TestDeploymentPlanServer(t *testing.T) {
	ts, err := servertest.NewTestServer()
	require.NoError(t, err)
	defer ts.Close()

	client := mrdspb.NewDeploymentPlansClient(ts.Conn())
	ctx := context.Background()

	// create
	resp, err := client.Create(ctx, &mrdspb.CreateDeploymentPlanRequest{Name: "test-DeploymentPlan"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "test-DeploymentPlan", resp.Record.Name)

	// get by metadata
	getResp, err := client.GetByMetadata(ctx, &mrdspb.GetDeploymentPlanByMetadataRequest{
		Metadata: resp.Record.Metadata,
	})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, "test-DeploymentPlan", getResp.Record.Name)

	// get by name
	getByNameResp, err := client.GetByName(ctx, &mrdspb.GetDeploymentPlanByNameRequest{Name: "test-DeploymentPlan"})
	require.NoError(t, err)
	require.NotNil(t, getByNameResp)
	require.Equal(t, "test-DeploymentPlan", getByNameResp.Record.Name)

	// update
	updateResp, err := client.UpdateStatus(ctx, &mrdspb.UpdateDeploymentPlanStatusRequest{
		Metadata: resp.Record.Metadata,
		Status: &mrdspb.DeploymentPlanStatus{
			State:   mrdspb.DeploymentPlanState_DeploymentPlanState_ACTIVE,
			Message: "test-message",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-DeploymentPlan", updateResp.Record.Name)
	require.Equal(t, mrdspb.DeploymentPlanState_DeploymentPlanState_ACTIVE, updateResp.Record.Status.State)

	// Create another
	resp2, err := client.Create(ctx, &mrdspb.CreateDeploymentPlanRequest{Name: "test-DeploymentPlan-2"})
	require.NoError(t, err)
	require.NotNil(t, resp2)

	// list
	listResp, err := client.List(ctx, &mrdspb.ListDeploymentPlanRequest{
		StateIn: []mrdspb.DeploymentPlanState{mrdspb.DeploymentPlanState_DeploymentPlanState_ACTIVE},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Records, 1)

	// Delete
	_, err = client.Delete(ctx, &mrdspb.DeleteDeploymentPlanRequest{Metadata: resp2.Record.Metadata})
	require.NoError(t, err)

	// Get deleted by name
	_, err = client.GetByName(ctx, &mrdspb.GetDeploymentPlanByNameRequest{Name: "test-DeploymentPlan-2"})
	require.Error(t, err)
}
