package metainstance_test

import (
	"context"
	"testing"

	"github.com/msanath/mrds/ledger/deploymentplan"
	ledgererrors "github.com/msanath/mrds/ledger/errors"
	"github.com/msanath/mrds/ledger/metainstance"
	"github.com/msanath/mrds/ledger/node"
	"github.com/msanath/mrds/pkg/sqlstorage/test"
	"github.com/stretchr/testify/require"
)

func CreateRequest() *metainstance.CreateRequest {

	return &metainstance.CreateRequest{
		Name:             "test-metainstance",
		DeploymentPlanID: "test-deployment-plan",
		DeploymentID:     "test-deployment",
	}
}

func TestOperationsLedger(t *testing.T) {
	// Pre-requiste - create a deployment Plan and add a deployment
	storage := test.TestSQLStorage(t)
	dl := deploymentplan.NewLedger(storage.DeploymentPlan)
	resp, err := dl.Create(context.Background(), &deploymentplan.CreateRequest{
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
	})
	require.NoError(t, err)
	dPlanUpdateResp, err := dl.AddDeployment(context.Background(), &deploymentplan.AddDeploymentRequest{
		Metadata:     resp.Record.Metadata,
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
	})
	require.NoError(t, err)

	// Create a node
	nl := node.NewLedger(storage.Node)
	nodeCreateResp, err := nl.Create(context.Background(), &node.CreateRequest{
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
		CapabilityIDs: []string{"test-capability", "test-capability-1"},
	})

	var lastUpdatedRecord metainstance.MetaInstanceRecord
	t.Run("Create Success", func(t *testing.T) {
		l := metainstance.NewLedger(storage.MetaInstance)

		req := &metainstance.CreateRequest{
			Name:             "test-metainstance",
			DeploymentPlanID: resp.Record.Metadata.ID,
			DeploymentID:     "test-deployment-1",
		}
		resp, err := l.Create(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-metainstance", resp.Record.Name)
		require.NotEmpty(t, resp.Record.Metadata.ID)
		require.Equal(t, uint64(0), resp.Record.Metadata.Version)
		require.Equal(t, metainstance.MetaInstanceStateActive, resp.Record.Status.State)
		lastUpdatedRecord = resp.Record
	})

	l := metainstance.NewLedger(storage.MetaInstance)
	t.Run("Create EmptyName Failure", func(t *testing.T) {

		req := &metainstance.CreateRequest{
			Name: "", // Empty name
		}
		resp, err := l.Create(context.Background(), req)

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})

	t.Run("GetByID Success", func(t *testing.T) {
		resp, err := l.GetByID(context.Background(), lastUpdatedRecord.Metadata.ID)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-metainstance", resp.Record.Name)
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

	t.Run("GetByName Success", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "test-metainstance")

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-metainstance", resp.Record.Name)
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

	t.Run("UpdateStatus Success", func(t *testing.T) {
		updateReq := &metainstance.UpdateStatusRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: metainstance.MetaInstanceStatus{
				State:   metainstance.MetaInstanceStateMarkedForDeletion,
				Message: "MetaInstance is marked for deletion",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, metainstance.MetaInstanceStateMarkedForDeletion, resp.Record.Status.State)
		lastUpdatedRecord = resp.Record
	})

	t.Run("UpdateStatus InvalidTransition Failure", func(t *testing.T) {
		updateReq := &metainstance.UpdateStatusRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: metainstance.MetaInstanceStatus{
				State:   metainstance.MetaInstanceStateActive, // Invalid transition
				Message: "Invalid state transition",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})

	t.Run("Update Deployment ID Success", func(t *testing.T) {
		dPlanUpdateResp, err = dl.AddDeployment(context.Background(), &deploymentplan.AddDeploymentRequest{
			Metadata:     dPlanUpdateResp.Record.Metadata,
			DeploymentID: "test-deployment-2",
			PayloadCoordinates: []deploymentplan.PayloadCoordinates{
				{
					PayloadName: "test-payload",
					Coordinates: map[string]string{
						"key1": "value1",
					},
				},
			},
			InstanceCount: 3,
		})
		require.NoError(t, err)

		resp, err := l.UpdateDeploymentID(context.Background(), &metainstance.UpdateDeploymentIDRequest{
			Metadata:     lastUpdatedRecord.Metadata,
			DeploymentID: "test-deployment-2",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-deployment-2", resp.Record.DeploymentID)
		lastUpdatedRecord = resp.Record
	})

	t.Run("Add RuntimeInstance Success", func(t *testing.T) {
		resp, err := l.AddRuntimeInstance(context.Background(), &metainstance.AddRuntimeInstanceRequest{
			Metadata: lastUpdatedRecord.Metadata,
			RuntimeInstance: metainstance.RuntimeInstance{
				ID:       "test-runtime-instance",
				NodeID:   nodeCreateResp.Record.Metadata.ID,
				IsActive: true,
				Status: metainstance.RuntimeInstanceStatus{
					State:   metainstance.RuntimeStateRunning,
					Message: "Runtime instance is running",
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Record.RuntimeInstances, 1)
		require.Equal(t, "test-runtime-instance", resp.Record.RuntimeInstances[0].ID)
		lastUpdatedRecord = resp.Record
	})

	t.Run("Update RuntimeInstance Status Success", func(t *testing.T) {
		resp, err := l.UpdateRuntimeStatus(context.Background(), &metainstance.UpdateRuntimeStatusRequest{
			Metadata:          lastUpdatedRecord.Metadata,
			RuntimeInstanceID: "test-runtime-instance",
			Status: metainstance.RuntimeInstanceStatus{
				State:   metainstance.RuntimeStateTerminated,
				Message: "Runtime instance is terminated",
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Check the status of the runtime instance
		require.Len(t, resp.Record.RuntimeInstances, 1)
		require.Equal(t, "test-runtime-instance", resp.Record.RuntimeInstances[0].ID)
		require.Equal(t, metainstance.RuntimeStateTerminated, resp.Record.RuntimeInstances[0].Status.State)
		lastUpdatedRecord = resp.Record
	})

	t.Run("Remove RuntimeInstance Success", func(t *testing.T) {
		resp, err := l.RemoveRuntimeInstance(context.Background(), &metainstance.RemoveRuntimeInstanceRequest{
			Metadata:          lastUpdatedRecord.Metadata,
			RuntimeInstanceID: "test-runtime-instance",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Check the status of the runtime instance
		require.Len(t, resp.Record.RuntimeInstances, 0)
		lastUpdatedRecord = resp.Record
	})

	t.Run("List Success", func(t *testing.T) {
		l := metainstance.NewLedger(storage.MetaInstance)

		req := &metainstance.CreateRequest{
			Name:             "test-metainstance-2",
			DeploymentPlanID: resp.Record.Metadata.ID,
			DeploymentID:     "test-deployment-2",
		}
		resp, err := l.Create(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// List the MetaInstances
		listReq := &metainstance.ListRequest{}
		listResp, err := l.List(context.Background(), listReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, listResp.Records, 2)
	})

	t.Run("Add Operation Success", func(t *testing.T) {
		resp, err := l.AddOperation(context.Background(), &metainstance.AddOperationRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Operation: metainstance.Operation{
				ID:       "test-operation",
				Type:     "test-operation-type",
				IntentID: "test-intent",
				Status: metainstance.OperationStatus{
					State:   metainstance.OperationStatePendingApproval,
					Message: "Operation is running",
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Record.Operations, 1)
		require.Equal(t, "test-operation", resp.Record.Operations[0].ID)
		lastUpdatedRecord = resp.Record
	})

	t.Run("Update Operation Status Success", func(t *testing.T) {
		resp, err := l.UpdateOperationStatus(context.Background(), &metainstance.UpdateOperationStatusRequest{
			Metadata:    lastUpdatedRecord.Metadata,
			OperationID: "test-operation",
			Status: metainstance.OperationStatus{
				State:   metainstance.OperationStateSucceeded,
				Message: "Operation is successful",
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Record.Operations, 1)
		require.Equal(t, "test-operation", resp.Record.Operations[0].ID)
		require.Equal(t, metainstance.OperationStateSucceeded, resp.Record.Operations[0].Status.State)
		lastUpdatedRecord = resp.Record
	})

	t.Run("Remove Operation Success", func(t *testing.T) {
		resp, err := l.RemoveOperation(context.Background(), &metainstance.RemoveOperationRequest{
			Metadata:    lastUpdatedRecord.Metadata,
			OperationID: "test-operation",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Record.Operations, 0)
		lastUpdatedRecord = resp.Record
	})

	t.Run("Delete Success", func(t *testing.T) {
		// Now, delete the MetaInstance
		err = l.Delete(context.Background(), &metainstance.DeleteRequest{Metadata: lastUpdatedRecord.Metadata})
		require.NoError(t, err)

		// Try to get the MetaInstance again
		_, err = l.GetByID(context.Background(), lastUpdatedRecord.Metadata.ID)
		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
	})

}
