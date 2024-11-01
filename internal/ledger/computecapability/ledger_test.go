package computecapability_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/internal/ledger/computecapability"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"
	"github.com/msanath/mrds/internal/sqlstorage/test"

	"github.com/stretchr/testify/require"
)

func TestLedgerCreate(t *testing.T) {

	t.Run("Create Success", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := computecapability.NewLedger(storage.ComputeCapability)

		req := &computecapability.CreateRequest{
			Name: "test-computecapability",
		}
		resp, err := l.Create(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-computecapability", resp.Record.Name)
		require.NotEmpty(t, resp.Record.Metadata.ID)
		require.Equal(t, uint64(0), resp.Record.Metadata.Version)
		require.Equal(t, computecapability.ComputeCapabilityStatePending, resp.Record.Status.State)
	})

	t.Run("Create EmptyName Failure", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := computecapability.NewLedger(storage.ComputeCapability)

		req := &computecapability.CreateRequest{
			Name: "", // Empty name
		}
		resp, err := l.Create(context.Background(), req)

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})
}

func TestLedgerGetByMetadata(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := computecapability.NewLedger(storage.ComputeCapability)

	req := &computecapability.CreateRequest{
		Name: "test-computecapability",
	}
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByMetadata Success", func(t *testing.T) {
		resp, err := l.GetByID(context.Background(), createResp.Record.Metadata.ID)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-computecapability", resp.Record.Name)
	})

	t.Run("GetByMetadata InvalidID Failure", func(t *testing.T) {
		resp, err := l.GetByID(context.Background(), "")

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})

	t.Run("GetByMetadata NotFound Failure", func(t *testing.T) {
		resp, err := l.GetByID(context.Background(), "unknown")

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})
}

func TestLedgerGetByName(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := computecapability.NewLedger(storage.ComputeCapability)

	req := &computecapability.CreateRequest{
		Name: "test-computecapability",
	}
	_, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByName Success", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "test-computecapability")

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-computecapability", resp.Record.Name)
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
	l := computecapability.NewLedger(storage.ComputeCapability)

	req := &computecapability.CreateRequest{
		Name: "test-computecapability",
	}
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	lastUpdatedRecord := createResp.Record
	t.Run("UpdateStatus Success", func(t *testing.T) {
		updateReq := &computecapability.UpdateStateRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: computecapability.ComputeCapabilityStatus{
				State:   computecapability.ComputeCapabilityStateInActive,
				Message: "ComputeCapability is inactive now",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, computecapability.ComputeCapabilityStateInActive, resp.Record.Status.State)
		lastUpdatedRecord = resp.Record
	})

	t.Run("UpdateStatus InvalidTransition Failure", func(t *testing.T) {
		updateReq := &computecapability.UpdateStateRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: computecapability.ComputeCapabilityStatus{
				State:   computecapability.ComputeCapabilityStatePending, // Invalid transition
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
		updateReq := &computecapability.UpdateStateRequest{
			Metadata: createResp.Record.Metadata, // This is the old metadata. Should cause a conflict.
			Status: computecapability.ComputeCapabilityStatus{
				State:   computecapability.ComputeCapabilityStateActive,
				Message: "ComputeCapability is active now",
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
	l := computecapability.NewLedger(storage.ComputeCapability)

	// Create two ComputeCapabilitys
	_, err := l.Create(context.Background(), &computecapability.CreateRequest{Name: "ComputeCapability1"})
	require.NoError(t, err)

	_, err = l.Create(context.Background(), &computecapability.CreateRequest{Name: "ComputeCapability2"})
	require.NoError(t, err)

	// List the ComputeCapabilitys
	listReq := &computecapability.ListRequest{}
	resp, err := l.List(context.Background(), listReq)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Records, 2)
}

func TestLedgerDelete(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := computecapability.NewLedger(storage.ComputeCapability)

	// First, create the ComputeCapability
	createResp, err := l.Create(context.Background(), &computecapability.CreateRequest{Name: "test-computecapability"})
	require.NoError(t, err)

	// Now, delete the ComputeCapability
	err = l.Delete(context.Background(), &computecapability.DeleteRequest{Metadata: createResp.Record.Metadata})
	require.NoError(t, err)

	// Try to get the ComputeCapability again
	_, err = l.GetByID(context.Background(), createResp.Record.Metadata.ID)
	require.Error(t, err)
	require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
	require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
}
