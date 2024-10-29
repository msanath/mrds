
package deploymentplan

import (
	"context"

	"github.com/msanath/mrds/internal/ledger/core"
)

// DeploymentPlan is a representation of the DeploymentPlan of an application.
type DeploymentPlanRecord struct {
	Metadata core.Metadata // Metadata is the metadata that identifies the DeploymentPlan. It is a combination of the DeploymentPlan's name and version.
	Name     string        // Name is the name of the DeploymentPlan.
	Status   DeploymentPlanStatus // Status is the status of the DeploymentPlan.
}

// DeploymentPlanState is the state of a DeploymentPlan.
type DeploymentPlanState string

const (
	DeploymentPlanStateUnknown  DeploymentPlanState = "DeploymentPlanState_UNKNOWN"
	DeploymentPlanStatePending  DeploymentPlanState = "DeploymentPlanState_PENDING"
	DeploymentPlanStateActive   DeploymentPlanState = "DeploymentPlanState_ACTIVE"
	DeploymentPlanStateInActive DeploymentPlanState = "DeploymentPlanState_INACTIVE"
)

// ToString returns the string representation of the DeploymentPlanState.
func (s DeploymentPlanState) ToString() string {
	return string(s)
}

// FromString converts a string into a DeploymentPlanState if valid, otherwise returns an error.
func DeploymentPlanStateFromString(s string) DeploymentPlanState {
	switch s {
	case string(DeploymentPlanStatePending):
		return DeploymentPlanStatePending
	case string(DeploymentPlanStateActive):
		return DeploymentPlanStateActive
	case string(DeploymentPlanStateInActive):
		return DeploymentPlanStateInActive
	default:
		return DeploymentPlanState(s) // Unknown state. Return as is.
	}
}

type DeploymentPlanStatus struct {
	State   DeploymentPlanState // State is the discrete condition of the resource.
	Message string       // Message is a human-readable description of the resource's state.
}

// Ledger provides the methods for managing DeploymentPlan records.
type Ledger interface {
	// Create creates a new DeploymentPlan.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// GetByMetadata retrieves a DeploymentPlan by its metadata.
	GetByMetadata(context.Context, *core.Metadata) (*GetResponse, error)
	// GetByName retrieves a DeploymentPlan by its name.
	GetByName(context.Context, string) (*GetResponse, error)
	// UpdateStatus updates the state and message of an existing DeploymentPlan.
	UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateResponse, error)
	// List returns a list of DeploymentPlan that match the provided filters.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Delete deletes a DeploymentPlan.
	Delete(context.Context, *DeleteRequest) error
}

// CreateRequest represents the DeploymentPlan creation request.
type CreateRequest struct {
	Name string
}

// CreateResponse represents the response after creating a new DeploymentPlan.
type CreateResponse struct {
	Record DeploymentPlanRecord
}

// UpdateStatusRequest represents the request to update the state and message of a DeploymentPlan.
type UpdateStatusRequest struct {
	Metadata core.Metadata
	Status   DeploymentPlanStatus
}

// GetResponse represents the response for fetching a DeploymentPlan.
type GetResponse struct {
	Record DeploymentPlanRecord
}

// UpdateResponse represents the response after updating the state of a DeploymentPlan.
type UpdateResponse struct {
	Record DeploymentPlanRecord
}

// ListRequest represents the request to list DeploymentPlans with filters.
type ListRequest struct {
	Filters DeploymentPlanListFilters
}

// DeploymentPlanFilters contains filters for querying the DeploymentPlan table.
type DeploymentPlanListFilters struct {
	IDIn       []string // IN condition
	NameIn     []string // IN condition
	VersionGte *uint64  // Greater than or equal condition
	VersionLte *uint64  // Less than or equal condition
	VersionEq  *uint64  // Equal condition

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
