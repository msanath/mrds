package sqlstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/msanath/mrds/ledger/core"
	"github.com/msanath/mrds/ledger/deploymentplan"
	ledgererrors "github.com/msanath/mrds/ledger/errors"
	"github.com/msanath/mrds/ledger/metainstance"
	"github.com/msanath/mrds/ledger/node"
	"github.com/msanath/mrds/pkg/sqlstorage/test"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const metaInstanceidPrefix = "metainstance"

func TestMetaInstanceRecordLifecycle(t *testing.T) {
	storage := test.TestSQLStorage(t)

	// Create a deployment plan
	err := storage.DeploymentPlan.Insert(context.Background(), deploymentplan.DeploymentPlanRecord{
		Metadata: core.Metadata{
			ID:      "dp1",
			Version: 1,
		},
		Name: "dp1",
		Applications: []deploymentplan.Application{
			{
				PayloadName: "app1",
				Resources: deploymentplan.ApplicationResources{
					Cores:  12,
					Memory: 64,
				},
			},
		},
	})
	require.NoError(t, err)

	// Add a deployment
	err = storage.DeploymentPlan.InsertDeployment(context.Background(), core.Metadata{
		ID:      "dp1",
		Version: 1,
	}, deploymentplan.Deployment{
		ID: "d1",
	})
	require.NoError(t, err)

	testRecord := metainstance.MetaInstanceRecord{
		Metadata: core.Metadata{
			ID:      fmt.Sprintf("%s-0", metaInstanceidPrefix),
			Version: 1,
		},
		Name: fmt.Sprintf("%s-0", metaInstanceidPrefix),
		Status: metainstance.MetaInstanceStatus{
			State:   metainstance.MetaInstanceStateRunning,
			Message: "",
		},
		DeploymentPlanID: "dp1",
		DeploymentID:     "d1",
	}
	repo := storage.MetaInstance

	ctx := context.Background()

	t.Run("Insert Success", func(t *testing.T) {
		err = repo.Insert(ctx, testRecord)
		require.NoError(t, err)
	})

	t.Run("Insert Duplicate Failure", func(t *testing.T) {
		err = repo.Insert(ctx, testRecord)
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordInsertConflict, err.(ledgererrors.LedgerError).Code)
	})

	t.Run("Get By Metadata Success", func(t *testing.T) {
		receivedRecord, err := repo.GetByID(ctx, testRecord.Metadata.ID)
		require.NoError(t, err)
		require.Equal(t, testRecord.Name, receivedRecord.Name)
		require.Equal(t, testRecord.Status, receivedRecord.Status)
		require.Equal(t, testRecord.DeploymentPlanID, receivedRecord.DeploymentPlanID)
		require.Equal(t, testRecord.DeploymentID, receivedRecord.DeploymentID)
		require.Equal(t, testRecord.Metadata, receivedRecord.Metadata)
	})

	t.Run("Get By Name Success", func(t *testing.T) {
		receivedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(testRecord, receivedRecord))
	})

	t.Run("Get By Name Failure", func(t *testing.T) {
		_, err := repo.GetByName(ctx, "unknown")
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code, err.Error())
	})

	t.Run("Update State Success", func(t *testing.T) {
		status := metainstance.MetaInstanceStatus{
			State:   "error",
			Message: "Needs attention",
		}

		err = repo.UpdateStatus(ctx, testRecord.Metadata, status)
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Equal(t, status, updatedRecord.Status)
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Add Operation Success", func(t *testing.T) {
		operation := metainstance.Operation{
			ID:       "op1",
			Type:     "create",
			IntentID: "intent1",
			Status: metainstance.OperationStatus{
				State:   metainstance.OperationStatePendingApproval,
				Message: "Needs attention",
			},
		}

		err = repo.InsertOperation(ctx, testRecord.Metadata, operation)
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.Operations, 1)
		require.Equal(t, operation, updatedRecord.Operations[0])
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Add Runtime Instance With invalid node ID Failure", func(t *testing.T) {
		runtimeInstance := metainstance.RuntimeInstance{
			ID:       "ri1",
			NodeID:   "unknown",
			IsActive: true,
			Status: metainstance.RuntimeInstanceStatus{
				State:   metainstance.RuntimeStateRunning,
				Message: "In progress",
			},
		}

		err = repo.InsertRuntimeInstance(ctx, testRecord.Metadata, runtimeInstance)
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
	})

	t.Run("Add Runtime Instance Success", func(t *testing.T) {
		// Create a node
		storage.Node.Insert(ctx, node.NodeRecord{
			Metadata: core.Metadata{
				ID:      "node1",
				Version: 1,
			},
			Name: "node1",
			TotalResources: node.Resources{
				Cores:  50,
				Memory: 256,
			},
			SystemReservedResources: node.Resources{
				Cores:  2,
				Memory: 8,
			},
			RemainingResources: node.Resources{
				Cores:  48,
				Memory: 248,
			},
			Status: node.NodeStatus{
				State:   node.NodeStateAllocated,
				Message: "Node is active",
			},
		})

		runtimeInstance := metainstance.RuntimeInstance{
			ID:       "ri1",
			NodeID:   "node1",
			IsActive: true,
			Status: metainstance.RuntimeInstanceStatus{
				State:   metainstance.RuntimeStateRunning,
				Message: "In progress",
			},
		}

		err = repo.InsertRuntimeInstance(ctx, testRecord.Metadata, runtimeInstance)
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.RuntimeInstances, 1)
		require.Equal(t, runtimeInstance, updatedRecord.RuntimeInstances[0])
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord

		node, err := storage.Node.GetByID(ctx, "node1")
		require.NoError(t, err)
		// Node has 48 cores and 248 memory remaining. The app uses 12 cores and 64 memory.
		// After the runtime instance is added, the remaining resources should be 48-12=36 cores and 248-64=184 memory.
		require.Equal(t, node.RemainingResources.Cores, uint32(36))
		require.Equal(t, node.RemainingResources.Memory, uint32(184))
	})

	t.Run("Add Runtime Instance when no remaining failure", func(t *testing.T) {
		// Create a node
		storage.Node.Insert(ctx, node.NodeRecord{
			Metadata: core.Metadata{
				ID:      "node2",
				Version: 1,
			},
			Name: "node2",
			TotalResources: node.Resources{
				Cores:  50,
				Memory: 256,
			},
			SystemReservedResources: node.Resources{
				Cores:  2,
				Memory: 8,
			},
			RemainingResources: node.Resources{
				Cores:  1,
				Memory: 1,
			},
			Status: node.NodeStatus{
				State:   node.NodeStateAllocated,
				Message: "Node is active",
			},
		})

		runtimeInstance := metainstance.RuntimeInstance{
			ID:       "ri2",
			NodeID:   "node2",
			IsActive: true,
			Status: metainstance.RuntimeInstanceStatus{
				State:   metainstance.RuntimeStateRunning,
				Message: "In progress",
			},
		}

		err = repo.InsertRuntimeInstance(ctx, testRecord.Metadata, runtimeInstance)
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordInsertConflict, err.(ledgererrors.LedgerError).Code)
		require.ErrorContains(t, err, "does not have enough")
	})

	t.Run("Update Runtime Instance Status Success", func(t *testing.T) {
		err = repo.UpdateRuntimeInstanceStatus(ctx, testRecord.Metadata, "ri1", metainstance.RuntimeInstanceStatus{
			State:   metainstance.RuntimeStateTerminated,
			Message: "Is Terminated",
		})
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.RuntimeInstances, 1)
		require.Equal(t, metainstance.RuntimeInstanceStatus{
			State:   metainstance.RuntimeStateTerminated,
			Message: "Is Terminated",
		}, updatedRecord.RuntimeInstances[0].Status)
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Delete Runtime Instance Success", func(t *testing.T) {
		err = repo.DeleteRuntimeInstance(ctx, testRecord.Metadata, "ri1")
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.RuntimeInstances, 0)
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord

		node, err := storage.Node.GetByID(ctx, "node1")
		require.NoError(t, err)
		// Node has 48 cores and 248 memory remaining. The app uses 12 cores and 64 memory.
		// After the runtime instance is deleted, the remaining resources should be back to 48 cores and 248 memory.
		require.Equal(t, node.RemainingResources.Cores, uint32(48))
		require.Equal(t, node.RemainingResources.Memory, uint32(248))
	})

	t.Run("Operation Status Update Success", func(t *testing.T) {
		err = repo.UpdateOperationStatus(ctx, testRecord.Metadata, "op1", metainstance.OperationStatus{
			State:   metainstance.OperationStateSucceeded,
			Message: "In progress",
		})
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.Operations, 1)
		require.Equal(t, metainstance.OperationStatus{
			State:   metainstance.OperationStateSucceeded,
			Message: "In progress",
		}, updatedRecord.Operations[0].Status)
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Operation Delete Success", func(t *testing.T) {
		err = repo.DeleteOperation(ctx, testRecord.Metadata, "op1")
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.Operations, 0)
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Delete Success", func(t *testing.T) {
		err = repo.Delete(ctx, testRecord.Metadata)
		require.NoError(t, err)

		_, err = repo.GetByName(ctx, testRecord.Name)
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
	})

	t.Run("Create More Resources", func(t *testing.T) {
		// Create 10 records.
		for i := range 10 {
			newRecord := testRecord
			newRecord.Metadata.ID = fmt.Sprintf("%s-%d", metaInstanceidPrefix, i+1)
			newRecord.Metadata.Version = 0
			newRecord.Name = fmt.Sprintf("%s-%d", metaInstanceidPrefix, i+1)
			newRecord.Status.State = metainstance.MetaInstanceStateRunning
			newRecord.Status.Message = fmt.Sprintf("%s-%d is active", metaInstanceidPrefix, i+1)

			if (i+1)%2 == 0 {
				newRecord.Status.State = metainstance.MetaInstanceStateTerminated
				newRecord.Status.Message = fmt.Sprintf("%s-%d is inactive", metaInstanceidPrefix, i+1)
			}

			err = repo.Insert(ctx, newRecord)
			require.NoError(t, err)
		}
	})

	// Add operations on every 3rd record.
	for i := 0; i < 10; i += 3 {
		operation := metainstance.Operation{
			ID:       fmt.Sprintf("op-%d", i+1),
			Type:     "create",
			IntentID: fmt.Sprintf("intent-%d", i+1),
			Status: metainstance.OperationStatus{
				State:   metainstance.OperationStatePendingApproval,
				Message: "Needs attention",
			},
		}

		err = repo.InsertOperation(ctx, core.Metadata{
			ID:      fmt.Sprintf("%s-%d", metaInstanceidPrefix, i+1),
			Version: 0,
		}, operation)

		require.NoError(t, err)
	}

	t.Run("List", func(t *testing.T) {
		records, err := repo.List(ctx, metainstance.MetaInstanceListFilters{})
		require.NoError(t, err)
		require.Len(t, records, 10)

		receivedIDs := []string{}
		for _, record := range records {
			receivedIDs = append(receivedIDs, record.Metadata.ID)
		}
		expectedIDs := []string{}
		for i := range 10 {
			expectedIDs = append(expectedIDs, fmt.Sprintf("%s-%d", metaInstanceidPrefix, i+1))

		}
		require.ElementsMatch(t, expectedIDs, receivedIDs)
		allRecords := records

		t.Run("List Success With Filter", func(t *testing.T) {
			records, err := repo.List(ctx, metainstance.MetaInstanceListFilters{
				StateIn: []metainstance.MetaInstanceState{metainstance.MetaInstanceStateRunning},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, metainstance.MetaInstanceStateRunning, record.Status.State)
			}
		})

		t.Run("List with Names Filter", func(t *testing.T) {
			records, err := repo.List(ctx, metainstance.MetaInstanceListFilters{
				NameIn: []string{allRecords[0].Name, allRecords[1].Name, allRecords[2].Name},
			})
			require.NoError(t, err)
			require.Len(t, records, 3)

			// Check if the returned records are the same as the first 3 computeCapabilitys.
			for i, record := range records {
				require.Equal(t, allRecords[i], record)
			}
		})

		t.Run("List with Limit", func(t *testing.T) {
			records, err := repo.List(ctx, metainstance.MetaInstanceListFilters{
				Limit: 3,
			})
			require.NoError(t, err)
			require.Len(t, records, 3)
		})

		t.Run("List with IncludeDeleted", func(t *testing.T) {
			// Get record with ID 1 and delete it.
			rec, err := repo.GetByName(ctx, fmt.Sprintf("%s-1", metaInstanceidPrefix))
			require.NoError(t, err)
			err = repo.Delete(ctx, rec.Metadata)
			require.NoError(t, err)

			records, err := repo.List(ctx, metainstance.MetaInstanceListFilters{
				IncludeDeleted: true,
			})
			require.NoError(t, err)
			require.Len(t, records, 11)
		})

		t.Run("List with StateNotIn", func(t *testing.T) {
			records, err := repo.List(ctx, metainstance.MetaInstanceListFilters{
				StateNotIn: []metainstance.MetaInstanceState{metainstance.MetaInstanceStateRunning},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, metainstance.MetaInstanceStateTerminated, record.Status.State)
			}
		})

		t.Run("Update State and check version", func(t *testing.T) {
			status := metainstance.MetaInstanceStatus{
				State:   metainstance.MetaInstanceStatePendingAllocation,
				Message: "Needs attention",
			}
			rec, err := repo.GetByName(ctx, fmt.Sprintf("%s-2", metaInstanceidPrefix))
			require.NoError(t, err)
			err = repo.UpdateStatus(ctx, rec.Metadata, status)
			require.NoError(t, err)
			ve := uint64(1)
			records, err := repo.List(ctx, metainstance.MetaInstanceListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 4)

			ve += 2
			records, err = repo.List(ctx, metainstance.MetaInstanceListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 0)
		})
	})
}
