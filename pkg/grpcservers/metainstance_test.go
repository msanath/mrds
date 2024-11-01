package grpcservers_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/msanath/mrds/gen/api/mrdspb"
	servertest "github.com/msanath/mrds/pkg/grpcservers/test"

	"github.com/stretchr/testify/require"
)

func TestMetaInstanceServer(t *testing.T) {
	ts, err := servertest.NewTestServer()
	require.NoError(t, err)
	defer ts.Close()

	client := mrdspb.NewMetaInstancesClient(ts.Conn())
	ctx := context.Background()

	// Create a deployment plan, a deployment and a node
	deploymentPlanClient := mrdspb.NewDeploymentPlansClient(ts.Conn())
	planResp, err := deploymentPlanClient.Create(ctx, &mrdspb.CreateDeploymentPlanRequest{
		Name:        "test-deployment-plan",
		Namespace:   "test-namespace",
		ServiceName: "test-service",
		Applications: []*mrdspb.Application{
			{
				PayloadName: "test-payload",
				Resources: &mrdspb.ApplicationResources{
					Cores:  1,
					Memory: 200,
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
	})
	require.NoError(t, err)
	require.NotNil(t, planResp)

	deploymentResp, err := deploymentPlanClient.AddDeployment(ctx, &mrdspb.AddDeploymentRequest{
		Metadata:     planResp.Record.Metadata,
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
	require.NotNil(t, deploymentResp)

	nodeClient := mrdspb.NewNodesClient(ts.Conn())
	nodeResp, err := nodeClient.Create(ctx, &mrdspb.CreateNodeRequest{
		Name:         "test-Node",
		UpdateDomain: "test-domain",
		TotalResources: &mrdspb.Resources{
			Cores:  64,
			Memory: 512,
		},
		SystemReservedResources: &mrdspb.Resources{
			Cores:  4,
			Memory: 32,
		},
		LocalVolumes: []*mrdspb.NodeLocalVolume{
			{
				MountPath:       "/var/lib/docker",
				StorageClass:    "SSD",
				StorageCapacity: 1024,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, nodeResp)

	// create
	resp, err := client.Create(ctx, &mrdspb.CreateMetaInstanceRequest{
		Name:             "test-metaInstance",
		DeploymentPlanId: planResp.Record.Metadata.Id,
		DeploymentId:     deploymentResp.Record.Deployments[0].Id,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "test-metaInstance", resp.Record.Name)

	// get by id
	getResp, err := client.GetByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{Id: resp.Record.Metadata.Id})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, "test-metaInstance", getResp.Record.Name)

	// get by name
	getByNameResp, err := client.GetByName(ctx, &mrdspb.GetMetaInstanceByNameRequest{Name: "test-metaInstance"})
	require.NoError(t, err)
	require.NotNil(t, getByNameResp)
	require.Equal(t, "test-metaInstance", getByNameResp.Record.Name)

	// update
	updateResp, err := client.UpdateStatus(ctx, &mrdspb.UpdateMetaInstanceStatusRequest{
		Metadata: resp.Record.Metadata,
		Status: &mrdspb.MetaInstanceStatus{
			State:   mrdspb.MetaInstanceState_MetaInstanceState_RUNNING,
			Message: "test-message",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-metaInstance", updateResp.Record.Name)
	require.Equal(t, mrdspb.MetaInstanceState_MetaInstanceState_RUNNING, updateResp.Record.Status.State)

	// Add runtime instance
	updateResp, err = client.AddRuntimeInstance(ctx, &mrdspb.AddRuntimeInstanceRequest{
		Metadata: updateResp.Record.Metadata,
		RuntimeInstance: &mrdspb.RuntimeInstance{
			Id:     uuid.New().String(),
			NodeId: nodeResp.Record.Metadata.Id,
			Status: &mrdspb.RuntimeInstanceStatus{
				State:   mrdspb.RuntimeInstanceState_RuntimeState_RUNNING,
				Message: "test-message",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-metaInstance", updateResp.Record.Name)
	require.Len(t, updateResp.Record.RuntimeInstances, 1)

	// Update runtime status
	updateResp, err = client.UpdateRuntimeStatus(ctx, &mrdspb.UpdateRuntimeStatusRequest{
		Metadata:          updateResp.Record.Metadata,
		RuntimeInstanceId: updateResp.Record.RuntimeInstances[0].Id,
		Status: &mrdspb.RuntimeInstanceStatus{
			State:   mrdspb.RuntimeInstanceState_RuntimeState_TERMINATED,
			Message: "test-message",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-metaInstance", updateResp.Record.Name)
	require.Equal(t, mrdspb.RuntimeInstanceState_RuntimeState_TERMINATED, updateResp.Record.RuntimeInstances[0].Status.State)

	// Remove runtime instance
	updateResp, err = client.RemoveRuntimeInstance(ctx, &mrdspb.RemoveRuntimeInstanceRequest{
		Metadata:          updateResp.Record.Metadata,
		RuntimeInstanceId: updateResp.Record.RuntimeInstances[0].Id,
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-metaInstance", updateResp.Record.Name)
	require.Len(t, updateResp.Record.RuntimeInstances, 0)

	// Add operation
	updateResp, err = client.AddOperation(ctx, &mrdspb.AddOperationRequest{
		Metadata: updateResp.Record.Metadata,
		Operation: &mrdspb.Operation{
			Id:       uuid.New().String(),
			Type:     "deploy",
			IntentId: "intent-id",
			Status: &mrdspb.OperationStatus{
				State:   mrdspb.OperationState_OperationState_APPROVED,
				Message: "test-message",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-metaInstance", updateResp.Record.Name)
	require.Len(t, updateResp.Record.Operations, 1)

	// Update operation status
	updateResp, err = client.UpdateOperationStatus(ctx, &mrdspb.UpdateOperationStatusRequest{
		Metadata:    updateResp.Record.Metadata,
		OperationId: updateResp.Record.Operations[0].Id,
		Status: &mrdspb.OperationStatus{
			State:   mrdspb.OperationState_OperationState_SUCCEEDED,
			Message: "test-message",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-metaInstance", updateResp.Record.Name)
	require.Equal(t, mrdspb.OperationState_OperationState_SUCCEEDED, updateResp.Record.Operations[0].Status.State)

	// Remove operation
	updateResp, err = client.RemoveOperation(ctx, &mrdspb.RemoveOperationRequest{
		Metadata:    updateResp.Record.Metadata,
		OperationId: updateResp.Record.Operations[0].Id,
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-metaInstance", updateResp.Record.Name)
	require.Len(t, updateResp.Record.Operations, 0)

	// list
	listResp, err := client.List(ctx, &mrdspb.ListMetaInstanceRequest{
		StateIn: []mrdspb.MetaInstanceState{mrdspb.MetaInstanceState_MetaInstanceState_RUNNING},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Records, 1)

	// Delete
	_, err = client.Delete(ctx, &mrdspb.DeleteMetaInstanceRequest{Metadata: updateResp.Record.Metadata})
	require.NoError(t, err)
}
