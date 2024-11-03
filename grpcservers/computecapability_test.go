package grpcservers_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/gen/api/mrdspb"
	testserver "github.com/msanath/mrds/test/server"

	"github.com/stretchr/testify/require"
)

func TestComputeCapabilityServer(t *testing.T) {
	ts, err := testserver.NewTestServer()
	require.NoError(t, err)
	defer ts.Close()

	client := mrdspb.NewComputeCapabilitiesClient(ts.Conn())
	ctx := context.Background()

	// create
	resp, err := client.Create(ctx, &mrdspb.CreateComputeCapabilityRequest{
		Name:  "test-ComputeCapability",
		Type:  "CPU",
		Score: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "test-ComputeCapability", resp.Record.Name)
	require.Equal(t, "CPU", resp.Record.Type)
	require.Equal(t, uint32(10), resp.Record.Score)

	// get by id
	getResp, err := client.GetByID(ctx, &mrdspb.GetComputeCapabilityByIDRequest{Id: resp.Record.Metadata.Id})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, "test-ComputeCapability", getResp.Record.Name)

	// get by name
	getByNameResp, err := client.GetByName(ctx, &mrdspb.GetComputeCapabilityByNameRequest{Name: "test-ComputeCapability"})
	require.NoError(t, err)
	require.NotNil(t, getByNameResp)
	require.Equal(t, "test-ComputeCapability", getByNameResp.Record.Name)

	// update
	updateResp, err := client.UpdateStatus(ctx, &mrdspb.UpdateComputeCapabilityStatusRequest{
		Metadata: resp.Record.Metadata,
		Status: &mrdspb.ComputeCapabilityStatus{
			State:   mrdspb.ComputeCapabilityState_ComputeCapabilityState_ACTIVE,
			Message: "test-message",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-ComputeCapability", updateResp.Record.Name)
	require.Equal(t, mrdspb.ComputeCapabilityState_ComputeCapabilityState_ACTIVE, updateResp.Record.Status.State)

	// Create another
	resp2, err := client.Create(ctx, &mrdspb.CreateComputeCapabilityRequest{Name: "test-ComputeCapability-2"})
	require.NoError(t, err)
	require.NotNil(t, resp2)

	// list
	listResp, err := client.List(ctx, &mrdspb.ListComputeCapabilityRequest{
		StateIn: []mrdspb.ComputeCapabilityState{mrdspb.ComputeCapabilityState_ComputeCapabilityState_ACTIVE},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Records, 1)

	// Delete
	_, err = client.Delete(ctx, &mrdspb.DeleteComputeCapabilityRequest{Metadata: resp2.Record.Metadata})
	require.NoError(t, err)

	// Get deleted by name
	_, err = client.GetByName(ctx, &mrdspb.GetComputeCapabilityByNameRequest{Name: "test-ComputeCapability-2"})
	require.Error(t, err)
}
