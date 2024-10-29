
package deploymentplan_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/internal/ledger/deploymentplan"
	"github.com/msanath/mrds/internal/ledger/core"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"
	"github.com/msanath/mrds/internal/sqlstorage/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedgerCreate(t *testing.T) {

	t.Run("Create Success", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := deploymentplan.NewLedger(storage.DeploymentPlan)

		req := &deploymentplan.CreateRequest{
			Name: "test-deploymentplan",
		}
		resp, err := l.Create(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-deploymentplan", resp.Record.Name)
		require.NotEmpty(t, resp.Record.Metadata.ID)
		require.Equal(t, uint64(0), resp.Record.Metadata.Version)
		require.Equal(t, deploymentplan.DeploymentPlanStatePending, resp.Record.Status.State)
	})

	t.Run("Create EmptyName Failure", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := deploymentplan.NewLedger(storage.DeploymentPlan)

		req := &deploymentplan.CreateRequest{
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
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	req := &deploymentplan.CreateRequest{
		Name: "test-deploymentplan",
	}
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByMetadata Success", func(t *testing.T) {
		resp, err := l.GetByMetadata(context.Background(), &createResp.Record.Metadata)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-deploymentplan", resp.Record.Name)
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
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	req := &deploymentplan.CreateRequest{
		Name: "test-deploymentplan",
	}
	_, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByName Success", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "test-deploymentplan")

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-deploymentplan", resp.Record.Name)
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
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	req := &deploymentplan.CreateRequest{
		Name: "test-deploymentplan",
	}
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	lastUpdatedRecord := createResp.Record
	t.Run("UpdateStatus Success", func(t *testing.T) {
		updateReq := &deploymentplan.UpdateStatusRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: deploymentplan.DeploymentPlanStatus{
				State:   deploymentplan.DeploymentPlanStateInActive,
				Message: "DeploymentPlan is inactive now",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, deploymentplan.DeploymentPlanStateInActive, resp.Record.Status.State)
		lastUpdatedRecord = resp.Record
	})

	t.Run("UpdateStatus InvalidTransition Failure", func(t *testing.T) {
		updateReq := &deploymentplan.UpdateStatusRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: deploymentplan.DeploymentPlanStatus{
				State:   deploymentplan.DeploymentPlanStatePending, // Invalid transition
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
		updateReq := &deploymentplan.UpdateStatusRequest{
			Metadata: createResp.Record.Metadata, // This is the old metadata. Should cause a conflict.
			Status: deploymentplan.DeploymentPlanStatus{
				State:   deploymentplan.DeploymentPlanStateActive,
				Message: "DeploymentPlan is active now",
			},
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
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	// Create two DeploymentPlans
	_, err := l.Create(context.Background(), &deploymentplan.CreateRequest{Name: "DeploymentPlan1"})
	assert.NoError(t, err)

	_, err = l.Create(context.Background(), &deploymentplan.CreateRequest{Name: "DeploymentPlan2"})
	assert.NoError(t, err)

	// List the DeploymentPlans
	listReq := &deploymentplan.ListRequest{}
	resp, err := l.List(context.Background(), listReq)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Records, 2)
}

func TestLedgerDelete(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	// First, create the DeploymentPlan
	createResp, err := l.Create(context.Background(), &deploymentplan.CreateRequest{Name: "test-deploymentplan"})
	assert.NoError(t, err)

	// Now, delete the DeploymentPlan
	err = l.Delete(context.Background(), &deploymentplan.DeleteRequest{Metadata: createResp.Record.Metadata})
	assert.NoError(t, err)

	// Try to get the DeploymentPlan again
	_, err = l.GetByMetadata(context.Background(), &createResp.Record.Metadata)
	assert.Error(t, err)
	require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
	require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
}
