package node

import (
	"time"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/node/types"
	"github.com/spf13/cobra"
)

func NewNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage nodes",
	}

	cmd.AddCommand(newNodeCreateCmd())
	cmd.AddCommand(newNodeListCmd())
	cmd.AddCommand(newNodeShowCmd())
	cmd.AddCommand(newAddToClusterCmd())

	return cmd
}

func convertGRPCNodeToDisplayNode(n *mrdspb.Node) types.DisplayNode {
	displayNode := types.DisplayNode{
		Metadata: types.DisplayMetadata{
			ID:      n.GetMetadata().GetId(),
			Version: int(n.GetMetadata().GetVersion()),
		},
		Name:         n.GetName(),
		UpdateDomain: n.GetUpdateDomain(),
		Status: types.DisplayNodeStatus{
			State:   n.GetStatus().GetState().String(),
			Message: n.GetStatus().GetMessage(),
		},
		CapabilityIDs: n.GetCapabilityIds(),
		TotalResources: types.DisplayResources{
			Cores:  int(n.GetTotalResources().GetCores()),
			Memory: int(n.GetTotalResources().GetMemory()),
		},
		SystemReserved: types.DisplayResources{
			Cores:  int(n.GetSystemReservedResources().GetCores()),
			Memory: int(n.GetSystemReservedResources().GetMemory()),
		},
		RemainingResources: types.DisplayResources{
			Cores:  int(n.GetRemainingResources().GetCores()),
			Memory: int(n.GetRemainingResources().GetMemory()),
		},
		ClusterID: n.GetClusterId(),
	}

	for _, volume := range n.GetLocalVolumes() {
		displayNode.LocalVolumes = append(displayNode.LocalVolumes, types.DisplayLocalVolume{
			MountPath:       volume.GetMountPath(),
			StorageClass:    volume.GetStorageClass(),
			StorageCapacity: int(volume.GetStorageCapacity()),
		})
	}

	for _, disruption := range n.GetDisruptions() {
		displayNode.Disruptions = append(displayNode.Disruptions, types.DisplayDisruption{
			ID:          disruption.GetId(),
			ShouldEvict: disruption.GetShouldEvict(),
			StartTime:   time.Now(), //  FIXME
			Status: types.DisplayDisruptionStatus{
				State:   disruption.GetStatus().GetState().String(),
				Message: disruption.GetStatus().GetMessage(),
			},
		})
	}

	return displayNode
}
