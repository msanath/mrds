package getter

import (
	"context"

	"github.com/msanath/mrds/ctl/metainstance/types"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"google.golang.org/grpc"
)

type Getter struct {
	deploymentPlansClient mrdspb.DeploymentPlansClient
	nodesClient           mrdspb.NodesClient
	metaInstancesClient   mrdspb.MetaInstancesClient
}

func NewGetter(conn *grpc.ClientConn) *Getter {
	return &Getter{
		deploymentPlansClient: mrdspb.NewDeploymentPlansClient(conn),
		nodesClient:           mrdspb.NewNodesClient(conn),
		metaInstancesClient:   mrdspb.NewMetaInstancesClient(conn),
	}
}

func (g *Getter) GetDisplayMetaInstances(ctx context.Context, m []*mrdspb.MetaInstance) ([]types.DisplayMetaInstance, error) {
	var displayMetaInstances []types.DisplayMetaInstance

	for _, metaInstance := range m {
		displayMetaInstance, err := g.GetDisplayMetaInstance(ctx, metaInstance)
		if err != nil {
			return nil, err
		}
		displayMetaInstances = append(displayMetaInstances, displayMetaInstance)
	}

	return displayMetaInstances, nil
}

func (g *Getter) GetDisplayMetaInstance(ctx context.Context, m *mrdspb.MetaInstance) (types.DisplayMetaInstance, error) {
	deploymentResp, err := g.deploymentPlansClient.GetByID(ctx, &mrdspb.GetDeploymentPlanByIDRequest{Id: m.DeploymentPlanId})
	if err != nil {
		return types.DisplayMetaInstance{}, err
	}

	displayMetaInstance := types.DisplayMetaInstance{
		Metadata: types.DisplayMetadata{
			ID:      m.GetMetadata().GetId(),
			Version: int(m.GetMetadata().GetVersion()),
		},
		Name:               m.GetName(),
		DeploymentPlanName: deploymentResp.Record.Name,
		DeploymentID:       m.GetDeploymentId(),
		Status: types.DisplayMetaInstanceStatus{
			State:   m.GetStatus().GetState().String(),
			Message: m.GetStatus().GetMessage(),
		},
	}

	// Convert RuntimeInstances
	for _, instance := range m.GetRuntimeInstances() {
		nodeResp, err := g.nodesClient.GetByID(ctx, &mrdspb.GetNodeByIDRequest{Id: instance.NodeId})
		if err != nil {
			return types.DisplayMetaInstance{}, err
		}
		displayMetaInstance.RuntimeInstances = append(displayMetaInstance.RuntimeInstances, types.DisplayRuntimeInstance{
			ID:       instance.GetId(),
			NodeName: nodeResp.Record.Name,
			IsActive: instance.GetIsActive(),
			Status: types.DisplayRuntimeInstanceStatus{
				State:   instance.GetStatus().GetState().String(),
				Message: instance.GetStatus().GetMessage(),
			},
		})
	}

	// Convert Operations
	for _, operation := range m.GetOperations() {
		displayMetaInstance.Operations = append(displayMetaInstance.Operations, types.DisplayOperation{
			ID:       operation.GetId(),
			Type:     operation.GetType().String(),
			IntentID: operation.GetIntentId(),
			Status: types.DisplayOperationStatus{
				State:   operation.GetStatus().GetState().String(),
				Message: operation.GetStatus().GetMessage(),
			},
		})
	}

	return displayMetaInstance, nil
}
