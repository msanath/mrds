package deploymentplan_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/internal/ledger/deploymentplan"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"
	"github.com/msanath/mrds/internal/sqlstorage/test"

	"github.com/stretchr/testify/require"
)

func deplomentRequest() *deploymentplan.CreateRequest {
	return &deploymentplan.CreateRequest{
		Name:        "test-deploymentplan",
		Namespace:   "test-namespace",
		ServiceName: "test-service",
		Applications: []deploymentplan.Application{
			{
				PayloadName: "test-payload",
				Resources: deploymentplan.ApplicationResources{
					Cores: 1,
				},
			},
		},
		MatchingComputeCapabilities: []deploymentplan.MatchingComputeCapability{
			{
				CapabilityType:  "test-capability",
				Comparator:      deploymentplan.ComparatorTypeIn,
				CapabilityNames: []string{"test-capability-name"},
			},
		},
	}
}

func TestLedgerCreate(t *testing.T) {

	t.Run("Create Success", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := deploymentplan.NewLedger(storage.DeploymentPlan)

		req := &deploymentplan.CreateRequest{
			Name:        "test-deploymentplan",
			Namespace:   "test-namespace",
			ServiceName: "test-service",
			Applications: []deploymentplan.Application{
				{
					PayloadName: "test-payload",
					Resources: deploymentplan.ApplicationResources{
						Cores: 1,
					},
				},
			},
			MatchingComputeCapabilities: []deploymentplan.MatchingComputeCapability{
				{
					CapabilityType:  "test-capability",
					Comparator:      deploymentplan.ComparatorTypeIn,
					CapabilityNames: []string{"test-capability-name"},
				},
			},
		}
		resp, err := l.Create(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-deploymentplan", resp.Record.Name)
		require.NotEmpty(t, resp.Record.Metadata.ID)
		require.Equal(t, uint64(0), resp.Record.Metadata.Version)
		require.Equal(t, deploymentplan.DeploymentPlanStateActive, resp.Record.Status.State)
		require.Empty(t, resp.Record.Status.Message)
		require.Equal(t, "test-namespace", resp.Record.Namespace)
		require.Equal(t, "test-service", resp.Record.ServiceName)
		require.Len(t, resp.Record.Applications, 1)
		require.Equal(t, "test-payload", resp.Record.Applications[0].PayloadName)
		require.Equal(t, uint32(1), resp.Record.Applications[0].Resources.Cores)
		require.Len(t, resp.Record.MatchingComputeCapabilities, 1)
		require.Equal(t, "test-capability", resp.Record.MatchingComputeCapabilities[0].CapabilityType)
		require.Equal(t, deploymentplan.ComparatorTypeIn, resp.Record.MatchingComputeCapabilities[0].Comparator)
		require.Len(t, resp.Record.MatchingComputeCapabilities[0].CapabilityNames, 1)
		require.Equal(t, "test-capability-name", resp.Record.MatchingComputeCapabilities[0].CapabilityNames[0])
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

func TestLedgerGetByID(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	req := deplomentRequest()
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByID Success", func(t *testing.T) {
		resp, err := l.GetByID(context.Background(), createResp.Record.Metadata.ID)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-deploymentplan", resp.Record.Name)
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
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	req := deplomentRequest()
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
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	req := deplomentRequest()
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	lastUpdatedRecord := createResp.Record
	t.Run("UpdateStatus Success", func(t *testing.T) {
		updateReq := &deploymentplan.UpdateStatusRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: deploymentplan.DeploymentPlanStatus{
				State:   deploymentplan.DeploymentPlanStateInactive,
				Message: "DeploymentPlan is inactive now",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, deploymentplan.DeploymentPlanStateInactive, resp.Record.Status.State)
		lastUpdatedRecord = resp.Record
	})

	t.Run("UpdateStatus InvalidTransition Failure", func(t *testing.T) {
		updateReq := &deploymentplan.UpdateStatusRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: deploymentplan.DeploymentPlanStatus{
				State:   deploymentplan.DeploymentPlanStateUnknown, // Invalid transition
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
		updateReq := &deploymentplan.UpdateStatusRequest{
			Metadata: createResp.Record.Metadata, // This is the old metadata. Should cause a conflict.
			Status: deploymentplan.DeploymentPlanStatus{
				State:   deploymentplan.DeploymentPlanStateActive,
				Message: "DeploymentPlan is active now",
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
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	// Create two DeploymentPlans
	req := deplomentRequest()
	_, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	req.Name = "DeploymentPlan2"
	_, err = l.Create(context.Background(), req)
	require.NoError(t, err)

	// List the DeploymentPlans
	listReq := &deploymentplan.ListRequest{}
	resp, err := l.List(context.Background(), listReq)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Records, 2)
}

func TestLedgerDelete(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	// First, create the DeploymentPlan
	req := deplomentRequest()
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	// Now, delete the DeploymentPlan
	err = l.Delete(context.Background(), &deploymentplan.DeleteRequest{Metadata: createResp.Record.Metadata})
	require.NoError(t, err)

	// Try to get the DeploymentPlan again
	_, err = l.GetByID(context.Background(), createResp.Record.Metadata.ID)
	require.Error(t, err)
	require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
	require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
}

func TestDeployment(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := deploymentplan.NewLedger(storage.DeploymentPlan)

	// First, create the DeploymentPlan
	req := deplomentRequest()
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)
	updatedRecord := createResp.Record

	t.Run("AddDeployment Success", func(t *testing.T) {
		// Now, add a deployment
		addDeploymentReq := &deploymentplan.AddDeploymentRequest{
			Metadata:     createResp.Record.Metadata,
			DeploymentID: "test-deployment-1",
			PayloadCoordinates: []deploymentplan.PayloadCoordinates{
				{
					PayloadName: "test-payload",
					Coordinates: map[string]string{
						"key1": "value1",
					},
				},
			},
			InstanceCount: 3,
		}
		resp, err := l.AddDeployment(context.Background(), addDeploymentReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, len(resp.Record.Deployments), 1)
		require.Equal(t, "test-deployment-1", resp.Record.Deployments[0].ID)
		updatedRecord = resp.Record
	})

	t.Run("AddDeployment MissingPayloadCoordinates Failure", func(t *testing.T) {
		addDeploymentReq := &deploymentplan.AddDeploymentRequest{
			Metadata:     updatedRecord.Metadata,
			DeploymentID: "test-deployment-2",
			PayloadCoordinates: []deploymentplan.PayloadCoordinates{
				{
					PayloadName: "test-payload", // Missing coordinates
				},
			},
			InstanceCount: 3,
		}
		resp, err := l.AddDeployment(context.Background(), addDeploymentReq)

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})

	t.Run("AddDeployment MissingPayloadName Failure", func(t *testing.T) {
		addDeploymentReq := &deploymentplan.AddDeploymentRequest{
			Metadata:     updatedRecord.Metadata,
			DeploymentID: "test-deployment-3",
			PayloadCoordinates: []deploymentplan.PayloadCoordinates{
				{
					Coordinates: map[string]string{
						"key1": "value1",
					},
				},
			},
			InstanceCount: 3,
		}
		resp, err := l.AddDeployment(context.Background(), addDeploymentReq)

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})

	t.Run("Adding new deployment when there is an existing deployment in progress Failure", func(t *testing.T) {
		addDeploymentReq := &deploymentplan.AddDeploymentRequest{
			Metadata:     updatedRecord.Metadata,
			DeploymentID: "test-deployment-4",
			PayloadCoordinates: []deploymentplan.PayloadCoordinates{
				{
					PayloadName: "test-payload",
					Coordinates: map[string]string{
						"key1": "value1",
					},
				},
			},
			InstanceCount: 3,
		}
		resp, err := l.AddDeployment(context.Background(), addDeploymentReq)

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})

	t.Run("UpdateStatus Success", func(t *testing.T) {
		updateReq := &deploymentplan.UpdateDeploymentStatusRequest{
			Metadata:     updatedRecord.Metadata,
			DeploymentID: "test-deployment-1",
			Status: deploymentplan.DeploymentStatus{
				State:   deploymentplan.DeploymentStateCompleted,
				Message: "Deployment is completed",
			},
		}

		resp, err := l.UpdateDeploymentStatus(context.Background(), updateReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, deploymentplan.DeploymentStateCompleted, resp.Record.Deployments[0].Status.State)
		updatedRecord = resp.Record
	})

	t.Run("Adding new deployment when there is an existing deployment in completed state Success", func(t *testing.T) {
		addDeploymentReq := &deploymentplan.AddDeploymentRequest{
			Metadata:     updatedRecord.Metadata,
			DeploymentID: "test-deployment-5",
			PayloadCoordinates: []deploymentplan.PayloadCoordinates{
				{
					PayloadName: "test-payload",
					Coordinates: map[string]string{
						"key1": "value1",
					},
				},
			},
			InstanceCount: 3,
		}
		resp, err := l.AddDeployment(context.Background(), addDeploymentReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, len(resp.Record.Deployments), 2)
		require.Equal(t, "test-deployment-5", resp.Record.Deployments[1].ID)
		updatedRecord = resp.Record
	})
}
