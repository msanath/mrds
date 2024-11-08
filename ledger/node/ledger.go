package node

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/ledger/core"
	ledgererrors "github.com/msanath/mrds/ledger/errors"

	"github.com/google/uuid"
)

// ledger implements the Ledger interface.
type ledger struct {
	repo Repository
}

// Repository provides the methods that the storage layer must implement to support the ledger.
type Repository interface {
	Insert(context.Context, NodeRecord) error
	GetByID(context.Context, string) (NodeRecord, error)
	GetByName(context.Context, string) (NodeRecord, error)
	UpdateStatus(context.Context, core.Metadata, NodeStatus, string) error
	Delete(context.Context, core.Metadata) error
	List(context.Context, NodeListFilters) ([]NodeRecord, error)

	InsertDisruption(context.Context, core.Metadata, Disruption) error
	DeleteDisruption(ctx context.Context, metadata core.Metadata, disruptionID string) error
	UpdateDisruptionStatus(ctx context.Context, metadata core.Metadata, disruptionID string, status DisruptionStatus) error

	InsertCapability(ctx context.Context, metadata core.Metadata, capabilityID string) error
	DeleteCapability(ctx context.Context, metadata core.Metadata, capabilityID string) error
}

// NewLedger creates a new Ledger instance.
func NewLedger(repo Repository) Ledger {
	return &ledger{repo: repo}
}

// Create creates a new Node.
func (l *ledger) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	// validate the request
	if req.Name == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Node name is required",
		)
	}
	if req.TotalResources.Cores == 0 || req.TotalResources.Memory == 0 {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"TotalResources must have non-zero values for Cores and Memory",
		)
	}
	if req.SystemReservedResources.Cores > req.TotalResources.Cores || req.SystemReservedResources.Memory > req.TotalResources.Memory {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"SystemReservedResources cannot be greater than TotalResources",
		)
	}
	if req.UpdateDomain == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"UpdateDomain is required",
		)
	}

	rec := NodeRecord{
		Metadata: core.Metadata{
			ID:      uuid.New().String(),
			Version: 0,
		},
		Name: req.Name,
		Status: NodeStatus{
			State:   NodeStateUnallocated,
			Message: "",
		},
		UpdateDomain:            req.UpdateDomain,
		TotalResources:          req.TotalResources,
		SystemReservedResources: req.SystemReservedResources,
		RemainingResources: Resources{
			Cores:  req.TotalResources.Cores - req.SystemReservedResources.Cores,
			Memory: req.TotalResources.Memory - req.SystemReservedResources.Memory,
		},
		CapabilityIDs: req.CapabilityIDs,
		LocalVolumes:  req.LocalVolumes,
	}

	err := l.repo.Insert(ctx, rec)
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		Record: rec,
	}, nil
}

// GetByID retrieves a Node by its metadata.
func (l *ledger) GetByID(ctx context.Context, id string) (*GetResponse, error) {
	// validate the request
	if id == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to fetch by ID",
		)
	}

	record, err := l.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Record: record,
	}, nil
}

// GetByName retrieves a Node by its name.
func (l *ledger) GetByName(ctx context.Context, name string) (*GetResponse, error) {
	// validate the request
	if name == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Name missing. Name is required to fetch by name",
		)
	}

	record, err := l.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Record: record,
	}, nil
}

var validStateTransitions = map[NodeState][]NodeState{
	NodeStateUnallocated: {NodeStateAllocating},
	NodeStateAllocating:  {NodeStateAllocated, NodeStateEvicted},
	NodeStateAllocated:   {NodeStateEvicted},
	NodeStateEvicted:     {NodeStateSanitizing},
	NodeStateSanitizing:  {NodeStateEvicted, NodeStateUnallocated},
}

// UpdateStatus updates the state and message of an existing Node.
func (l *ledger) UpdateStatus(ctx context.Context, req *UpdateStatusRequest) (*UpdateResponse, error) {
	// validate the request
	if req.Metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to update state",
		)
	}

	existingRecord, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		if ledgererrors.IsLedgerError(err) && ledgererrors.AsLedgerError(err).Code == ledgererrors.ErrRecordNotFound {
			return nil, ledgererrors.NewLedgerError(
				ledgererrors.ErrRecordInsertConflict,
				"Either record does not exist or version mismatch resulted in conflict. Check and retry.",
			)
		}
	}

	// validate the state transition
	validStates, ok := validStateTransitions[existingRecord.Status.State]
	if !ok {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			fmt.Sprintf("Invalid state transition from %s to %s", existingRecord.Status.State, req.Status.State),
		)
	}
	var valid bool
	for _, state := range validStates {
		if state == req.Status.State {
			valid = true
			break
		}
	}
	if !valid {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			fmt.Sprintf("Invalid state transition from %s to %s", existingRecord.Status.State, req.Status.State),
		)
	}

	if req.Status.State == NodeStateAllocating || req.Status.State == NodeStateAllocated {
		if req.ClusterID == "" {
			return nil, ledgererrors.NewLedgerError(
				ledgererrors.ErrRequestInvalid,
				"ClusterID is required when transitioning to NodeStateAllocating or NodeStateAllocated",
			)
		}
		if req.Status.State == NodeStateAllocated {
			if req.ClusterID != existingRecord.ClusterID {
				return nil, ledgererrors.NewLedgerError(
					ledgererrors.ErrRequestInvalid,
					"ClusterID cannot be changed when transitioning to NodeStateAllocated",
				)
			}
		}
	}
	if req.Status.State != NodeStateAllocating && req.Status.State != NodeStateAllocated {
		if req.ClusterID != "" {
			return nil, ledgererrors.NewLedgerError(
				ledgererrors.ErrRequestInvalid,
				"ClusterID is only allowed when transitioning to NodeStateAllocating or NodeStateAllocated",
			)
		}
	}

	if req.ClusterID != "" && !(req.Status.State == NodeStateAllocating || req.Status.State == NodeStateAllocated) {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ClusterID is only allowed when transitioning to NodeStateAllocating",
		)
	}

	err = l.repo.UpdateStatus(ctx, req.Metadata, req.Status, req.ClusterID)
	if err != nil {
		return nil, err
	}

	record, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Record: record,
	}, nil
}

// List returns a list of Nodes that match the provided filters.
func (l *ledger) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	records, err := l.repo.List(ctx, req.Filters)
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Records: records,
	}, nil
}

func (l *ledger) Delete(ctx context.Context, req *DeleteRequest) error {
	return l.repo.Delete(ctx, req.Metadata)
}

// AddDisruption adds a disruption to a Node.
func (l *ledger) AddDisruption(ctx context.Context, req *AddDisruptionRequest) (*UpdateResponse, error) {
	// validate the request
	if req.Metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to add a disruption",
		)
	}
	if req.Disruption.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Disruption ID is required to add a disruption",
		)
	}
	if req.Disruption.Status.State != DisruptionStateScheduled {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Disruption must be in scheduled state to add",
		)
	}

	err := l.repo.InsertDisruption(ctx, req.Metadata, req.Disruption)
	if err != nil {
		return nil, err
	}

	record, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Record: record,
	}, nil
}

var validDisruptionStateTransitions = map[DisruptionState][]DisruptionState{
	DisruptionStateScheduled: {DisruptionStateApproved, DisruptionStateCompleted},
	DisruptionStateApproved:  {DisruptionStateCompleted},
}

// UpdateDisruptionStatus updates the status of a disruption on a Node.
func (l *ledger) UpdateDisruptionStatus(ctx context.Context, req *UpdateDisruptionStatusRequest) (*UpdateResponse, error) {
	// validate the request
	if req.Metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to update disruption status",
		)
	}
	if req.DisruptionID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Disruption ID is required to update disruption status",
		)
	}

	existingRecord, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		if ledgererrors.IsLedgerError(err) && ledgererrors.AsLedgerError(err).Code == ledgererrors.ErrRecordNotFound {
			return nil, ledgererrors.NewLedgerError(
				ledgererrors.ErrRecordInsertConflict,
				"Either record does not exist or version mismatch resulted in conflict. Check and retry.",
			)
		}
	}
	found := false
	for _, existingDisruption := range existingRecord.Disruptions {
		if existingDisruption.ID == req.DisruptionID {
			found = true
			validStates, ok := validDisruptionStateTransitions[existingDisruption.Status.State]
			if !ok {
				return nil, ledgererrors.NewLedgerError(
					ledgererrors.ErrRequestInvalid,
					fmt.Sprintf("Invalid state transition from %s to %s", existingRecord.Status.State, req.Status.State),
				)
			}
			var valid bool
			for _, state := range validStates {
				if state == req.Status.State {
					valid = true
					break
				}
			}
			if !valid {
				return nil, ledgererrors.NewLedgerError(
					ledgererrors.ErrRequestInvalid,
					fmt.Sprintf("Invalid state transition from %s to %s", existingRecord.Status.State, req.Status.State),
				)
			}
		}
	}
	if !found {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Disruption not found",
		)
	}

	err = l.repo.UpdateDisruptionStatus(ctx, req.Metadata, req.DisruptionID, req.Status)
	if err != nil {
		return nil, err
	}

	record, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Record: record,
	}, nil
}

// RemoveDisruption removes a disruption from a Node.
func (l *ledger) RemoveDisruption(ctx context.Context, req *RemoveDisruptionRequest) (*UpdateResponse, error) {
	// validate the request
	if req.Metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to remove a disruption",
		)
	}
	if req.DisruptionID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Disruption ID is required to remove a disruption",
		)
	}
	existingRecord, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		if ledgererrors.IsLedgerError(err) && ledgererrors.AsLedgerError(err).Code == ledgererrors.ErrRecordNotFound {
			return nil, ledgererrors.NewLedgerError(
				ledgererrors.ErrRecordInsertConflict,
				"Either record does not exist or version mismatch resulted in conflict. Check and retry.",
			)
		}
	}
	found := false
	for _, existingDisruption := range existingRecord.Disruptions {
		if existingDisruption.ID == req.DisruptionID {
			found = true
			if existingDisruption.Status.State != DisruptionStateCompleted {
				return nil, ledgererrors.NewLedgerError(
					ledgererrors.ErrRequestInvalid,
					"Disruption must be in completed state to remove",
				)
			}
		}
	}
	if !found {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Disruption not found",
		)
	}

	err = l.repo.DeleteDisruption(ctx, req.Metadata, req.DisruptionID)
	if err != nil {
		return nil, err
	}

	record, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Record: record,
	}, nil
}

// AddCapability adds a capability to a Node.
func (l *ledger) AddCapability(ctx context.Context, req *AddCapabilityRequest) (*UpdateResponse, error) {
	// validate the request
	if req.Metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to add a capability",
		)
	}
	if req.CapabilityID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Capability ID is required to add a capability",
		)
	}

	err := l.repo.InsertCapability(ctx, req.Metadata, req.CapabilityID)
	if err != nil {
		return nil, err
	}

	record, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Record: record,
	}, nil
}

// RemoveCapability removes a capability from a Node.
func (l *ledger) RemoveCapability(ctx context.Context, req *RemoveCapabilityRequest) (*UpdateResponse, error) {
	// validate the request
	if req.Metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to remove a capability",
		)
	}
	if req.CapabilityID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Capability ID is required to remove a capability",
		)
	}

	err := l.repo.DeleteCapability(ctx, req.Metadata, req.CapabilityID)
	if err != nil {
		return nil, err
	}

	record, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Record: record,
	}, nil
}
