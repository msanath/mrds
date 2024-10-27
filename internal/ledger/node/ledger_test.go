package node_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/internal/ledger/core"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"
	"github.com/msanath/mrds/internal/ledger/node"
	"github.com/msanath/mrds/internal/sqlstorage/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedgerCreate(t *testing.T) {

	t.Run("Create Success", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := node.NewLedger(storage.Node)

		req := &node.CreateRequest{
			Name:         "test-node",
			UpdateDomain: "test-domain",
			TotalResources: node.Resources{
				Cores:  64,
				Memory: 512,
			},
			SystemReservedResources: node.Resources{
				Cores:  4,
				Memory: 32,
			},
		}
		resp, err := l.Create(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-node", resp.Record.Name)
		require.NotEmpty(t, resp.Record.Metadata.ID)
		require.Equal(t, uint64(0), resp.Record.Metadata.Version)
		require.Equal(t, node.NodeStateUnallocated, resp.Record.Status.State)
	})

	t.Run("Create Failures", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := node.NewLedger(storage.Node)

		testCases := []struct {
			name string
			req  *node.CreateRequest
		}{
			{
				name: "EmptyName",
				req: &node.CreateRequest{
					Name: "", // Empty name
				},
			},
			{
				name: "EmptyUpdateDomain",
				req: &node.CreateRequest{
					Name: "test-node",
					TotalResources: node.Resources{
						Cores:  64,
						Memory: 512,
					},
					SystemReservedResources: node.Resources{
						Cores:  4,
						Memory: 32,
					},
				},
			},
			{
				name: "EmptyTotalResources",
				req: &node.CreateRequest{
					Name:         "test-node",
					UpdateDomain: "test-domain",
					SystemReservedResources: node.Resources{
						Cores:  4,
						Memory: 32,
					},
				},
			},
			{
				name: "SystemReservedResourcesGreaterThanTotalResources",
				req: &node.CreateRequest{
					Name:         "test-node",
					UpdateDomain: "test-domain",
					TotalResources: node.Resources{
						Cores:  64,
						Memory: 512,
					},
					SystemReservedResources: node.Resources{
						Cores:  65,
						Memory: 513,
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				resp, err := l.Create(context.Background(), tc.req)

				require.Error(t, err)
				require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
				require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
				require.Nil(t, resp)
			})
		}
	})
}

func TestLedgerGetByMetadata(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := node.NewLedger(storage.Node)

	req := &node.CreateRequest{
		Name:         "test-node",
		UpdateDomain: "test-domain",
		TotalResources: node.Resources{
			Cores:  64,
			Memory: 512,
		},
		SystemReservedResources: node.Resources{
			Cores:  4,
			Memory: 32,
		},
	}
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByMetadata Success", func(t *testing.T) {
		resp, err := l.GetByMetadata(context.Background(), &createResp.Record.Metadata)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-node", resp.Record.Name)
	})

	t.Run("GetByMetadata InvalidID Failure", func(t *testing.T) {
		resp, err := l.GetByMetadata(context.Background(), &core.Metadata{ID: ""})

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})

	t.Run("GetByMetadata NotFound Failure", func(t *testing.T) {
		resp, err := l.GetByMetadata(context.Background(), &core.Metadata{ID: "unknown"})

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})
}

func TestLedgerGetByName(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := node.NewLedger(storage.Node)

	req := &node.CreateRequest{
		Name:         "test-node",
		UpdateDomain: "test-domain",
		TotalResources: node.Resources{
			Cores:  64,
			Memory: 512,
		},
		SystemReservedResources: node.Resources{
			Cores:  4,
			Memory: 32,
		},
	}
	_, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByName Success", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "test-node")

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-node", resp.Record.Name)
	})

	t.Run("GetByName InvalidName Failure", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "")

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})

	t.Run("GetByName NotFound Failure", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "unknown")

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})
}

func TestLedgerUpdateStatus(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := node.NewLedger(storage.Node)

	req := &node.CreateRequest{
		Name:         "test-node",
		UpdateDomain: "test-domain",
		TotalResources: node.Resources{
			Cores:  64,
			Memory: 512,
		},
		SystemReservedResources: node.Resources{
			Cores:  4,
			Memory: 32,
		},
	}
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	lastUpdatedRecord := createResp.Record
	t.Run("UpdateStatus Success", func(t *testing.T) {
		updateReq := &node.UpdateStatusRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: node.NodeStatus{
				State:   node.NodeStateAllocating,
				Message: "Node is Allocating now",
			},
			ClusterID: "test-cluster",
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, node.NodeStateAllocating, resp.Record.Status.State)
		lastUpdatedRecord = resp.Record
	})

	t.Run("UpdateStatus InvalidTransition Failure", func(t *testing.T) {
		updateReq := &node.UpdateStatusRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: node.NodeStatus{
				State:   node.NodeStateUnallocated, // Invalid transition
				Message: "Invalid state transition",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})

	t.Run("Update conflict Failure", func(t *testing.T) {
		updateReq := &node.UpdateStatusRequest{
			Metadata: createResp.Record.Metadata, // Using an older version.
			Status: node.NodeStatus{
				State:   node.NodeStateAllocating,
				Message: "Node is Allocating now",
			},
			ClusterID: "test-cluster-2",
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordInsertConflict, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})
}

func TestLedgerList(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := node.NewLedger(storage.Node)

	// Create two Nodes
	req := &node.CreateRequest{
		Name:         "test-node",
		UpdateDomain: "test-domain",
		TotalResources: node.Resources{
			Cores:  64,
			Memory: 512,
		},
		SystemReservedResources: node.Resources{
			Cores:  4,
			Memory: 32,
		},
	}
	_, err := l.Create(context.Background(), req)
	assert.NoError(t, err)

	req1 := &node.CreateRequest{
		Name:         "test-node-1",
		UpdateDomain: "test-domain",
		TotalResources: node.Resources{
			Cores:  64,
			Memory: 512,
		},
		SystemReservedResources: node.Resources{
			Cores:  4,
			Memory: 32,
		},
	}
	_, err = l.Create(context.Background(), req1)
	assert.NoError(t, err)

	// List the Nodes
	listReq := &node.ListRequest{}
	resp, err := l.List(context.Background(), listReq)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Records, 2)
}

func TestLedgerDelete(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := node.NewLedger(storage.Node)

	// First, create the Node
	req := &node.CreateRequest{
		Name:         "test-node",
		UpdateDomain: "test-domain",
		TotalResources: node.Resources{
			Cores:  64,
			Memory: 512,
		},
		SystemReservedResources: node.Resources{
			Cores:  4,
			Memory: 32,
		},
	}
	createResp, err := l.Create(context.Background(), req)
	assert.NoError(t, err)

	// Now, delete the Node
	err = l.Delete(context.Background(), &node.DeleteRequest{Metadata: createResp.Record.Metadata})
	assert.NoError(t, err)

	// Try to get the Node again
	_, err = l.GetByMetadata(context.Background(), &createResp.Record.Metadata)
	assert.Error(t, err)
	require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
	require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
}