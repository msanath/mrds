package cluster

import (
	"context"

	"github.com/msanath/mrds/ledger/core"
)

// Cluster is a representation of the Cluster of an application.
type ClusterRecord struct {
	Metadata core.Metadata // Metadata is the metadata that identifies the Cluster. It is a combination of the Cluster's name and version.
	Name     string        // Name is the name of the Cluster.
	Status   ClusterStatus // Status is the status of the Cluster.
}

// ClusterState is the state of a Cluster.
type ClusterState string

const (
	ClusterStateUnknown  ClusterState = "ClusterState_UNKNOWN"
	ClusterStatePending  ClusterState = "ClusterState_PENDING"
	ClusterStateActive   ClusterState = "ClusterState_ACTIVE"
	ClusterStateInActive ClusterState = "ClusterState_INACTIVE"
)

// ToString returns the string representation of the ClusterState.
func (s ClusterState) ToString() string {
	return string(s)
}

// FromString converts a string into a ClusterState if valid, otherwise returns an error.
func ClusterStateFromString(s string) ClusterState {
	switch s {
	case string(ClusterStatePending):
		return ClusterStatePending
	case string(ClusterStateActive):
		return ClusterStateActive
	case string(ClusterStateInActive):
		return ClusterStateInActive
	default:
		return ClusterState(s) // Unknown state. Return as is.
	}
}

type ClusterStatus struct {
	State   ClusterState // State is the discrete condition of the resource.
	Message string       // Message is a human-readable description of the resource's state.
}

// Ledger provides the methods for managing Cluster records.
type Ledger interface {
	// Create creates a new Cluster.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// GetByID retrieves a Cluster by its ID.
	GetByID(context.Context, string) (*GetResponse, error)
	// GetByName retrieves a Cluster by its name.
	GetByName(context.Context, string) (*GetResponse, error)
	// UpdateStatus updates the state and message of an existing Cluster.
	UpdateStatus(context.Context, *UpdateStateRequest) (*UpdateResponse, error)
	// List returns a list of Cluster that match the provided filters.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Delete deletes a Cluster.
	Delete(context.Context, *DeleteRequest) error
}

// CreateRequest represents the Cluster creation request.
type CreateRequest struct {
	Name string
}

// CreateResponse represents the response after creating a new Cluster.
type CreateResponse struct {
	Record ClusterRecord
}

// UpdateStateRequest represents the request to update the state and message of a Cluster.
type UpdateStateRequest struct {
	Metadata core.Metadata
	Status   ClusterStatus
}

// GetResponse represents the response for fetching a Cluster.
type GetResponse struct {
	Record ClusterRecord
}

// UpdateResponse represents the response after updating the state of a Cluster.
type UpdateResponse struct {
	Record ClusterRecord
}

// ListRequest represents the request to list Clusters with filters.
type ListRequest struct {
	Filters ClusterListFilters
}

// ClusterFilters contains filters for querying the Cluster table.
type ClusterListFilters struct {
	IDIn       []string // IN condition
	NameIn     []string // IN condition
	VersionGte *uint64  // Greater than or equal condition
	VersionLte *uint64  // Less than or equal condition
	VersionEq  *uint64  // Equal condition

	IncludeDeleted bool   // IncludeDeleted indicates whether to include soft-deleted records in the result.
	Limit          uint32 // Limit is the maximum number of records to return.

	StateIn    []ClusterState
	StateNotIn []ClusterState
}

// ListResponse represents the response to a list request.
type ListResponse struct {
	Records []ClusterRecord
}

type DeleteRequest struct {
	Metadata core.Metadata
}
