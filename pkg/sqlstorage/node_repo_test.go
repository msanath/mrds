package sqlstorage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/msanath/mrds/ledger/core"
	ledgererrors "github.com/msanath/mrds/ledger/errors"
	"github.com/msanath/mrds/ledger/node"
	"github.com/msanath/mrds/pkg/sqlstorage/test"

	"github.com/stretchr/testify/require"
)

const nodeidPrefix = "node"

func TestNodeRecordLifecycle(t *testing.T) {
	storage := test.TestSQLStorage(t)

	testRecord := node.NodeRecord{
		Metadata: core.Metadata{
			ID:      fmt.Sprintf("%s1", nodeidPrefix),
			Version: 1,
		},
		Name: fmt.Sprintf("%s1", nodeidPrefix),
		Status: node.NodeStatus{
			State:   node.NodeStateUnallocated,
			Message: "",
		},
		UpdateDomain: "test-domain",
		TotalResources: node.Resources{
			Cores:  64,
			Memory: 512,
		},
		SystemReservedResources: node.Resources{
			Cores:  4,
			Memory: 32,
		},
		RemainingResources: node.Resources{
			Cores:  60,
			Memory: 480,
		},
		LocalVolumes: []node.LocalVolume{
			{
				MountPath:       "/var/lib/foo",
				StorageClass:    "SSD",
				StorageCapacity: 100,
			},
			{
				MountPath:       "/var/lib/bar",
				StorageClass:    "HDD",
				StorageCapacity: 200,
			},
		},
		CapabilityIDs: []string{"capability-1", "capability-2"},
	}
	repo := storage.Node
	ctx := context.Background()
	var err error

	t.Run("Insert Success", func(t *testing.T) {
		err = repo.Insert(ctx, testRecord)
		require.NoError(t, err)
	})

	t.Run("Insert Duplicate Failure", func(t *testing.T) {
		newRecord := testRecord
		// Change the ID to create a duplicate record.
		newRecord.Metadata.ID = fmt.Sprintf("%s2", nodeidPrefix)
		err = repo.Insert(ctx, newRecord)
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordInsertConflict, err.(ledgererrors.LedgerError).Code)
	})

	t.Run("Get By Metadata Success", func(t *testing.T) {
		receivedRecord, err := repo.GetByID(ctx, testRecord.Metadata.ID)
		require.NoError(t, err)
		require.Equal(t, testRecord.Metadata, receivedRecord.Metadata)
		require.Equal(t, testRecord.Name, receivedRecord.Name)
		require.Equal(t, testRecord.Status, receivedRecord.Status)
		require.Equal(t, testRecord.UpdateDomain, receivedRecord.UpdateDomain)
		require.Equal(t, testRecord.TotalResources, receivedRecord.TotalResources)
		require.Equal(t, testRecord.SystemReservedResources, receivedRecord.SystemReservedResources)
		require.Equal(t, testRecord.RemainingResources, receivedRecord.RemainingResources)
		require.ElementsMatch(t, testRecord.LocalVolumes, receivedRecord.LocalVolumes)
		require.ElementsMatch(t, testRecord.CapabilityIDs, receivedRecord.CapabilityIDs)
	})

	t.Run("Get By Name Success", func(t *testing.T) {
		receivedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Equal(t, testRecord.Metadata, receivedRecord.Metadata)
		require.Equal(t, testRecord.Name, receivedRecord.Name)
		require.Equal(t, testRecord.Status, receivedRecord.Status)
		require.Equal(t, testRecord.UpdateDomain, receivedRecord.UpdateDomain)
		require.Equal(t, testRecord.TotalResources, receivedRecord.TotalResources)
		require.Equal(t, testRecord.SystemReservedResources, receivedRecord.SystemReservedResources)
		require.Equal(t, testRecord.RemainingResources, receivedRecord.RemainingResources)
		require.ElementsMatch(t, testRecord.LocalVolumes, receivedRecord.LocalVolumes)
		require.ElementsMatch(t, testRecord.CapabilityIDs, receivedRecord.CapabilityIDs)

	})

	t.Run("Get By Name Failure", func(t *testing.T) {
		_, err := repo.GetByName(ctx, "unknown")
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code, err.Error())
	})

	t.Run("Update State Success", func(t *testing.T) {
		status := node.NodeStatus{
			State:   "error",
			Message: "Needs attention",
		}

		err = repo.UpdateStatus(ctx, testRecord.Metadata, status, "")
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Equal(t, status, updatedRecord.Status)
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Add disruption", func(t *testing.T) {
		disruption := node.Disruption{
			ID:          "disruption-1",
			StartTime:   time.Now().Truncate(time.Second),
			ShouldEvict: true,
			Status: node.DisruptionStatus{
				State:   node.DisruptionStateScheduled,
				Message: "Scheduled",
			},
		}

		err = repo.InsertDisruption(ctx, testRecord.Metadata, disruption)
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.Disruptions, 1)
		require.Equal(t, disruption, updatedRecord.Disruptions[0])
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Update disruption", func(t *testing.T) {
		err := repo.UpdateDisruptionStatus(ctx, testRecord.Metadata, "disruption-1", node.DisruptionStatus{
			State:   node.DisruptionStateApproved,
			Message: "Approved",
		})
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.Disruptions, 1)
		require.Equal(t, node.DisruptionStateApproved, updatedRecord.Disruptions[0].Status.State)
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Add capability", func(t *testing.T) {
		capability := "capability-3"
		err = repo.InsertCapability(ctx, testRecord.Metadata, capability)
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Contains(t, updatedRecord.CapabilityIDs, capability)
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Delete capability", func(t *testing.T) {
		err := repo.DeleteCapability(ctx, testRecord.Metadata, "capability-3")
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.NotContains(t, updatedRecord.CapabilityIDs, "capability-3")
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Delete disruption", func(t *testing.T) {
		err := repo.DeleteDisruption(ctx, testRecord.Metadata, "disruption-1")
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.Disruptions, 0)
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
			newRecord.Metadata.ID = fmt.Sprintf("%s-%d", nodeidPrefix, i+1)
			newRecord.Metadata.Version = 0
			newRecord.Name = fmt.Sprintf("%s-%d", nodeidPrefix, i+1)
			newRecord.Status.State = node.NodeStateUnallocated
			newRecord.Status.Message = fmt.Sprintf("%s-%d is active", nodeidPrefix, i+1)

			if (i+1)%2 == 0 {
				newRecord.Status.State = node.NodeStateAllocating
				newRecord.Status.Message = fmt.Sprintf("%s-%d is inactive", nodeidPrefix, i+1)
			}
			if (i+1)%3 == 0 {
				newRecord.RemainingResources.Cores = 20
				newRecord.RemainingResources.Memory = 160
			}
			if (i+1)%4 == 0 {
				newRecord.Status.State = node.NodeStateAllocated
				newRecord.ClusterID = "cluster-1"
			}

			err = repo.Insert(ctx, newRecord)
			require.NoError(t, err)
		}
	})

	t.Run("List", func(t *testing.T) {
		records, err := repo.List(ctx, node.NodeListFilters{})
		require.NoError(t, err)
		require.Len(t, records, 10)

		receivedIDs := []string{}
		for _, record := range records {
			receivedIDs = append(receivedIDs, record.Metadata.ID)

		}
		expectedIDs := []string{}
		for i := range 10 {
			expectedIDs = append(expectedIDs, fmt.Sprintf("%s-%d", nodeidPrefix, i+1))

		}
		require.ElementsMatch(t, expectedIDs, receivedIDs)
		allRecords := records

		t.Run("List Success With Filter", func(t *testing.T) {
			records, err := repo.List(ctx, node.NodeListFilters{
				StateIn: []node.NodeState{node.NodeStateUnallocated},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, node.NodeStateUnallocated, record.Status.State)
			}
		})

		t.Run("List with Names Filter", func(t *testing.T) {
			records, err := repo.List(ctx, node.NodeListFilters{
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
			records, err := repo.List(ctx, node.NodeListFilters{
				Limit: 3,
			})
			require.NoError(t, err)
			require.Len(t, records, 3)
		})

		t.Run("List with IncludeDeleted", func(t *testing.T) {
			err = repo.Delete(ctx, allRecords[0].Metadata)
			require.NoError(t, err)

			records, err := repo.List(ctx, node.NodeListFilters{
				IncludeDeleted: true,
			})
			require.NoError(t, err)
			require.Len(t, records, 11)
		})

		t.Run("List with StateNotIn", func(t *testing.T) {
			records, err := repo.List(ctx, node.NodeListFilters{
				StateNotIn: []node.NodeState{node.NodeStateUnallocated},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.NotEqual(t, node.NodeStateUnallocated, record.Status.State)
			}
		})

		t.Run("Update State and check version", func(t *testing.T) {
			status := node.NodeStatus{
				State:   node.NodeStateEvicted,
				Message: "Needs attention",
			}

			err = repo.UpdateStatus(ctx, allRecords[1].Metadata, status, "")
			require.NoError(t, err)

			ve := uint64(1)
			records, err := repo.List(ctx, node.NodeListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 1)

			ve += 1
			records, err = repo.List(ctx, node.NodeListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 0)
		})

		t.Run("List Remaining Resources", func(t *testing.T) {
			remainingCores := uint32(21)
			records, err := repo.List(ctx, node.NodeListFilters{
				RemainingCoresGte: &remainingCores,
			})
			require.NoError(t, err)
			require.Len(t, records, 6)
			for _, record := range records {
				require.GreaterOrEqual(t, record.RemainingResources.Cores, uint32(21))
			}
		})

		t.Run("List by ClusterID", func(t *testing.T) {
			records, err := repo.List(ctx, node.NodeListFilters{
				ClusterIDIn: []string{"cluster-1"},
			})
			require.NoError(t, err)
			require.Len(t, records, 2)
			for _, record := range records {
				require.Equal(t, "cluster-1", record.ClusterID)
			}
		})
	})
}
