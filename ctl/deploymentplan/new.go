package deploymentplan

import (
	"encoding/json"

	"github.com/msanath/mrds/ctl/deploymentplan/types"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/spf13/cobra"
)

func NewDeploymentPlanCmd() *cobra.Command {
	cmd := cobra.Command{
		Use: "deployment",
	}

	cmd.AddCommand(newDeploymentPlanCreateCmd())
	cmd.AddCommand(newDeploymentPlanListCmd())
	cmd.AddCommand(newDeploymentPlanShowCmd())
	cmd.AddCommand(newAddDeploymentCmd())
	cmd.AddCommand(newCancelDeploymentCmd())
	cmd.AddCommand(newApproveOperationCmd())
	cmd.AddCommand(newOperateOption())

	return &cmd
}

func convertGRPCDeploymentPlanToDisplayDeploymentPlan(d *mrdspb.DeploymentPlanRecord) types.DisplayDeploymentPlan {
	displayDeploymentPlan := types.DisplayDeploymentPlan{
		Metadata: types.DisplayMetadata{
			ID:      d.GetMetadata().GetId(),
			Version: int(d.GetMetadata().GetVersion()),
		},
		Name:        d.GetName(),
		Namespace:   d.GetNamespace(),
		ServiceName: d.GetServiceName(),
		Status: types.DisplayDeploymentPlanStatus{
			State:   d.GetStatus().GetState().String(),
			Message: d.GetStatus().GetMessage(),
		},
	}

	// Convert MatchingComputeCapabilities
	for _, capability := range d.GetMatchingComputeCapabilities() {
		displayDeploymentPlan.MatchingComputeCapabilities = append(displayDeploymentPlan.MatchingComputeCapabilities, types.DisplayMatchingComputeCapability{
			CapabilityType:  capability.GetCapabilityType(),
			Comparator:      capability.GetComparator().String(),
			CapabilityNames: capability.GetCapabilityNames(),
		})
	}

	// Convert Applications
	for _, app := range d.GetApplications() {
		displayApp := types.DisplayApplication{
			PayloadName: app.GetPayloadName(),
			Resources: types.DisplayApplicationResources{
				Cores:  int(app.GetResources().GetCores()),
				Memory: int(app.GetResources().GetMemory()),
			},
		}

		// Convert Ports
		for _, port := range app.GetPorts() {
			displayApp.Ports = append(displayApp.Ports, types.DisplayApplicationPort{
				Protocol: port.GetProtocol(),
				Port:     int(port.GetPort()),
			})
		}

		// Convert PersistentVolumes
		for _, volume := range app.GetPersistentVolumes() {
			displayApp.PersistentVolumes = append(displayApp.PersistentVolumes, types.DisplayApplicationPersistentVolume{
				StorageClass: volume.GetStorageClass(),
				Capacity:     int(volume.GetCapacity()),
				MountPath:    volume.GetMountPath(),
			})
		}

		displayDeploymentPlan.Applications = append(displayDeploymentPlan.Applications, displayApp)
	}

	// Convert Deployments
	for _, deployment := range d.GetDeployments() {
		displayDeployment := types.DisplayDeployment{
			ID:            deployment.GetId(),
			InstanceCount: int(deployment.GetInstanceCount()),
			Status: types.DisplayDeploymentStatus{
				State:   deployment.GetStatus().GetState().String(),
				Message: deployment.GetStatus().GetMessage(),
			},
		}

		// Convert PayloadCoordinates
		for _, coord := range deployment.GetPayloadCoordinates() {
			jsonStr, _ := json.Marshal(coord.GetCoordinates())
			displayDeployment.PayloadCoordinates = append(displayDeployment.PayloadCoordinates, types.DisplayPayloadCoordinates{
				PayloadName: coord.GetPayloadName(),
				Coordinates: string(jsonStr),
			})
		}

		displayDeploymentPlan.Deployments = append(displayDeploymentPlan.Deployments, displayDeployment)
	}

	return displayDeploymentPlan
}
