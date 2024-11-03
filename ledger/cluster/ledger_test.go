package cluster_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/ledger/cluster"
	ledgererrors "github.com/msanath/mrds/ledger/errors"
	"github.com/msanath/mrds/pkg/sqlstorage/test"

	"github.com/stretchr/testify/require"
)

func TestLedgerCreate(t *testing.T) {

	t.Run("Create Success", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := cluster.NewLedger(storage.Cluster)

		req := &cluster.CreateRequest{
			Name: "test-cluster",
		}
		resp, err := l.Create(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-cluster", resp.Record.Name)
		require.NotEmpty(t, resp.Record.Metadata.ID)
		require.Equal(t, uint64(0), resp.Record.Metadata.Version)
		require.Equal(t, cluster.ClusterStatePending, resp.Record.Status.State)
	})

	t.Run("Create EmptyName Failure", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := cluster.NewLedger(storage.Cluster)

		req := &cluster.CreateRequest{
			Name: "", // Empty name
		}
		resp, err := l.Create(context.Background(), req)

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})
}

func TestLedgerGetByID(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := cluster.NewLedger(storage.Cluster)

	req := &cluster.CreateRequest{
		Name: "test-cluster",
	}
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByID Success", func(t *testing.T) {
		resp, err := l.GetByID(context.Background(), createResp.Record.Metadata.ID)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-cluster", resp.Record.Name)
	})

	t.Run("GetByID InvalidID Failure", func(t *testing.T) {
		resp, err := l.GetByID(context.Background(), "")

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})

	t.Run("GetByID NotFound Failure", func(t *testing.T) {
		resp, err := l.GetByID(context.Background(), "unknown")

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})
}

func TestLedgerGetByName(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := cluster.NewLedger(storage.Cluster)

	req := &cluster.CreateRequest{
		Name: "test-cluster",
	}
	_, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByName Success", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "test-cluster")

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-cluster", resp.Record.Name)
	})

	t.Run("GetByName InvalidName Failure", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "")

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})

	t.Run("GetByName NotFound Failure", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "unknown")

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})
}

func TestLedgerUpdateStatus(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := cluster.NewLedger(storage.Cluster)

	req := &cluster.CreateRequest{
		Name: "test-cluster",
	}
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	lastUpdatedRecord := createResp.Record
	t.Run("UpdateStatus Success", func(t *testing.T) {
		updateReq := &cluster.UpdateStateRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: cluster.ClusterStatus{
				State:   cluster.ClusterStateInActive,
				Message: "Cluster is inactive now",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, cluster.ClusterStateInActive, resp.Record.Status.State)
		lastUpdatedRecord = resp.Record
	})

	t.Run("UpdateStatus InvalidTransition Failure", func(t *testing.T) {
		updateReq := &cluster.UpdateStateRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: cluster.ClusterStatus{
				State:   cluster.ClusterStatePending, // Invalid transition
				Message: "Invalid state transition",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})

	t.Run("Update conflict Failure", func(t *testing.T) {
		updateReq := &cluster.UpdateStateRequest{
			Metadata: createResp.Record.Metadata, // This is the old metadata. Should cause a conflict.
			Status: cluster.ClusterStatus{
				State:   cluster.ClusterStateActive,
				Message: "Cluster is active now",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordInsertConflict, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})
}

func TestLedgerList(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := cluster.NewLedger(storage.Cluster)

	// Create two Clusters
	_, err := l.Create(context.Background(), &cluster.CreateRequest{Name: "Cluster1"})
	require.NoError(t, err)

	_, err = l.Create(context.Background(), &cluster.CreateRequest{Name: "Cluster2"})
	require.NoError(t, err)

	// List the Clusters
	listReq := &cluster.ListRequest{}
	resp, err := l.List(context.Background(), listReq)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Records, 2)
}

func TestLedgerDelete(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := cluster.NewLedger(storage.Cluster)

	// First, create the Cluster
	createResp, err := l.Create(context.Background(), &cluster.CreateRequest{Name: "test-cluster"})
	require.NoError(t, err)

	// Now, delete the Cluster
	err = l.Delete(context.Background(), &cluster.DeleteRequest{Metadata: createResp.Record.Metadata})
	require.NoError(t, err)

	// Try to get the Cluster again
	_, err = l.GetByID(context.Background(), createResp.Record.Metadata.ID)
	require.Error(t, err)
	require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
	require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
}
