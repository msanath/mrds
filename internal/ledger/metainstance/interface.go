package metainstance

import (
	"context"

	"github.com/msanath/mrds/internal/ledger/core"
)

// MetaInstance is a representation of the MetaInstance of an application.
type MetaInstanceRecord struct {
	Metadata core.Metadata      // Metadata is the metadata that identifies the MetaInstance. It is a combination of the MetaInstance's name and version.
	Name     string             // Name is the name of the MetaInstance.
	Status   MetaInstanceStatus // Status is the status of the MetaInstance.

	DeploymentPlanID string // DeploymentPlanID is the ID of the DeploymentPlan to which the MetaInstance belongs.
	DeploymentID     string // DeploymentID is the ID of the Deployment which this MetaInstance is currently responsible for.

	RuntimeInstances []RuntimeInstance // RuntimeInstances is a list of runtime instances that are currently running on the MetaInstance.
	Operations       []Operation       // Operations is a list of operations that are currently pending on the MetaInstance.
}

// MetaInstanceState is the state of a MetaInstance.
type MetaInstanceState string

const (
	MetaInstanceStateUnknown           MetaInstanceState = "MetaInstanceState_UNKNOWN"
	MetaInstanceStatePendingAllocation MetaInstanceState = "MetaInstanceState_PENDING_ALLOCATION"
	MetaInstanceStateRunning           MetaInstanceState = "MetaInstanceState_RUNNING"
	MetaInstanceStateTerminated        MetaInstanceState = "MetaInstanceState_TERMINATED"
)

type MetaInstanceStatus struct {
	State   MetaInstanceState // State is the discrete condition of the resource.
	Message string            // Message is a human-readable description of the resource's state.
}

type RuntimeInstance struct {
	ID       string
	NodeID   string
	IsActive bool
	Status   RuntimeInstanceStatus
}

type RuntimeInstanceStatus struct {
	State   RuntimeInstanceState
	Message string
}

type RuntimeInstanceState string

const (
	RuntimeStateUnknown    RuntimeInstanceState = "RuntimeState_UNKNOWN"
	RuntimeStatePending    RuntimeInstanceState = "RuntimeState_PENDING"
	RuntimeStateRunning    RuntimeInstanceState = "RuntimeState_RUNNING"
	RuntimeStateTerminated RuntimeInstanceState = "RuntimeState_TERMINATED"
)

type Operation struct {
	ID       string // The unique ID of the operation.
	Type     string // The type of operation.
	IntentID string // The ID of the intent that triggered this operation.
	Status   OperationStatus
}

type OperationStatus struct {
	State   OperationState
	Message string
}

type OperationState string

const (
	OperationStateUnknown         OperationState = "OperationState_UNKNOWN"
	OperationStatePreparing       OperationState = "OperationState_PREPARING"
	OperationStatePendingApproval OperationState = "OperationState_PENDING_APPROVAL"
	OperationStateApproved        OperationState = "OperationState_APPROVED"
	OperationStateSucceeded       OperationState = "OperationState_SUCCEEDED"
	OperationStateFailed          OperationState = "OperationState_FAILED"
)

// Ledger provides the methods for managing MetaInstance records.
type Ledger interface {
	// Create creates a new MetaInstance.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// GetByID retrieves a MetaInstance by its ID
	GetByID(context.Context, string) (*GetResponse, error)
	// GetByName retrieves a MetaInstance by its name.
	GetByName(context.Context, string) (*GetResponse, error)
	// UpdateStatus updates the state and message of an existing MetaInstance.
	UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateResponse, error)
	// UpdateDeploymentID updates the DeploymentID of an existing MetaInstance.
	UpdateDeploymentID(context.Context, *UpdateDeploymentIDRequest) (*UpdateResponse, error)
	// List returns a list of MetaInstance that match the provided filters.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Delete deletes a MetaInstance.
	Delete(context.Context, *DeleteRequest) error

	AddRuntimeInstance(context.Context, *AddRuntimeInstanceRequest) (*UpdateResponse, error)
	UpdateRuntimeStatus(context.Context, *UpdateRuntimeStatusRequest) (*UpdateResponse, error)
	RemoveRuntimeInstance(context.Context, *RemoveRuntimeInstanceRequest) (*UpdateResponse, error)

	AddOperation(context.Context, *AddOperationRequest) (*UpdateResponse, error)
	UpdateOperationStatus(context.Context, *UpdateOperationStatusRequest) (*UpdateResponse, error)
	RemoveOperation(context.Context, *RemoveOperationRequest) (*UpdateResponse, error)
}

// CreateRequest represents the MetaInstance creation request.
type CreateRequest struct {
	Name             string
	DeploymentPlanID string
	DeploymentID     string
}

// CreateResponse represents the response after creating a new MetaInstance.
type CreateResponse struct {
	Record MetaInstanceRecord
}

// UpdateStatusRequest represents the request to update the state and message of a MetaInstance.
type UpdateStatusRequest struct {
	Metadata core.Metadata
	Status   MetaInstanceStatus
}

type UpdateDeploymentIDRequest struct {
	Metadata     core.Metadata
	DeploymentID string
}

// GetResponse represents the response for fetching a MetaInstance.
type GetResponse struct {
	Record MetaInstanceRecord
}

// UpdateResponse represents the response after updating the state of a MetaInstance.
type UpdateResponse struct {
	Record MetaInstanceRecord
}

// ListRequest represents the request to list MetaInstances with filters.
type ListRequest struct {
	Filters MetaInstanceListFilters
}

// MetaInstanceFilters contains filters for querying the MetaInstance table.
type MetaInstanceListFilters struct {
	IDIn               []string // IN condition
	NameIn             []string // IN condition
	VersionGte         *uint64  // Greater than or equal condition
	VersionLte         *uint64  // Less than or equal condition
	VersionEq          *uint64  // Equal condition
	DeploymentIDIn     []string // IN condition
	DeploymentPlanIDIn []string // IN condition

	IncludeDeleted bool   // IncludeDeleted indicates whether to include soft-deleted records in the result.
	Limit          uint32 // Limit is the maximum number of records to return.

	StateIn    []MetaInstanceState
	StateNotIn []MetaInstanceState
}

// ListResponse represents the response to a list request.
type ListResponse struct {
	Records []MetaInstanceRecord
}

type DeleteRequest struct {
	Metadata core.Metadata
}

type AddRuntimeInstanceRequest struct {
	Metadata        core.Metadata
	RuntimeInstance RuntimeInstance
}

type UpdateRuntimeStatusRequest struct {
	Metadata          core.Metadata
	RuntimeInstanceID string
	Status            RuntimeInstanceStatus
}

type RemoveRuntimeInstanceRequest struct {
	Metadata          core.Metadata
	RuntimeInstanceID string
}

type AddOperationRequest struct {
	Metadata  core.Metadata
	Operation Operation
}

type UpdateOperationStatusRequest struct {
	Metadata    core.Metadata
	OperationID string
	Status      OperationStatus
}

type RemoveOperationRequest struct {
	Metadata    core.Metadata
	OperationID string
}
