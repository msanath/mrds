package grpcservers_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/msanath/mrds/gen/api/mrdspb"
	servertest "github.com/msanath/mrds/pkg/grpcservers/test"

	"github.com/stretchr/testify/require"
)

func TestDeploymentPlanServer(t *testing.T) {
	ts, err := servertest.NewTestServer()
	require.NoError(t, err)
	defer ts.Close()

	client := mrdspb.NewDeploymentPlansClient(ts.Conn())
	ctx := context.Background()

	req := &mrdspb.CreateDeploymentPlanRequest{
		Name:        "test-deployment-plan",
		Namespace:   "test-namespace",
		ServiceName: "test-service",
		Applications: []*mrdspb.Application{
			{
				PayloadName: "test-payload",
				Resources: &mrdspb.ApplicationResources{
					Cores:  1,
					Memory: 1023,
				},
				PersistentVolumes: []*mrdspb.ApplicationPersistentVolume{
					{
						MountPath:    "/data",
						Capacity:     1024,
						StorageClass: "SSD",
					},
				},
				Ports: []*mrdspb.ApplicationPort{
					{
						Protocol: "TCP",
						Port:     80,
					},
				},
			},
		},
	}
	// create
	resp, err := client.Create(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "test-deployment-plan", resp.Record.Name)
	require.Equal(t, "test-namespace", resp.Record.Namespace)
	require.Equal(t, "test-service", resp.Record.ServiceName)
	require.Len(t, resp.Record.Applications, 1)

	// get by id
	getResp, err := client.GetByID(ctx, &mrdspb.GetDeploymentPlanByIDRequest{Id: resp.Record.Metadata.Id})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, "test-deployment-plan", getResp.Record.Name)

	// get by name
	getByNameResp, err := client.GetByName(ctx, &mrdspb.GetDeploymentPlanByNameRequest{Name: "test-deployment-plan"})
	require.NoError(t, err)
	require.NotNil(t, getByNameResp)
	require.Equal(t, "test-deployment-plan", getByNameResp.Record.Name)

	// update status
	updateResp, err := client.UpdateStatus(ctx, &mrdspb.UpdateDeploymentPlanStatusRequest{
		Metadata: resp.Record.Metadata,
		Status: &mrdspb.DeploymentPlanStatus{
			State:   mrdspb.DeploymentPlanState_DeploymentPlanState_INACTIVE,
			Message: "Deployment is now active",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-deployment-plan", updateResp.Record.Name)
	require.Equal(t, mrdspb.DeploymentPlanState_DeploymentPlanState_INACTIVE, updateResp.Record.Status.State)

	// Add deployment
	updateResp, err = client.AddDeployment(ctx, &mrdspb.AddDeploymentRequest{
		Metadata:     updateResp.Record.Metadata,
		DeploymentId: uuid.New().String(),
		PayloadCoordinates: []*mrdspb.PayloadCoordinates{
			{
				PayloadName: "test-payload",
				Coordinates: map[string]string{
					"key": "value",
				},
			},
		},
		InstanceCount: 1,
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-deployment-plan", updateResp.Record.Name)
	require.Len(t, updateResp.Record.Deployments, 1)

	// Update deployment status
	updateResp, err = client.UpdateDeploymentStatus(ctx, &mrdspb.UpdateDeploymentStatusRequest{
		Metadata:     updateResp.Record.Metadata,
		DeploymentId: updateResp.Record.Deployments[0].Id,
		Status: &mrdspb.DeploymentStatus{
			State:   mrdspb.DeploymentState_DeploymentState_IN_PROGRESS,
			Message: "Deployment completed",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-deployment-plan", updateResp.Record.Name)
	require.Equal(t, mrdspb.DeploymentState_DeploymentState_IN_PROGRESS, updateResp.Record.Deployments[0].Status.State)

	// list
	listResp, err := client.List(ctx, &mrdspb.ListDeploymentPlanRequest{
		Filters: &mrdspb.DeploymentPlanListFilters{
			NameIn: []string{"test-deployment-plan"},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Records, 1)

	// Delete
	_, err = client.Delete(ctx, &mrdspb.DeleteDeploymentPlanRequest{Metadata: updateResp.Record.Metadata})
	require.NoError(t, err)

	// Get deleted by name
	_, err = client.GetByName(ctx, &mrdspb.GetDeploymentPlanByNameRequest{Name: "test-deployment-plan"})
	require.Error(t, err)
}
