package deployment

import (
	"context"

	"github.com/msanath/mrds/internal/ledger/core"
)

// Deployment is a representation of the Deployment of an application.
type DeploymentRecord struct {
	Metadata core.Metadata    // Metadata is the metadata that identifies the Deployment. It is a combination of the Deployment's name and version.
	Name     string           // Name is the name of the Deployment.
	Status   DeploymentStatus // Status is the status of the Deployment.
}

// DeploymentState is the state of a Deployment.
type DeploymentState string

const (
	DeploymentStateUnknown  DeploymentState = "DeploymentState_UNKNOWN"
	DeploymentStatePending  DeploymentState = "DeploymentState_PENDING"
	DeploymentStateActive   DeploymentState = "DeploymentState_ACTIVE"
	DeploymentStateInActive DeploymentState = "DeploymentState_INACTIVE"
)

// ToString returns the string representation of the DeploymentState.
func (s DeploymentState) ToString() string {
	return string(s)
}

// FromString converts a string into a DeploymentState if valid, otherwise returns an error.
func DeploymentStateFromString(s string) DeploymentState {
	switch s {
	case string(DeploymentStatePending):
		return DeploymentStatePending
	case string(DeploymentStateActive):
		return DeploymentStateActive
	case string(DeploymentStateInActive):
		return DeploymentStateInActive
	default:
		return DeploymentState(s) // Unknown state. Return as is.
	}
}

type DeploymentStatus struct {
	State   DeploymentState // State is the discrete condition of the resource.
	Message string          // Message is a human-readable description of the resource's state.
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
	// Delete deletes a Deployment.
	Delete(context.Context, *DeleteRequest) error
}

// CreateRequest represents the Deployment creation request.
type CreateRequest struct {
	Name string
}

// CreateResponse represents the response after creating a new Deployment.
type CreateResponse struct {
	Record DeploymentRecord
}

// UpdateStatusRequest represents the request to update the state and message of a Deployment.
type UpdateStatusRequest struct {
	Metadata core.Metadata
	Status   DeploymentStatus
}

// GetResponse represents the response for fetching a Deployment.
type GetResponse struct {
	Record DeploymentRecord
}

// UpdateResponse represents the response after updating the state of a Deployment.
type UpdateResponse struct {
	Record DeploymentRecord
}

// ListRequest represents the request to list Deployments with filters.
type ListRequest struct {
	Filters DeploymentListFilters
}

// DeploymentFilters contains filters for querying the Deployment table.
type DeploymentListFilters struct {
	IDIn       []string // IN condition
	NameIn     []string // IN condition
	VersionGte *uint64  // Greater than or equal condition
	VersionLte *uint64  // Less than or equal condition
	VersionEq  *uint64  // Equal condition

	IncludeDeleted bool   // IncludeDeleted indicates whether to include soft-deleted records in the result.
	Limit          uint32 // Limit is the maximum number of records to return.

	StateIn    []DeploymentState
	StateNotIn []DeploymentState
}

// ListResponse represents the response to a list request.
type ListResponse struct {
	Records []DeploymentRecord
}

type DeleteRequest struct {
	Metadata core.Metadata
}
