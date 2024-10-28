package node

import (
	"context"
	"time"

	"github.com/msanath/mrds/internal/ledger/core"
)

// Node is a representation of the Node of an application.
type NodeRecord struct {
	Metadata core.Metadata // Metadata is the metadata that identifies the Node. It is a combination of the Node's name and version.
	Name     string        // Name is the name of the Node.
	Status   NodeStatus    // Status is the status of the Node.

	ClusterID               string    // ClusterID is the ID of the Cluster to which the Node belongs.
	UpdateDomain            string    // UpdateDomain is the update domain of the Node.
	TotalResources          Resources // TotalResources is the total resources available on the Node.
	SystemReservedResources Resources // SystemReservedResources is the resources reserved for system use.
	RemainingResources      Resources // RemainingResources is the resources available for application use.
}

type Resources struct {
	Cores  uint32
	Memory uint32
}

// NodeState is the state of a Node.
type NodeState string

const (
	NodeStateUnknown     NodeState = "NodeState_UNKNOWN"
	NodeStateUnallocated NodeState = "NodeState_UNALLOCATED"
	NodeStateAllocating  NodeState = "NodeState_ALLOCATING"
	NodeStateAllocated   NodeState = "NodeState_ALLOCATED"
	NodeStateEvicted     NodeState = "NodeState_EVICTED"
	NodeStateSanitizing  NodeState = "NodeState_SANITIZING"
)

// ToString returns the string representation of the NodeState.
func (s NodeState) ToString() string {
	return string(s)
}

// FromString converts a string into a NodeState if valid, otherwise returns an error.
func NodeStateFromString(s string) NodeState {
	switch s {
	case string(NodeStateUnallocated):
		return NodeStateUnallocated
	case string(NodeStateAllocated):
		return NodeStateAllocated
	case string(NodeStateEvicted):
		return NodeStateEvicted
	default:
		return NodeState(s) // Unknown state. Return as is.
	}
}

type NodeStatus struct {
	State   NodeState // State is the discrete condition of the resource.
	Message string    // Message is a human-readable description of the resource's state.
}

type NodeDisruption struct {
	ID        string
	EvictNode bool
	StartTime time.Time
	Status    DisruptionStatus
}

type DisruptionStatus struct {
	State   DisruptionState
	Message string
}

type DisruptionState string

const (
	DisruptionStateUnknown   DisruptionState = "DisruptionState_UNKNOWN"
	DisruptionStateScheduled DisruptionState = "DisruptionState_SCHEDULED"
	DisruptionStateApproved  DisruptionState = "DisruptionState_APPROVED"
	DisruptionStateCompleted DisruptionState = "DisruptionState_COMPLETED"
)

// Ledger provides the methods for managing Node records.
type Ledger interface {
	// Create creates a new Node.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// GetByMetadata retrieves a Node by its metadata.
	GetByMetadata(context.Context, *core.Metadata) (*GetResponse, error)
	// GetByName retrieves a Node by its name.
	GetByName(context.Context, string) (*GetResponse, error)
	// UpdateStatus updates the state and message of an existing Node.
	UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateResponse, error)
	// List returns a list of Node that match the provided filters.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Delete deletes a Node.
	Delete(context.Context, *DeleteRequest) error
}

// CreateRequest represents the Node creation request.
type CreateRequest struct {
	Name                    string
	UpdateDomain            string    // UpdateDomain is the update domain of the Node.
	TotalResources          Resources // TotalResources is the total resources available on the Node.
	SystemReservedResources Resources // SystemReservedResources is the resources reserved for system use.
}

// CreateResponse represents the response after creating a new Node.
type CreateResponse struct {
	Record NodeRecord
}

// UpdateStatusRequest represents the request to update the state and message of a Node.
type UpdateStatusRequest struct {
	Metadata  core.Metadata
	Status    NodeStatus
	ClusterID string
}

// GetResponse represents the response for fetching a Node.
type GetResponse struct {
	Record NodeRecord
}

// UpdateResponse represents the response after updating the state of a Node.
type UpdateResponse struct {
	Record NodeRecord
}

// ListRequest represents the request to list Nodes with filters.
type ListRequest struct {
	Filters NodeListFilters
}

// NodeFilters contains filters for querying the Node table.
type NodeListFilters struct {
	IDIn       []string // IN condition
	NameIn     []string // IN condition
	VersionGte *uint64  // Greater than or equal condition
	VersionLte *uint64  // Less than or equal condition
	VersionEq  *uint64  // Equal condition

	IncludeDeleted bool   // IncludeDeleted indicates whether to include soft-deleted records in the result.
	Limit          uint32 // Limit is the maximum number of records to return.

	StateIn            []NodeState
	StateNotIn         []NodeState
	RemainingCoresGte  *uint32
	RemainingCoresLte  *uint32
	RemainingMemoryGte *uint32
	RemainingMemoryLte *uint32
	ClusterIDIn        []string
	UpdateDomainIn     []string
}

// ListResponse represents the response to a list request.
type ListResponse struct {
	Records []NodeRecord
}

type DeleteRequest struct {
	Metadata core.Metadata
}
