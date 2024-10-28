package sqlstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/msanath/mrds/internal/ledger/core"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"
	"github.com/msanath/mrds/internal/ledger/metainstance"
	"github.com/msanath/mrds/internal/sqlstorage/test"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const metaInstanceidPrefix = "metainstance"

func TestMetaInstanceRecordLifecycle(t *testing.T) {
	storage := test.TestSQLStorage(t)

	testRecord := metainstance.MetaInstanceRecord{
		Metadata: core.Metadata{
			ID:      fmt.Sprintf("%s1", metaInstanceidPrefix),
			Version: 1,
		},
		Name: fmt.Sprintf("%s1", metaInstanceidPrefix),
		Status: metainstance.MetaInstanceStatus{
			State:   metainstance.MetaInstanceStateActive,
			Message: "",
		},
	}
	repo := storage.MetaInstance

	ctx := context.Background()
	var err error

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
		status := metainstance.MetaInstanceStatus{
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
			newRecord.Metadata.ID = fmt.Sprintf("%s-%d", metaInstanceidPrefix, i+1)
			newRecord.Metadata.Version = 0
			newRecord.Name = fmt.Sprintf("%s-%d", metaInstanceidPrefix, i+1)
			newRecord.Status.State = metainstance.MetaInstanceStateActive
			newRecord.Status.Message = fmt.Sprintf("%s-%d is active", metaInstanceidPrefix, i+1)

			if (i+1)%2 == 0 {
				newRecord.Status.State = metainstance.MetaInstanceStateInActive
				newRecord.Status.Message = fmt.Sprintf("%s-%d is inactive", metaInstanceidPrefix, i+1)
			}

			err = repo.Insert(ctx, newRecord)
			require.NoError(t, err)
		}
	})

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
				StateIn: []metainstance.MetaInstanceState{metainstance.MetaInstanceStateActive},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, metainstance.MetaInstanceStateActive, record.Status.State)
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
			err = repo.Delete(ctx, allRecords[0].Metadata)
			require.NoError(t, err)

			records, err := repo.List(ctx, metainstance.MetaInstanceListFilters{
				IncludeDeleted: true,
			})
			require.NoError(t, err)
			require.Len(t, records, 11)
		})

		t.Run("List with StateNotIn", func(t *testing.T) {
			records, err := repo.List(ctx, metainstance.MetaInstanceListFilters{
				StateNotIn: []metainstance.MetaInstanceState{metainstance.MetaInstanceStateActive},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, metainstance.MetaInstanceStateInActive, record.Status.State)
			}
		})

		t.Run("Update State and check version", func(t *testing.T) {
			status := metainstance.MetaInstanceStatus{
				State:   metainstance.MetaInstanceStatePending,
				Message: "Needs attention",
			}

			err = repo.UpdateState(ctx, allRecords[1].Metadata, status)
			require.NoError(t, err)
			ve := uint64(1)
			records, err := repo.List(ctx, metainstance.MetaInstanceListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 1)

			ve += 1
			records, err = repo.List(ctx, metainstance.MetaInstanceListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 0)
		})
	})
}
