package computecapability

import (
	"context"

	"github.com/msanath/mrds/ledger/core"
)

// ComputeCapability is a representation of the ComputeCapability of an application.
type ComputeCapabilityRecord struct {
	Metadata core.Metadata           // Metadata is the metadata that identifies the ComputeCapability. It is a combination of the ComputeCapability's name and version.
	Name     string                  // Name is the name of the ComputeCapability.
	Status   ComputeCapabilityStatus // Status is the status of the ComputeCapability.

	Type  string // Type is the type of the ComputeCapability.
	Score uint32 // Score is the score of the ComputeCapability. This is relative to the other ComputeCapabilities of the same type.
}

// ComputeCapabilityState is the state of a ComputeCapability.
type ComputeCapabilityState string

const (
	ComputeCapabilityStateUnknown  ComputeCapabilityState = "ComputeCapabilityState_UNKNOWN"
	ComputeCapabilityStatePending  ComputeCapabilityState = "ComputeCapabilityState_PENDING"
	ComputeCapabilityStateActive   ComputeCapabilityState = "ComputeCapabilityState_ACTIVE"
	ComputeCapabilityStateInActive ComputeCapabilityState = "ComputeCapabilityState_INACTIVE"
)

// ToString returns the string representation of the ComputeCapabilityState.
func (s ComputeCapabilityState) ToString() string {
	return string(s)
}

// FromString converts a string into a ComputeCapabilityState if valid, otherwise returns an error.
func ComputeCapabilityStateFromString(s string) ComputeCapabilityState {
	switch s {
	case string(ComputeCapabilityStatePending):
		return ComputeCapabilityStatePending
	case string(ComputeCapabilityStateActive):
		return ComputeCapabilityStateActive
	case string(ComputeCapabilityStateInActive):
		return ComputeCapabilityStateInActive
	default:
		return ComputeCapabilityState(s) // Unknown state. Return as is.
	}
}

type ComputeCapabilityStatus struct {
	State   ComputeCapabilityState // State is the discrete condition of the resource.
	Message string                 // Message is a human-readable description of the resource's state.
}

// Ledger provides the methods for managing ComputeCapability records.
type Ledger interface {
	// Create creates a new ComputeCapability.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// GetByMetadata retrieves a ComputeCapability by its metadata.
	GetByID(context.Context, string) (*GetResponse, error)
	// GetByName retrieves a ComputeCapability by its name.
	GetByName(context.Context, string) (*GetResponse, error)
	// UpdateStatus updates the state and message of an existing ComputeCapability.
	UpdateStatus(context.Context, *UpdateStateRequest) (*UpdateResponse, error)
	// List returns a list of ComputeCapability that match the provided filters.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Delete deletes a ComputeCapability.
	Delete(context.Context, *DeleteRequest) error
}

// CreateRequest represents the ComputeCapability creation request.
type CreateRequest struct {
	Name  string
	Type  string
	Score uint32
}

// CreateResponse represents the response after creating a new ComputeCapability.
type CreateResponse struct {
	Record ComputeCapabilityRecord
}

// UpdateStateRequest represents the request to update the state and message of a ComputeCapability.
type UpdateStateRequest struct {
	Metadata core.Metadata
	Status   ComputeCapabilityStatus
}

// GetResponse represents the response for fetching a ComputeCapability.
type GetResponse struct {
	Record ComputeCapabilityRecord
}

// UpdateResponse represents the response after updating the state of a ComputeCapability.
type UpdateResponse struct {
	Record ComputeCapabilityRecord
}

// ListRequest represents the request to list ComputeCapabilitys with filters.
type ListRequest struct {
	Filters ComputeCapabilityListFilters
}

// ComputeCapabilityFilters contains filters for querying the ComputeCapability table.
type ComputeCapabilityListFilters struct {
	IDIn       []string // IN condition
	NameIn     []string // IN condition
	VersionGte *uint64  // Greater than or equal condition
	VersionLte *uint64  // Less than or equal condition
	VersionEq  *uint64  // Equal condition

	IncludeDeleted bool   // IncludeDeleted indicates whether to include soft-deleted records in the result.
	Limit          uint32 // Limit is the maximum number of records to return.

	StateIn    []ComputeCapabilityState
	StateNotIn []ComputeCapabilityState
	TypeIn     []string
}

// ListResponse represents the response to a list request.
type ListResponse struct {
	Records []ComputeCapabilityRecord
}

type DeleteRequest struct {
	Metadata core.Metadata
}
