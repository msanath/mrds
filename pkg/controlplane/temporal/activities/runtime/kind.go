// This file implements a local k8s kind cluster runtime, which is meant to be used for testing purposes.
package runtime

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type RuntimeActivity interface {
	StartInstance(ctx context.Context, req *RuntimeActivityParams) (*mrdspb.UpdateMetaInstanceResponse, error)
	StopInstance(ctx context.Context, req *RuntimeActivityParams) (*mrdspb.UpdateMetaInstanceResponse, error)
}

type RuntimeActivityParams struct {
	MetaInstanceID    string
	RuntimeInstanceID string
}

type KindRuntime struct {
	metaInstancesClient  mrdspb.MetaInstancesClient
	deploymentPlanClient mrdspb.DeploymentPlansClient
	nodesClient          mrdspb.NodesClient
}

func NewKindRuntime(
	metaInstancesClient mrdspb.MetaInstancesClient,
	deploymentPlanClient mrdspb.DeploymentPlansClient,
	nodesClient mrdspb.NodesClient,
	registry worker.Registry,
) RuntimeActivity {
	a := &KindRuntime{
		metaInstancesClient:  metaInstancesClient,
		deploymentPlanClient: deploymentPlanClient,
		nodesClient:          nodesClient,
	}

	registry.RegisterActivity(a.StartInstance)
	registry.RegisterActivity(a.StopInstance)

	return a
}

func (k *KindRuntime) StartInstance(ctx context.Context, req *RuntimeActivityParams) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("StartInstance", "request", req)

	runtimeDetails, err := k.getRuntimeDetails(ctx, req.MetaInstanceID, req.RuntimeInstanceID)
	if err != nil {
		return nil, err
	}

	// Update the runtime state to running.
	updateResp, err := k.metaInstancesClient.UpdateRuntimeStatus(ctx, &mrdspb.UpdateRuntimeStatusRequest{
		Metadata:          runtimeDetails.MetaInstance.Metadata,
		RuntimeInstanceId: req.RuntimeInstanceID,
		Status: &mrdspb.RuntimeInstanceStatus{
			State: mrdspb.RuntimeInstanceState_RuntimeState_RUNNING,
		},
	})

	return updateResp, err
}

func (k *KindRuntime) StopInstance(ctx context.Context, req *RuntimeActivityParams) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("StopInstance", "request", req)

	runtimeDetails, err := k.getRuntimeDetails(ctx, req.MetaInstanceID, req.RuntimeInstanceID)
	if err != nil {
		return nil, err
	}

	// Update the runtime state to stopped.
	updateResp, err := k.metaInstancesClient.UpdateRuntimeStatus(ctx, &mrdspb.UpdateRuntimeStatusRequest{
		Metadata:          runtimeDetails.MetaInstance.Metadata,
		RuntimeInstanceId: req.RuntimeInstanceID,
		Status: &mrdspb.RuntimeInstanceStatus{
			State: mrdspb.RuntimeInstanceState_RuntimeState_TERMINATED,
		},
	})

	return updateResp, err
}

type runtimeDetails struct {
	MetaInstance    *mrdspb.MetaInstance
	RuntimeInstance *mrdspb.RuntimeInstance
	DeploymentPlan  *mrdspb.DeploymentPlanRecord
	Node            *mrdspb.Node
}

func (k *KindRuntime) getRuntimeDetails(ctx context.Context, metaInstanceID string, runtimeInstanceID string) (*runtimeDetails, error) {
	// Get the corresponding meta instance
	metaInstanceGetResp, err := k.metaInstancesClient.GetByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{
		Id: metaInstanceID,
	})
	if err != nil {
		return nil, err
	}

	found := false
	runtimeInstance := &mrdspb.RuntimeInstance{}
	for _, ri := range metaInstanceGetResp.Record.RuntimeInstances {
		if ri.Id == runtimeInstanceID {
			found = true
			runtimeInstance = ri
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("runtime instance %s not found in meta instance %s", metaInstanceID, metaInstanceID)
	}

	// Get the node
	nodeGetResp, err := k.nodesClient.GetByID(ctx, &mrdspb.GetNodeByIDRequest{
		Id: runtimeInstance.NodeId,
	})
	if err != nil {
		return nil, err
	}
	node := nodeGetResp.Record
	activity.GetLogger(ctx).Info("Node", "node", node)

	// Get the corresponding deployment plan
	deploymentPlanGetResp, err := k.deploymentPlanClient.GetByID(ctx, &mrdspb.GetDeploymentPlanByIDRequest{
		Id: metaInstanceGetResp.Record.DeploymentPlanId,
	})
	if err != nil {
		return nil, err
	}
	plan := deploymentPlanGetResp.Record
	activity.GetLogger(ctx).Info("Deployment Plan", "plan", plan)

	return &runtimeDetails{
		MetaInstance:    metaInstanceGetResp.Record,
		RuntimeInstance: runtimeInstance,
		DeploymentPlan:  plan,
		Node:            node,
	}, nil
}
