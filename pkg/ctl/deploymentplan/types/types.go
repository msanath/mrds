package types

import "github.com/msanath/mrds/pkg/ctl/metainstance/types"

// DisplayDeploymentPlan represents the display version of the DeploymentPlanRecord.
//
//go:generate /Users/sanath/projects/gondolf/bin/cligen DisplayDeploymentPlan types types_gen.go
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
	MetaInstances []types.DisplayMetaInstance `json:"meta_instances,omitempty" doNotGen:"true"`
}
