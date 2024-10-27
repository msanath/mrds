package sqlstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/msanath/mrds/internal/ledger/core"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"
	"github.com/msanath/mrds/internal/ledger/node"
	"github.com/msanath/mrds/internal/sqlstorage"

	"github.com/msanath/gondolf/pkg/simplesql/test"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const nodeidPrefix = "node"

func TestNodeRecordLifecycle(t *testing.T) {
	db, err := test.NewTestSQLiteDB()
	require.NoError(t, err)

	storage, err := sqlstorage.NewSQLStorage(db, true)
	require.NoError(t, err)

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
	}
	repo := storage.Node
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
		receivedRecord, err := repo.GetByMetadata(ctx, testRecord.Metadata)
		require.NoError(t, err)
		require.Equal(t, testRecord, receivedRecord)
		require.Equal(t, testRecord.RemainingResources, receivedRecord.RemainingResources)
	})

	t.Run("Get By Metadata failure", func(t *testing.T) {
		metadata := core.Metadata{
			ID:      testRecord.Metadata.ID,
			Version: 2, // Different version
		}
		_, err := repo.GetByMetadata(ctx, metadata)
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code, err.Error())
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
		status := node.NodeStatus{
			State:   "error",
			Message: "Needs attention",
		}

		err = repo.UpdateState(ctx, testRecord.Metadata, status)
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Equal(t, status, updatedRecord.Status)
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

			err = repo.UpdateState(ctx, allRecords[1].Metadata, status)
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
