package sqlstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/deploymentplan"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"
	"github.com/msanath/mrds/internal/sqlstorage"

	"github.com/msanath/gondolf/pkg/simplesql/test"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const deploymentPlanidPrefix = "deploymentplan"

func TestDeploymentPlanRecordLifecycle(t *testing.T) {
	db, err := test.NewTestSQLiteDB()
	require.NoError(t, err)

	storage, err := sqlstorage.NewSQLStorage(db, true)
	require.NoError(t, err)

	testRecord := deploymentplan.DeploymentPlanRecord{
		Metadata: core.Metadata{
			ID:      fmt.Sprintf("%s1", deploymentPlanidPrefix),
			Version: 1,
		},
		Name: fmt.Sprintf("%s1", deploymentPlanidPrefix),
		Status: deploymentplan.DeploymentPlanStatus{
			State:   deploymentplan.DeploymentPlanStateActive,
			Message: "",
		},
		Namespace:   "default",
		ServiceName: "service1",
		MatchingComputeCapabilities: []deploymentplan.MatchingComputeCapability{
			{
				CapabilityType:  "capability1",
				Comparator:      deploymentplan.ComparatorTypeIn,
				CapabilityNames: []string{"name1", "name2"},
			},
			{
				CapabilityType:  "capability2",
				Comparator:      deploymentplan.ComparatorTypeNotIn,
				CapabilityNames: []string{"name3", "name4"},
			},
		},
		Applications: []deploymentplan.Application{
			{
				PayloadName: "app1",
				Resources: deploymentplan.ApplicationResources{
					Cores:  1,
					Memory: 200,
				},
				Ports: []deploymentplan.ApplicationPort{
					{
						Protocol: "TCP",
						Port:     8080,
					},
					{
						Protocol: "UDP",
						Port:     8081,
					},
				},
			},
			{
				PayloadName: "app2",
				Resources: deploymentplan.ApplicationResources{
					Cores:  1,
					Memory: 200,
				},
				PersistentVolumes: []deploymentplan.ApplicationPersistentVolume{
					{
						StorageClass: "SSD",
						Capacity:     1024,
						MountPath:    "/data",
					},
				},
			},
		},
	}
	repo := storage.DeploymentPlan

	testDeploymentPlanCRUD(t, repo, testRecord)
}

func testDeploymentPlanCRUD(t *testing.T, repo deploymentplan.Repository, testRecord deploymentplan.DeploymentPlanRecord) {
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
		require.Equal(t, testRecord.Name, receivedRecord.Name)
		require.Equal(t, testRecord.Status, receivedRecord.Status)
		require.Equal(t, testRecord.Namespace, receivedRecord.Namespace)
		require.Equal(t, testRecord.ServiceName, receivedRecord.ServiceName)
		require.ElementsMatch(t, testRecord.Applications, receivedRecord.Applications)
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
		status := deploymentplan.DeploymentPlanStatus{
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

	t.Run("Insert Deployment Success", func(t *testing.T) {
		err = repo.InsertDeployment(ctx, testRecord.Metadata, deploymentplan.Deployment{
			ID: "deployment1",
			Status: deploymentplan.DeploymentStatus{
				State:   deploymentplan.DeploymentStateInProgress,
				Message: "Deployment in progress",
			},
			InstanceCount: 20,
			PayloadCoordinates: []deploymentplan.PayloadCoordinates{
				{
					PayloadName: "app1",
					Coordinates: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
				{
					PayloadName: "app2",
					Coordinates: map[string]string{
						"key3": "value3",
						"key4": "value4",
					},
				},
			},
		})
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.Deployments, 1)
		require.Equal(t, "deployment1", updatedRecord.Deployments[0].ID)
		require.Equal(t, deploymentplan.DeploymentStateInProgress, updatedRecord.Deployments[0].Status.State)
		require.Equal(t, "Deployment in progress", updatedRecord.Deployments[0].Status.Message)
		require.Equal(t, uint32(20), updatedRecord.Deployments[0].InstanceCount)
		testRecord = updatedRecord
	})

	t.Run("Update Deployment State Success", func(t *testing.T) {
		err = repo.UpdateDeploymentStatus(ctx, testRecord.Metadata, "deployment1", deploymentplan.DeploymentStatus{
			State:   deploymentplan.DeploymentStateFailed,
			Message: "Deployment failed",
		})
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Len(t, updatedRecord.Deployments, 1)
		require.Equal(t, "deployment1", updatedRecord.Deployments[0].ID)
		require.Equal(t, deploymentplan.DeploymentStateFailed, updatedRecord.Deployments[0].Status.State)
		require.Equal(t, "Deployment failed", updatedRecord.Deployments[0].Status.Message)
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
			newRecord.Metadata.ID = fmt.Sprintf("%s-%d", deploymentPlanidPrefix, i+1)
			newRecord.Metadata.Version = 0
			newRecord.Name = fmt.Sprintf("%s-%d", deploymentPlanidPrefix, i+1)
			newRecord.Status.State = deploymentplan.DeploymentPlanStateActive
			newRecord.Status.Message = fmt.Sprintf("%s-%d is active", deploymentPlanidPrefix, i+1)
			newRecord.Applications = testRecord.Applications

			if (i+1)%2 == 0 {
				newRecord.Status.State = deploymentplan.DeploymentPlanStateInactive
				newRecord.Status.Message = fmt.Sprintf("%s-%d is inactive", deploymentPlanidPrefix, i+1)
			}

			err = repo.Insert(ctx, newRecord)
			require.NoError(t, err)
		}
	})

	t.Run("List", func(t *testing.T) {
		records, err := repo.List(ctx, deploymentplan.DeploymentPlanListFilters{})
		require.NoError(t, err)
		require.Len(t, records, 10)

		receivedIDs := []string{}
		for _, record := range records {
			receivedIDs = append(receivedIDs, record.Metadata.ID)

		}
		expectedIDs := []string{}
		for i := range 10 {
			expectedIDs = append(expectedIDs, fmt.Sprintf("%s-%d", deploymentPlanidPrefix, i+1))

		}
		require.ElementsMatch(t, expectedIDs, receivedIDs)
		allRecords := records

		t.Run("List Success With Filter", func(t *testing.T) {
			records, err := repo.List(ctx, deploymentplan.DeploymentPlanListFilters{
				StateIn: []deploymentplan.DeploymentPlanState{deploymentplan.DeploymentPlanStateActive},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, deploymentplan.DeploymentPlanStateActive, record.Status.State)
				require.Equal(t, 2, len(record.Applications))
			}
		})

		t.Run("List with Names Filter", func(t *testing.T) {
			records, err := repo.List(ctx, deploymentplan.DeploymentPlanListFilters{
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
			records, err := repo.List(ctx, deploymentplan.DeploymentPlanListFilters{
				Limit: 3,
			})
			require.NoError(t, err)
			require.Len(t, records, 3)
		})

		t.Run("List with IncludeDeleted", func(t *testing.T) {
			err = repo.Delete(ctx, allRecords[0].Metadata)
			require.NoError(t, err)

			records, err := repo.List(ctx, deploymentplan.DeploymentPlanListFilters{
				IncludeDeleted: true,
			})
			require.NoError(t, err)
			require.Len(t, records, 11)
		})

		t.Run("List with StateNotIn", func(t *testing.T) {
			records, err := repo.List(ctx, deploymentplan.DeploymentPlanListFilters{
				StateNotIn: []deploymentplan.DeploymentPlanState{deploymentplan.DeploymentPlanStateActive},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, deploymentplan.DeploymentPlanStateInactive, record.Status.State)
			}
		})

		t.Run("Update State and check version", func(t *testing.T) {
			status := deploymentplan.DeploymentPlanStatus{
				State:   deploymentplan.DeploymentPlanStateInactive,
				Message: "Needs attention",
			}

			err = repo.UpdateState(ctx, allRecords[1].Metadata, status)
			require.NoError(t, err)
			ve := uint64(1)
			records, err := repo.List(ctx, deploymentplan.DeploymentPlanListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 1)

			ve += 1
			records, err = repo.List(ctx, deploymentplan.DeploymentPlanListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 0)
		})
	})
}
