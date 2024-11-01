package sqlstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/msanath/mrds/internal/ledger/computecapability"
	"github.com/msanath/mrds/internal/ledger/core"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"
	"github.com/msanath/mrds/internal/sqlstorage/test"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const computeCapabilityidPrefix = "computecapability"

func TestComputeCapabilityRecordLifecycle(t *testing.T) {
	storage := test.TestSQLStorage(t)

	testRecord := computecapability.ComputeCapabilityRecord{
		Metadata: core.Metadata{
			ID:      fmt.Sprintf("%s1", computeCapabilityidPrefix),
			Version: 1,
		},
		Name: fmt.Sprintf("%s1", computeCapabilityidPrefix),
		Status: computecapability.ComputeCapabilityStatus{
			State:   computecapability.ComputeCapabilityStateActive,
			Message: "",
		},
		Type:  "CPU",
		Score: 10,
	}
	repo := storage.ComputeCapability

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
		receivedRecord, err := repo.GetByID(ctx, testRecord.Metadata.ID)
		require.NoError(t, err)
		require.Equal(t, testRecord, receivedRecord)
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
		status := computecapability.ComputeCapabilityStatus{
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
			newRecord.Metadata.ID = fmt.Sprintf("%s-%d", computeCapabilityidPrefix, i+1)
			newRecord.Metadata.Version = 0
			newRecord.Name = fmt.Sprintf("%s-%d", computeCapabilityidPrefix, i+1)
			newRecord.Status.State = computecapability.ComputeCapabilityStateActive
			newRecord.Status.Message = fmt.Sprintf("%s-%d is active", computeCapabilityidPrefix, i+1)

			if (i+1)%2 == 0 {
				newRecord.Status.State = computecapability.ComputeCapabilityStateInActive
				newRecord.Status.Message = fmt.Sprintf("%s-%d is inactive", computeCapabilityidPrefix, i+1)
			}

			// Change capacity type and score for every 3rd record.
			if (i+1)%3 == 0 {
				newRecord.Type = "GPU"
				newRecord.Score = 20
			}

			err = repo.Insert(ctx, newRecord)
			require.NoError(t, err)
		}
	})

	t.Run("List", func(t *testing.T) {
		records, err := repo.List(ctx, computecapability.ComputeCapabilityListFilters{})
		require.NoError(t, err)
		require.Len(t, records, 10)

		receivedIDs := []string{}
		for _, record := range records {
			receivedIDs = append(receivedIDs, record.Metadata.ID)

		}
		expectedIDs := []string{}
		for i := range 10 {
			expectedIDs = append(expectedIDs, fmt.Sprintf("%s-%d", computeCapabilityidPrefix, i+1))

		}
		require.ElementsMatch(t, expectedIDs, receivedIDs)
		allRecords := records

		t.Run("List Success With Filter", func(t *testing.T) {
			records, err := repo.List(ctx, computecapability.ComputeCapabilityListFilters{
				StateIn: []computecapability.ComputeCapabilityState{computecapability.ComputeCapabilityStateActive},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, computecapability.ComputeCapabilityStateActive, record.Status.State)
			}
		})

		t.Run("List with Names Filter", func(t *testing.T) {
			records, err := repo.List(ctx, computecapability.ComputeCapabilityListFilters{
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
			records, err := repo.List(ctx, computecapability.ComputeCapabilityListFilters{
				Limit: 3,
			})
			require.NoError(t, err)
			require.Len(t, records, 3)
		})

		t.Run("List with IncludeDeleted", func(t *testing.T) {
			err = repo.Delete(ctx, allRecords[0].Metadata)
			require.NoError(t, err)

			records, err := repo.List(ctx, computecapability.ComputeCapabilityListFilters{
				IncludeDeleted: true,
			})
			require.NoError(t, err)
			require.Len(t, records, 11)
		})

		t.Run("List with StateNotIn", func(t *testing.T) {
			records, err := repo.List(ctx, computecapability.ComputeCapabilityListFilters{
				StateNotIn: []computecapability.ComputeCapabilityState{computecapability.ComputeCapabilityStateActive},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, computecapability.ComputeCapabilityStateInActive, record.Status.State)
			}
		})

		t.Run("Update State and check version", func(t *testing.T) {
			status := computecapability.ComputeCapabilityStatus{
				State:   computecapability.ComputeCapabilityStatePending,
				Message: "Needs attention",
			}

			err = repo.UpdateState(ctx, allRecords[1].Metadata, status)
			require.NoError(t, err)
			ve := uint64(1)
			records, err := repo.List(ctx, computecapability.ComputeCapabilityListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 1)

			ve += 1
			records, err = repo.List(ctx, computecapability.ComputeCapabilityListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 0)
		})

		t.Run("List by Type", func(t *testing.T) {
			records, err := repo.List(ctx, computecapability.ComputeCapabilityListFilters{
				TypeIn: []string{"GPU"},
			})
			require.NoError(t, err)
			require.Len(t, records, 3)
		})
	})
}
