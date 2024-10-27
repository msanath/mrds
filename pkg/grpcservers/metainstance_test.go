package grpcservers_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/gen/api/mrdspb"
	servertest "github.com/msanath/mrds/pkg/grpcservers/test"

	"github.com/stretchr/testify/require"
)

func TestMetaInstanceServer(t *testing.T) {
	ts, err := servertest.NewTestServer()
	require.NoError(t, err)
	defer ts.Close()

	client := mrdspb.NewMetaInstancesClient(ts.Conn())
	ctx := context.Background()

	// create
	resp, err := client.Create(ctx, &mrdspb.CreateMetaInstanceRequest{Name: "test-MetaInstance"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "test-MetaInstance", resp.Record.Name)

	// get by metadata
	getResp, err := client.GetByMetadata(ctx, &mrdspb.GetMetaInstanceByMetadataRequest{
		Metadata: resp.Record.Metadata,
	})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, "test-MetaInstance", getResp.Record.Name)

	// get by name
	getByNameResp, err := client.GetByName(ctx, &mrdspb.GetMetaInstanceByNameRequest{Name: "test-MetaInstance"})
	require.NoError(t, err)
	require.NotNil(t, getByNameResp)
	require.Equal(t, "test-MetaInstance", getByNameResp.Record.Name)

	// update
	updateResp, err := client.UpdateStatus(ctx, &mrdspb.UpdateMetaInstanceStatusRequest{
		Metadata: resp.Record.Metadata,
		Status: &mrdspb.MetaInstanceStatus{
			State:   mrdspb.MetaInstanceState_MetaInstanceState_ACTIVE,
			Message: "test-message",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-MetaInstance", updateResp.Record.Name)
	require.Equal(t, mrdspb.MetaInstanceState_MetaInstanceState_ACTIVE, updateResp.Record.Status.State)

	// Create another
	resp2, err := client.Create(ctx, &mrdspb.CreateMetaInstanceRequest{Name: "test-MetaInstance-2"})
	require.NoError(t, err)
	require.NotNil(t, resp2)

	// list
	listResp, err := client.List(ctx, &mrdspb.ListMetaInstanceRequest{
		StateIn: []mrdspb.MetaInstanceState{mrdspb.MetaInstanceState_MetaInstanceState_ACTIVE},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Records, 1)

	// Delete
	_, err = client.Delete(ctx, &mrdspb.DeleteMetaInstanceRequest{Metadata: resp2.Record.Metadata})
	require.NoError(t, err)

	// Get deleted by name
	_, err = client.GetByName(ctx, &mrdspb.GetMetaInstanceByNameRequest{Name: "test-MetaInstance-2"})
	require.Error(t, err)
}
