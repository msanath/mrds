package computecapability

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/internal/ledger/core"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"

	"github.com/google/uuid"
)

// ledger implements the Ledger interface.
type ledger struct {
	repo Repository
}

// Repository provides the methods that the storage layer must implement to support the ledger.
type Repository interface {
	Insert(context.Context, ComputeCapabilityRecord) error
	GetByMetadata(context.Context, core.Metadata) (ComputeCapabilityRecord, error)
	GetByName(context.Context, string) (ComputeCapabilityRecord, error)
	UpdateState(context.Context, core.Metadata, ComputeCapabilityStatus) error
	Delete(context.Context, core.Metadata) error
	List(context.Context, ComputeCapabilityListFilters) ([]ComputeCapabilityRecord, error)
}

// NewLedger creates a new Ledger instance.
func NewLedger(repo Repository) Ledger {
	return &ledger{repo: repo}
}

// Create creates a new ComputeCapability.
func (l *ledger) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	// validate the request
	if req.Name == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ComputeCapability name is required",
		)
	}

	rec := ComputeCapabilityRecord{
		Metadata: core.Metadata{
			ID:      uuid.New().String(),
			Version: 0,
		},
		Name: req.Name,
		Status: ComputeCapabilityStatus{
			State:   ComputeCapabilityStatePending,
			Message: "",
		},
		Type:  req.Type,
		Score: req.Score,
	}

	err := l.repo.Insert(ctx, rec)
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		Record: rec,
	}, nil
}

// GetByMetadata retrieves a ComputeCapability by its metadata.
func (l *ledger) GetByMetadata(ctx context.Context, metadata *core.Metadata) (*GetResponse, error) {
	// validate the request
	if metadata == nil || metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to fetch by metadata",
		)
	}

	record, err := l.repo.GetByMetadata(ctx, *metadata)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Record: record,
	}, nil
}

// GetByName retrieves a ComputeCapability by its name.
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

var validStateTransitions = map[ComputeCapabilityState][]ComputeCapabilityState{
	ComputeCapabilityStatePending:  {ComputeCapabilityStateActive, ComputeCapabilityStateInActive},
	ComputeCapabilityStateActive:   {ComputeCapabilityStateInActive},
	ComputeCapabilityStateInActive: {ComputeCapabilityStateActive},
}

// UpdateStatus updates the state and message of an existing ComputeCapability.
func (l *ledger) UpdateStatus(ctx context.Context, req *UpdateStateRequest) (*UpdateResponse, error) {
	// validate the request
	if req.Metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to update state",
		)
	}

	existingRecord, err := l.repo.GetByMetadata(ctx, req.Metadata)
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

	err = l.repo.UpdateState(ctx, req.Metadata, req.Status)
	if err != nil {
		return nil, err
	}

	record, err := l.repo.GetByMetadata(ctx, core.Metadata{
		ID:      req.Metadata.ID,
		Version: req.Metadata.Version + 1,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Record: record,
	}, nil
}

// List returns a list of ComputeCapabilitys that match the provided filters.
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
