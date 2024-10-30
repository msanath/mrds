package types

import "time"

// DisplayNode is a display representation of the NodeRecord
//
//go:generate /Users/sanath/projects/gondolf/bin/cligen DisplayNode types types_gen.go
type DisplayNode struct {
	Metadata           DisplayMetadata      `json:"metadata,omitempty"`
	Name               string               `json:"name,omitempty" displayName:"Node Name" columnTag:"name"`
	Status             DisplayNodeStatus    `json:"status,omitempty"`
	ClusterID          string               `json:"cluster_id,omitempty" displayName:"Cluster ID" columnTag:"cluster_id"`
	UpdateDomain       string               `json:"update_domain,omitempty" displayName:"Update Domain" columnTag:"update_domain"`
	TotalResources     DisplayResources     `json:"total_resources,omitempty"`
	SystemReserved     DisplayResources     `json:"system_reserved,omitempty" displayName:"System Reserved Resources"`
	RemainingResources DisplayResources     `json:"remaining_resources,omitempty" displayName:"Remaining Resources"`
	LocalVolumes       []DisplayLocalVolume `json:"local_volumes,omitempty"`
	CapabilityIDs      []string             `json:"capability_ids,omitempty" displayName:"Capability IDs" columnTag:"capability_ids"`
	Disruptions        []DisplayDisruption  `json:"disruptions,omitempty"`
}

// DisplayMetadata is the display representation of the core.Metadata in NodeRecord
type DisplayMetadata struct {
	ID        string `json:"id,omitempty" displayName:"Node ID" columnTag:"node_id"`
	Version   int    `json:"version,omitempty" displayName:"Version" columnTag:"version"`
	IsDeleted bool   `json:"is_deleted,omitempty" displayName:"Is Deleted" columnTag:"is_deleted" redTexts:"true" greenTexts:"false"`
}

// DisplayNodeStatus represents the display of the NodeStatus
type DisplayNodeStatus struct {
	State   string `json:"state,omitempty" displayName:"Node State" columnTag:"state" greenTexts:"NodeState_ALLOCATED,NodeState_UNALLOCATED" redTexts:"NodeState_EVICTED,NodeState_SANITIZING"`
	Message string `json:"message,omitempty" displayName:"Status Message" columnTag:"message"`
}

// DisplayResources represents the resources available to the Node
type DisplayResources struct {
	Cores  int `json:"cores,omitempty" displayName:"Cores" columnTag:"cores"`
	Memory int `json:"memory,omitempty" displayName:"Memory (MB)" columnTag:"memory"`
}

// DisplayLocalVolume represents each local volume attached to the Node
type DisplayLocalVolume struct {
	MountPath       string `json:"mount_path,omitempty" displayName:"Mount Path"`
	StorageClass    string `json:"storage_class,omitempty" displayName:"Storage Class"`
	StorageCapacity int    `json:"storage_capacity,omitempty" displayName:"Storage Capacity (GB)"`
}

// DisplayDisruption represents the display version of Disruption in NodeRecord
type DisplayDisruption struct {
	ID          string                  `json:"disruption_id,omitempty" displayName:"Disruption ID"`
	ShouldEvict bool                    `json:"should_evict,omitempty" displayName:"Should Evict" redTexts:"true" greenTexts:"false"`
	StartTime   time.Time               `json:"start_time,omitempty" displayName:"Start Time"`
	Status      DisplayDisruptionStatus `json:"status,omitempty"`
}

// DisplayDisruptionStatus represents the display version of the disruption status
type DisplayDisruptionStatus struct {
	State   string `json:"state,omitempty" displayName:"Disruption State" greenTexts:"DisruptionState_COMPLETED" redTexts:"DisruptionState_SCHEDULED,DisruptionState_UNKNOWN"`
	Message string `json:"message,omitempty" displayName:"Disruption Status Message"`
}
