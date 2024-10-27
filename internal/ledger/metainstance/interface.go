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
}

// MetaInstanceState is the state of a MetaInstance.
type MetaInstanceState string

const (
	MetaInstanceStateUnknown  MetaInstanceState = "MetaInstanceState_UNKNOWN"
	MetaInstanceStatePending  MetaInstanceState = "MetaInstanceState_PENDING"
	MetaInstanceStateActive   MetaInstanceState = "MetaInstanceState_ACTIVE"
	MetaInstanceStateInActive MetaInstanceState = "MetaInstanceState_INACTIVE"
)

// ToString returns the string representation of the MetaInstanceState.
func (s MetaInstanceState) ToString() string {
	return string(s)
}

// FromString converts a string into a MetaInstanceState if valid, otherwise returns an error.
func MetaInstanceStateFromString(s string) MetaInstanceState {
	switch s {
	case string(MetaInstanceStatePending):
		return MetaInstanceStatePending
	case string(MetaInstanceStateActive):
		return MetaInstanceStateActive
	case string(MetaInstanceStateInActive):
		return MetaInstanceStateInActive
	default:
		return MetaInstanceState(s) // Unknown state. Return as is.
	}
}

type MetaInstanceStatus struct {
	State   MetaInstanceState // State is the discrete condition of the resource.
	Message string            // Message is a human-readable description of the resource's state.
}

// Ledger provides the methods for managing MetaInstance records.
type Ledger interface {
	// Create creates a new MetaInstance.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// GetByMetadata retrieves a MetaInstance by its metadata.
	GetByMetadata(context.Context, *core.Metadata) (*GetResponse, error)
	// GetByName retrieves a MetaInstance by its name.
	GetByName(context.Context, string) (*GetResponse, error)
	// UpdateStatus updates the state and message of an existing MetaInstance.
	UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateResponse, error)
	// List returns a list of MetaInstance that match the provided filters.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Delete deletes a MetaInstance.
	Delete(context.Context, *DeleteRequest) error
}

// CreateRequest represents the MetaInstance creation request.
type CreateRequest struct {
	Name string
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
	IDIn       []string // IN condition
	NameIn     []string // IN condition
	VersionGte *uint64  // Greater than or equal condition
	VersionLte *uint64  // Less than or equal condition
	VersionEq  *uint64  // Equal condition

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
