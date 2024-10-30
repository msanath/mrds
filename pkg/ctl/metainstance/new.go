package metainstance

import (
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/metainstance/types"
	"github.com/spf13/cobra"
)

func NewInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Manage instances of a deployment plan",
	}

	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newMetaInstanceListCmd())
	cmd.AddCommand(newAddRuntimeInstanceCmd())
	cmd.AddCommand(newStopInstanceCmd())
	cmd.AddCommand(newRestartInstanceCmd())
	cmd.AddCommand(newSwapInstanceCmd())
	cmd.AddCommand(newApproveOperationCmd())
	cmd.AddCommand(newCompleteOperationCmd())

	return cmd
}

func convertGRPCMetaInstanceToDisplayMetaInstance(m *mrdspb.MetaInstance) types.DisplayMetaInstance {
	displayMetaInstance := types.DisplayMetaInstance{
		Metadata: types.DisplayMetadata{
			ID:      m.GetMetadata().GetId(),
			Version: int(m.GetMetadata().GetVersion()),
		},
		Name:             m.GetName(),
		DeploymentPlanID: m.GetDeploymentPlanId(),
		DeploymentID:     m.GetDeploymentId(),
		Status: types.DisplayMetaInstanceStatus{
			State:   m.GetStatus().GetState().String(),
			Message: m.GetStatus().GetMessage(),
		},
	}

	// Convert RuntimeInstances
	for _, instance := range m.GetRuntimeInstances() {
		displayMetaInstance.RuntimeInstances = append(displayMetaInstance.RuntimeInstances, types.DisplayRuntimeInstance{
			ID:       instance.GetId(),
			NodeID:   instance.GetNodeId(),
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
			Type:     operation.GetType(),
			IntentID: operation.GetIntentId(),
			Status: types.DisplayOperationStatus{
				State:   operation.GetStatus().GetState().String(),
				Message: operation.GetStatus().GetMessage(),
			},
		})
	}

	return displayMetaInstance
}
