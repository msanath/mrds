package types

// DisplayMetaInstance represents the display version of the MetaInstanceRecord.
//
//go:generate /Users/sanath/projects/gondolf/bin/cligen DisplayMetaInstance types types_gen.go
type DisplayMetaInstance struct {
	Metadata           DisplayMetadata           `json:"metadata,omitempty"`
	Name               string                    `json:"name,omitempty" displayName:"Meta Instance Name" columnTag:"name"`
	Status             DisplayMetaInstanceStatus `json:"status,omitempty"`
	DeploymentPlanName string                    `json:"deployment_plan_name,omitempty" displayName:"Deployment Plan Name" columnTag:"deployment_plan_name"`
	DeploymentID       string                    `json:"deployment_id,omitempty" displayName:"Deployment ID" columnTag:"deployment_id"`
	RuntimeInstances   []DisplayRuntimeInstance  `json:"runtime_instances,omitempty"`
	Operations         []DisplayOperation        `json:"operations,omitempty"`
}

// DisplayMetadata is the display version of core.Metadata in MetaInstanceRecord.
type DisplayMetadata struct {
	ID        string `json:"id,omitempty" displayName:"ID" columnTag:"id"`
	Version   int    `json:"version,omitempty" displayName:"Version" columnTag:"version"`
	IsDeleted bool   `json:"is_deleted,omitempty" displayName:"Is Deleted" columnTag:"is_deleted" redTexts:"true" greenTexts:"false"`
}

// DisplayMetaInstanceStatus represents the display version of MetaInstanceStatus.
type DisplayMetaInstanceStatus struct {
	State   string `json:"state,omitempty" displayName:"State" columnTag:"state" greenTexts:"MetaInstanceState_RUNNING" redTexts:"MetaInstanceState_TERMINATED,MetaInstanceState_PENDING_ALLOCATION"`
	Message string `json:"message,omitempty" displayName:"Status Message" columnTag:"message"`
}

// DisplayRuntimeInstance represents the display version of RuntimeInstance.
type DisplayRuntimeInstance struct {
	ID       string                       `json:"id,omitempty" displayName:"Instance ID"`
	NodeName string                       `json:"node_name,omitempty" displayName:"Node Name"`
	IsActive bool                         `json:"is_active,omitempty" displayName:"Is Active" redTexts:"false" greenTexts:"true"`
	Status   DisplayRuntimeInstanceStatus `json:"status,omitempty"`
}

// DisplayRuntimeInstanceStatus represents the display version of RuntimeInstanceStatus.
type DisplayRuntimeInstanceStatus struct {
	State   string `json:"state,omitempty" displayName:"Runtime State" greenTexts:"RuntimeState_RUNNING" redTexts:"RuntimeState_TERMINATED,RuntimeState_PENDING"`
	Message string `json:"message,omitempty" displayName:"Status Message"`
}

// DisplayOperation represents the display version of Operation.
type DisplayOperation struct {
	ID       string                 `json:"id,omitempty" displayName:"Operation ID"`
	Type     string                 `json:"type,omitempty" displayName:"Operation Type"`
	IntentID string                 `json:"intent_id,omitempty" displayName:"Intent ID"`
	Status   DisplayOperationStatus `json:"status,omitempty"`
}

// DisplayOperationStatus represents the display version of OperationStatus.
type DisplayOperationStatus struct {
	State   string `json:"state,omitempty" displayName:"Operation State" yellowTexts:"OperationState_PREPARING,OperationState_PENDING_APPROVAL,OperationState_APPROVED" greenTexts:"OperationState_SUCCEEDED" redTexts:"OperationState_FAILED"`
	Message string `json:"message,omitempty" displayName:"Status Message"`
}
