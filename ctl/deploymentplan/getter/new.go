package getter

import (
	"context"
	"encoding/json"

	"github.com/msanath/mrds/ctl/deploymentplan/types"
	"github.com/msanath/mrds/ctl/metainstance/getter"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"google.golang.org/grpc"
)

type Getter struct {
	deploymentPlansClient mrdspb.DeploymentPlansClient
	nodesClient           mrdspb.NodesClient
	metaInstancesClient   mrdspb.MetaInstancesClient

	metaInstancesGetter *getter.Getter
}

func NewGetter(conn *grpc.ClientConn) *Getter {

	metaInstancesGetter := getter.NewGetter(conn)

	return &Getter{
		deploymentPlansClient: mrdspb.NewDeploymentPlansClient(conn),
		nodesClient:           mrdspb.NewNodesClient(conn),
		metaInstancesClient:   mrdspb.NewMetaInstancesClient(conn),
		metaInstancesGetter:   metaInstancesGetter,
	}
}

func (g *Getter) ConvertGRPCDeploymentPlanToDisplayDeploymentPlan(ctx context.Context, d *mrdspb.DeploymentPlanRecord) (types.DisplayDeploymentPlan, error) {
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

	listMetaInstancesResp, err := g.metaInstancesClient.List(ctx, &mrdspb.ListMetaInstanceRequest{
		DeploymentPlanIdIn: []string{d.GetMetadata().GetId()},
	})
	if err != nil {
		return displayDeploymentPlan, err
	}

	for _, instance := range listMetaInstancesResp.GetRecords() {
		for _, operation := range instance.Operations {
			switch operation.Type {
			case mrdspb.OperationType_OperationType_CREATE:
				switch operation.Status.State {
				case mrdspb.OperationState_OperationState_PENDING_APPROVAL:
					displayDeploymentPlan.InstanceSummary.NumCreateOperationsPending++
				case mrdspb.OperationState_OperationState_APPROVED:
					displayDeploymentPlan.InstanceSummary.NumCreateOperationsApproved++
				case mrdspb.OperationState_OperationState_FAILED:
					displayDeploymentPlan.InstanceSummary.NumCreateOperationsFailed++
				case mrdspb.OperationState_OperationState_SUCCEEDED:
					displayDeploymentPlan.InstanceSummary.NumCreateOperationsSucceeded++
				}
			case mrdspb.OperationType_OperationType_STOP:
				switch operation.Status.State {
				case mrdspb.OperationState_OperationState_PENDING_APPROVAL:
					displayDeploymentPlan.InstanceSummary.NumStopOperationsPending++
				case mrdspb.OperationState_OperationState_APPROVED:
					displayDeploymentPlan.InstanceSummary.NumStopOperationsApproved++
				case mrdspb.OperationState_OperationState_FAILED:
					displayDeploymentPlan.InstanceSummary.NumStopOperationsFailed++
				case mrdspb.OperationState_OperationState_SUCCEEDED:
					displayDeploymentPlan.InstanceSummary.NumStopOperationsSucceeded++
				}
			case mrdspb.OperationType_OperationType_RESTART:
				switch operation.Status.State {
				case mrdspb.OperationState_OperationState_PENDING_APPROVAL:
					displayDeploymentPlan.InstanceSummary.NumRestartOperationsPending++
				case mrdspb.OperationState_OperationState_APPROVED:
					displayDeploymentPlan.InstanceSummary.NumRestartOperationsApproved++
				case mrdspb.OperationState_OperationState_FAILED:
					displayDeploymentPlan.InstanceSummary.NumRestartOperationsFailed++
				case mrdspb.OperationState_OperationState_SUCCEEDED:
					displayDeploymentPlan.InstanceSummary.NumRestartOperationsSucceeded++
				}
			case mrdspb.OperationType_OperationType_UPDATE:
				switch operation.Status.State {
				case mrdspb.OperationState_OperationState_PENDING_APPROVAL:
					displayDeploymentPlan.InstanceSummary.NumUpdateOperationsPending++
				case mrdspb.OperationState_OperationState_APPROVED:
					displayDeploymentPlan.InstanceSummary.NumUpdateOperationsApproved++
				case mrdspb.OperationState_OperationState_FAILED:
					displayDeploymentPlan.InstanceSummary.NumUpdateOperationsFailed++
				case mrdspb.OperationState_OperationState_SUCCEEDED:
					displayDeploymentPlan.InstanceSummary.NumUpdateOperationsSucceeded++
				}
			case mrdspb.OperationType_OperationType_RELOCATE:
				switch operation.Status.State {
				case mrdspb.OperationState_OperationState_PENDING:
					displayDeploymentPlan.InstanceSummary.NumRelocateOperationsPending++
				case mrdspb.OperationState_OperationState_APPROVED:
					displayDeploymentPlan.InstanceSummary.NumRelocateOperationsApproved++
				case mrdspb.OperationState_OperationState_FAILED:
					displayDeploymentPlan.InstanceSummary.NumRelocateOperationsFailed++
				case mrdspb.OperationState_OperationState_SUCCEEDED:
					displayDeploymentPlan.InstanceSummary.NumRelocateOperationsSucceeded++
				}
			case mrdspb.OperationType_OperationType_DELETE:
				switch operation.Status.State {
				case mrdspb.OperationState_OperationState_PENDING:
					displayDeploymentPlan.InstanceSummary.NumDeleteOperationsPending++
				case mrdspb.OperationState_OperationState_APPROVED:
					displayDeploymentPlan.InstanceSummary.NumDeleteOperationsApproved++
				case mrdspb.OperationState_OperationState_FAILED:
					displayDeploymentPlan.InstanceSummary.NumDeleteOperationsFailed++
				case mrdspb.OperationState_OperationState_SUCCEEDED:
					displayDeploymentPlan.InstanceSummary.NumDeleteOperationsSucceeded++
				}
			}
		}
		for _, runtimeInstance := range instance.RuntimeInstances {
			switch runtimeInstance.Status.State {
			case mrdspb.RuntimeInstanceState_RuntimeState_RUNNING:
				displayDeploymentPlan.InstanceSummary.NumRunningInstances++
			case mrdspb.RuntimeInstanceState_RuntimeState_PENDING:
				displayDeploymentPlan.InstanceSummary.NumPendingInstances++
			case mrdspb.RuntimeInstanceState_RuntimeState_FAILED:
				displayDeploymentPlan.InstanceSummary.NumFailedInstances++
			}
		}
	}
	displayDeploymentPlan.InstanceSummary.NumTotalInstances = len(listMetaInstancesResp.GetRecords())

	displayMetaInstances, err := g.metaInstancesGetter.GetDisplayMetaInstances(ctx, listMetaInstancesResp.GetRecords())
	if err != nil {
		return displayDeploymentPlan, err
	}

	displayDeploymentPlan.InstanceSummary.MetaInstances = displayMetaInstances
	return displayDeploymentPlan, nil
}
