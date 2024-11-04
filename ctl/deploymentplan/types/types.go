package types

import "github.com/msanath/mrds/ctl/metainstance/types"

// DisplayDeploymentPlan represents the display version of the DeploymentPlanRecord.
type DisplayDeploymentPlan struct {
	Metadata                    DisplayMetadata                    `json:"metadata,omitempty"`
	Name                        string                             `json:"name,omitempty" displayName:"Deployment Plan Name" columnTag:"name"`
	Status                      DisplayDeploymentPlanStatus        `json:"status,omitempty"`
	Namespace                   string                             `json:"namespace,omitempty" displayName:"Namespace" columnTag:"namespace"`
	ServiceName                 string                             `json:"service_name,omitempty" displayName:"Service Name" columnTag:"service_name"`
	MatchingComputeCapabilities []DisplayMatchingComputeCapability `json:"matching_compute_capabilities,omitempty"`
	Applications                []DisplayApplication               `json:"applications,omitempty"`
	Deployments                 []DisplayDeployment                `json:"deployments,omitempty"`
	InstanceSummary             DisplayInstanceSummary             `json:"instance_summary,omitempty"`
}

// DisplayMetadata is the display representation of the core.Metadata in NodeRecord
type DisplayMetadata struct {
	ID        string `json:"id,omitempty" displayName:"Node ID" columnTag:"node_id"`
	Version   int    `json:"version,omitempty" displayName:"Version" columnTag:"version"`
	IsDeleted bool   `json:"is_deleted,omitempty" displayName:"Is Deleted" columnTag:"is_deleted" redTexts:"true" greenTexts:"false"`
}

// DisplayDeploymentPlanStatus represents the display version of DeploymentPlanStatus
type DisplayDeploymentPlanStatus struct {
	State   string `json:"state,omitempty" displayName:"Deployment Plan State" columnTag:"state" greenTexts:"DeploymentPlanState_ACTIVE" redTexts:"DeploymentPlanState_INACTIVE,DeploymentPlanState_UNKNOWN"`
	Message string `json:"message,omitempty" displayName:"Status Message" columnTag:"message"`
}

// DisplayMatchingComputeCapability represents the display version of MatchingComputeCapability
type DisplayMatchingComputeCapability struct {
	CapabilityType  string   `json:"capability_type,omitempty" displayName:"Capability Type"`
	Comparator      string   `json:"comparator,omitempty" displayName:"Comparator"`
	CapabilityNames []string `json:"capability_names,omitempty" displayName:"Capability Names"`
}

// DisplayApplication represents the display version of Application
type DisplayApplication struct {
	PayloadName       string                               `json:"payload_name,omitempty" displayName:"Payload Name"`
	Resources         DisplayApplicationResources          `json:"resources,omitempty"`
	Ports             []DisplayApplicationPort             `json:"ports,omitempty"`
	PersistentVolumes []DisplayApplicationPersistentVolume `json:"persistent_volumes,omitempty"`
}

// DisplayApplicationResources represents the resources available to the Application
type DisplayApplicationResources struct {
	Cores  int `json:"cores,omitempty" displayName:"Cores"`
	Memory int `json:"memory,omitempty" displayName:"Memory (MB)"`
}

// DisplayApplicationPort represents the ports used by the Application
type DisplayApplicationPort struct {
	Protocol string `json:"protocol,omitempty" displayName:"Protocol"`
	Port     int    `json:"port,omitempty" displayName:"Port"`
}

// DisplayApplicationPersistentVolume represents the persistent volumes used by the Application
type DisplayApplicationPersistentVolume struct {
	StorageClass string `json:"storage_class,omitempty" displayName:"Storage Class"`
	Capacity     int    `json:"capacity,omitempty" displayName:"Capacity (GB)"`
	MountPath    string `json:"mount_path,omitempty" displayName:"Mount Path"`
}

// DisplayDeployment represents the display version of a Deployment
type DisplayDeployment struct {
	ID                 string                      `json:"id,omitempty" displayName:"Deployment ID"`
	Status             DisplayDeploymentStatus     `json:"status,omitempty"`
	PayloadCoordinates []DisplayPayloadCoordinates `json:"payload_coordinates,omitempty"`
	InstanceCount      int                         `json:"instance_count,omitempty" displayName:"Instance Count"`
}

// DisplayDeploymentStatus represents the display version of DeploymentStatus
type DisplayDeploymentStatus struct {
	State   string `json:"state,omitempty" displayName:"Deployment State" greenTexts:"DeploymentState_COMPLETED" redTexts:"DeploymentState_FAILED,DeploymentState_PENDING"`
	Message string `json:"message,omitempty" displayName:"Status Message"`
}

// DisplayPayloadCoordinates represents the coordinates for a Payload in a Deployment
type DisplayPayloadCoordinates struct {
	PayloadName string `json:"payload_name,omitempty" displayName:"Payload Name"`
	Coordinates string `json:"coordinates,omitempty" displayName:"Coordinates"`
}

// DisplayInstanceSummary represents the summary of instances in a DeploymentPlan
type DisplayInstanceSummary struct {
	NumTotalInstances   int `json:"num_total_instances,omitempty" displayName:"Total Instances"`
	NumRunningInstances int `json:"num_running_instances,omitempty" displayName:"Running Instances"`
	NumPendingInstances int `json:"num_pending_instances,omitempty" displayName:"Pending Instances"`
	NumFailedInstances  int `json:"num_failed_instances,omitempty" displayName:"Failed Instances"`

	NumCreateOperationsPending   int `json:"num_create_operations_pending,omitempty" displayName:"Create Operations Pending"`
	NumStopOperationsPending     int `json:"num_stop_operations_pending,omitempty" displayName:"Stop Operations Pending"`
	NumRestartOperationsPending  int `json:"num_restart_operations_pending,omitempty" displayName:"Restart Operations Pending"`
	NumUpdateOperationsPending   int `json:"num_update_operations_pending,omitempty" displayName:"Update Operations Pending"`
	NumRelocateOperationsPending int `json:"num_relocate_operations_pending,omitempty" displayName:"Relocate Operations Pending"`
	NumDeleteOperationsPending   int `json:"num_delete_operations_pending,omitempty" displayName:"Delete Operations Pending"`

	NumCreateOperationsApproved   int `json:"num_create_operations_approved,omitempty" displayName:"Create Operations Approved"`
	NumStopOperationsApproved     int `json:"num_stop_operations_approved,omitempty" displayName:"Stop Operations Approved"`
	NumRestartOperationsApproved  int `json:"num_restart_operations_approved,omitempty" displayName:"Restart Operations Approved"`
	NumUpdateOperationsApproved   int `json:"num_update_operations_approved,omitempty" displayName:"Update Operations Approved"`
	NumRelocateOperationsApproved int `json:"num_relocate_operations_approved,omitempty" displayName:"Relocate Operations Approved"`
	NumDeleteOperationsApproved   int `json:"num_delete_operations_approved,omitempty" displayName:"Delete Operations Approved"`

	NumCreateOperationsFailed   int `json:"num_create_operations_failed,omitempty" displayName:"Create Operations Failed"`
	NumStopOperationsFailed     int `json:"num_stop_operations_failed,omitempty" displayName:"Stop Operations Failed"`
	NumRestartOperationsFailed  int `json:"num_restart_operations_failed,omitempty" displayName:"Restart Operations Failed"`
	NumUpdateOperationsFailed   int `json:"num_update_operations_failed,omitempty" displayName:"Update Operations Failed"`
	NumRelocateOperationsFailed int `json:"num_relocate_operations_failed,omitempty" displayName:"Relocate Operations Failed"`
	NumDeleteOperationsFailed   int `json:"num_delete_operations_failed,omitempty" displayName:"Delete Operations Failed"`

	NumCreateOperationsSucceeded   int `json:"num_create_operations_succeeded,omitempty" displayName:"Create Operations Succeeded"`
	NumStopOperationsSucceeded     int `json:"num_stop_operations_succeeded,omitempty" displayName:"Stop Operations Succeeded"`
	NumRestartOperationsSucceeded  int `json:"num_restart_operations_succeeded,omitempty" displayName:"Restart Operations Succeeded"`
	NumUpdateOperationsSucceeded   int `json:"num_update_operations_succeeded,omitempty" displayName:"Update Operations Succeeded"`
	NumRelocateOperationsSucceeded int `json:"num_relocate_operations_succeeded,omitempty" displayName:"Relocate Operations Succeeded"`
	NumDeleteOperationsSucceeded   int `json:"num_delete_operations_succeeded,omitempty" displayName:"Delete Operations Succeeded"`

	MetaInstances []types.DisplayMetaInstance `json:"meta_instances,omitempty" doNotGen:"true"`
}
