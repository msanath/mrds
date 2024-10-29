package deploymentplan

import (
	"context"

	"github.com/msanath/mrds/internal/ledger/core"
)

// Deployment is a representation of a workload which is expected to be deployed.
type DeploymentPlanRecord struct {
	Metadata core.Metadata        // Metadata is the metadata that identifies the Deployment. It is a combination of the Deployment's name and version.
	Name     string               // Name is the name for this Deployment.
	Status   DeploymentPlanStatus // Status is the status of the Deployment.

	Namespace                   string                      // Namespace is the namespace of the Deployment.
	ServiceName                 string                      // ServiceName is the name of the service associated with the Deployment. Certs will be issued for this service.
	MatchingComputeCapabilities []MatchingComputeCapability // MatchingCapabilities is a list of capabilities that the Deployment requires.
	Applications                []Application               // Applications is a list of applications that the Deployment requires.

	Deployments []Deployment // Deployment is an instantiation of the DeploymentPlan.
}

type DeploymentPlanStatus struct {
	State   DeploymentPlanState // State is the discrete condition of the resource.
	Message string              // Message is a human-readable description of the resource's state.
}

// DeploymentPlanState is the state of a Deployment.
type DeploymentPlanState string

const (
	DeploymentPlanStateUnknown  DeploymentPlanState = "DeploymentPlanState_UNKNOWN"
	DeploymentPlanStateActive   DeploymentPlanState = "DeploymentPlanState_ACTIVE"
	DeploymentPlanStateInactive DeploymentPlanState = "DeploymentPlanState_INACTIVE"
)

type MatchingComputeCapability struct {
	CapabilityType  string
	Comparator      ComparatorType
	CapabilityNames []string
}

type ComparatorType string

const (
	ComparatorTypeIn    ComparatorType = "ComparatorType_IN"
	ComparatorTypeNotIn ComparatorType = "ComparatorType_NOT_IN"
	ComparatorTypeGte   ComparatorType = "ComparatorType_GTE"
	ComparatorTypeLte   ComparatorType = "ComparatorType_LTE"
)

type Application struct {
	PayloadName       string
	Resources         ApplicationResources
	Ports             []ApplicationPort
	PersistentVolumes []ApplicationPersistentVolume
}

type ApplicationResources struct {
	Cores  uint32
	Memory uint32
}

type ApplicationPort struct {
	Protocol string
	Port     uint32
}

type ApplicationPersistentVolume struct {
	StorageClass string
	Capacity     uint32
	MountPath    string
}

type Deployment struct {
	ID                 string               // ID is the ID of the DeploymentPlan.
	Status             DeploymentStatus     // Status is the status of the Deployment.
	PayloadCoordinates []PayloadCoordinates // PayloadCoordinates is a list of coordinates for the payloads that the Deployment requires.
	InstanceCount      uint32               // InstanceCount is the number of instances of the Deployment.
}

// DeploymentState is the state of a Deployment.
type DeploymentState string

const (
	DeploymentStateUnknown    DeploymentState = "DeploymentState_UNKNOWN"
	DeploymentStatePending    DeploymentState = "DeploymentState_PENDING"
	DeploymentStateInProgress DeploymentState = "DeploymentState_IN_PROGRESS"
	DeploymentStateFailed     DeploymentState = "DeploymentState_FAILED"
	DeploymentStatePaused     DeploymentState = "DeploymentState_PAUSED"
	DeploymentStateCompleted  DeploymentState = "DeploymentState_COMPLETED"
)

type DeploymentStatus struct {
	State   DeploymentState // State is the discrete condition of the resource.
	Message string          // Message is a human-readable description of the resource's state.
}

type PayloadCoordinates struct {
	PayloadName string
	Coordinates map[string]string
}

// Ledger provides the methods for managing Deployment records.
type Ledger interface {
	// Create creates a new Deployment.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// GetByMetadata retrieves a Deployment by its metadata.
	GetByMetadata(context.Context, *core.Metadata) (*GetResponse, error)
	// GetByName retrieves a Deployment by its name.
	GetByName(context.Context, string) (*GetResponse, error)
	// UpdateStatus updates the state and message of an existing Deployment.
	UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateResponse, error)
	// List returns a list of Deployment that match the provided filters.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Delete soft-deletes a Deployment.
	Delete(context.Context, *DeleteRequest) error

	AddDeployment(context.Context, *AddDeploymentRequest) (*UpdateResponse, error)
	UpdateDeploymentStatus(context.Context, *UpdateDeploymentStatusRequest) (*UpdateResponse, error)
}

// CreateRequest represents the Deployment creation request.
type CreateRequest struct {
	Name                        string
	Namespace                   string                      // Namespace is the namespace of the Deployment.
	ServiceName                 string                      // ServiceName is the name of the service associated with the Deployment. Certs will be issued for this service.
	MatchingComputeCapabilities []MatchingComputeCapability // MatchingCapabilities is a list of capabilities that the Deployment requires.
	Applications                []Application               // Applications is a list of applications that the Deployment requires.
}

// CreateResponse represents the response after creating a new Deployment.
type CreateResponse struct {
	Record DeploymentPlanRecord
}

// UpdateStatusRequest represents the request to update the state and message of a Deployment.
type UpdateStatusRequest struct {
	Metadata core.Metadata
	Status   DeploymentPlanStatus
}

// GetResponse represents the response for fetching a Deployment.
type GetResponse struct {
	Record DeploymentPlanRecord
}

// UpdateResponse represents the response after updating the state of a Deployment.
type UpdateResponse struct {
	Record DeploymentPlanRecord
}

// ListRequest represents the request to list Deployments with filters.
type ListRequest struct {
	Filters DeploymentPlanListFilters
}

// DeploymentFilters contains filters for querying the Deployment table.
type DeploymentPlanListFilters struct {
	IDIn          []string // IN condition
	NameIn        []string // IN condition
	VersionGte    *uint64  // Greater than or equal condition
	VersionLte    *uint64  // Less than or equal condition
	VersionEq     *uint64  // Equal condition
	PayloadNameIn []string // IN condition
	ServiceNameIn []string // IN condition

	DeploymentPlanIDIn     []string // IN condition
	DeploymentPlanStatusIn []DeploymentState

	IncludeDeleted bool   // IncludeDeleted indicates whether to include soft-deleted records in the result.
	Limit          uint32 // Limit is the maximum number of records to return.

	StateIn    []DeploymentPlanState
	StateNotIn []DeploymentPlanState
}

// ListResponse represents the response to a list request.
type ListResponse struct {
	Records []DeploymentPlanRecord
}

type DeleteRequest struct {
	Metadata core.Metadata
}

type AddDeploymentRequest struct {
	Metadata           core.Metadata
	DeploymentID       string
	PayloadCoordinates []PayloadCoordinates
	InstanceCount      uint32
}

type UpdateDeploymentStatusRequest struct {
	Metadata     core.Metadata
	DeploymentID string
	Status       DeploymentStatus
}
