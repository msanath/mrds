package metainstance

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/internal/ledger/core"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"

	"github.com/google/uuid"
)

// ledger implements the Ledger interface.
type ledger struct {
	metaInstanceRepo Repository
}

// Repository provides the methods that the storage layer must implement to support the ledger.
type Repository interface {
	Insert(context.Context, MetaInstanceRecord) error
	GetByMetadata(context.Context, core.Metadata) (MetaInstanceRecord, error)
	GetByName(context.Context, string) (MetaInstanceRecord, error)
	UpdateState(context.Context, core.Metadata, MetaInstanceStatus) error
	UpdateDeploymentID(context.Context, core.Metadata, string) error
	Delete(context.Context, core.Metadata) error
	List(context.Context, MetaInstanceListFilters) ([]MetaInstanceRecord, error)

	InsertOperation(ctx context.Context, metadata core.Metadata, operation Operation) error
	UpdateOperationStatus(ctx context.Context, metadata core.Metadata, operationID string, status OperationStatus) error
	DeleteOperation(ctx context.Context, metadata core.Metadata, operationID string) error

	InsertRuntimeInstance(ctx context.Context, metadata core.Metadata, instance RuntimeInstance) error
	UpdateRuntimeInstanceStatus(ctx context.Context, metadata core.Metadata, instanceID string, status RuntimeInstanceStatus) error
	DeleteRuntimeInstance(ctx context.Context, metadata core.Metadata, instanceID string) error
}

// NewLedger creates a new Ledger instance.
func NewLedger(metaInstanceRepo Repository) Ledger {
	return &ledger{
		metaInstanceRepo: metaInstanceRepo,
	}
}

// Create creates a new MetaInstance.
func (l *ledger) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	// validate the request
	if req.Name == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"MetaInstance name is required",
		)
	}
	if req.DeploymentPlanID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"DeploymentPlanID is required",
		)
	}
	if req.DeploymentID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"DeploymentID is required",
		)
	}

	rec := MetaInstanceRecord{
		Metadata: core.Metadata{
			ID:      uuid.New().String(),
			Version: 0,
		},
		Name: req.Name,
		Status: MetaInstanceStatus{
			State:   MetaInstanceStatePendingAllocation,
			Message: "",
		},
		DeploymentPlanID: req.DeploymentPlanID,
		DeploymentID:     req.DeploymentID,
	}

	// TODO: Check instance count and create runtime instances

	err := l.metaInstanceRepo.Insert(ctx, rec)
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		Record: rec,
	}, nil
}

// GetByMetadata retrieves a MetaInstance by its metadata.
func (l *ledger) GetByMetadata(ctx context.Context, metadata *core.Metadata) (*GetResponse, error) {
	// validate the request
	if metadata == nil || metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to fetch by metadata",
		)
	}

	record, err := l.metaInstanceRepo.GetByMetadata(ctx, *metadata)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Record: record,
	}, nil
}

// GetByName retrieves a MetaInstance by its name.
func (l *ledger) GetByName(ctx context.Context, name string) (*GetResponse, error) {
	// validate the request
	if name == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Name missing. Name is required to fetch by name",
		)
	}

	record, err := l.metaInstanceRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Record: record,
	}, nil
}

var validStateTransitions = map[MetaInstanceState][]MetaInstanceState{
	MetaInstanceStatePendingAllocation: {MetaInstanceStateRunning, MetaInstanceStateTerminated},
	MetaInstanceStateRunning:           {MetaInstanceStateTerminated},
	MetaInstanceStateTerminated:        {MetaInstanceStateRunning},
}

// UpdateStatus updates the state and message of an existing MetaInstance.
func (l *ledger) UpdateStatus(ctx context.Context, req *UpdateStatusRequest) (*UpdateResponse, error) {
	// validate the request
	if req.Metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to update state",
		)
	}

	existingRecord, err := l.metaInstanceRepo.GetByMetadata(ctx, req.Metadata)
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

	err = l.metaInstanceRepo.UpdateState(ctx, req.Metadata, req.Status)
	if err != nil {
		return nil, err
	}

	record, err := l.metaInstanceRepo.GetByMetadata(ctx, core.Metadata{
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

func (l *ledger) UpdateDeploymentID(ctx context.Context, req *UpdateDeploymentIDRequest) (*UpdateResponse, error) {
	if req.DeploymentID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"DeploymentID is required",
		)
	}

	err := l.metaInstanceRepo.UpdateDeploymentID(ctx, req.Metadata, req.DeploymentID)
	if err != nil {
		return nil, err
	}

	record, err := l.metaInstanceRepo.GetByMetadata(ctx, core.Metadata{
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

// List returns a list of MetaInstances that match the provided filters.
func (l *ledger) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	records, err := l.metaInstanceRepo.List(ctx, req.Filters)
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Records: records,
	}, nil
}

func (l *ledger) Delete(ctx context.Context, req *DeleteRequest) error {
	return l.metaInstanceRepo.Delete(ctx, req.Metadata)
}

// AddRuntimeInstance adds a runtime instance to the MetaInstance.
func (l *ledger) AddRuntimeInstance(ctx context.Context, req *AddRuntimeInstanceRequest) (*UpdateResponse, error) {
	err := l.metaInstanceRepo.InsertRuntimeInstance(ctx, req.Metadata, req.RuntimeInstance)
	if err != nil {
		return nil, err
	}

	record, err := l.metaInstanceRepo.GetByMetadata(ctx, core.Metadata{
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

// UpdateRuntimeStatus updates the state and message of a runtime instance in the MetaInstance.
func (l *ledger) UpdateRuntimeStatus(ctx context.Context, req *UpdateRuntimeStatusRequest) (*UpdateResponse, error) {
	err := l.metaInstanceRepo.UpdateRuntimeInstanceStatus(ctx, req.Metadata, req.RuntimeInstanceID, req.Status)
	if err != nil {
		return nil, err
	}

	record, err := l.metaInstanceRepo.GetByMetadata(ctx, core.Metadata{
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

// RemoveRuntimeInstance removes a runtime instance from the MetaInstance.
func (l *ledger) RemoveRuntimeInstance(ctx context.Context, req *RemoveRuntimeInstanceRequest) (*UpdateResponse, error) {
	err := l.metaInstanceRepo.DeleteRuntimeInstance(ctx, req.Metadata, req.RuntimeInstanceID)
	if err != nil {
		return nil, err
	}

	record, err := l.metaInstanceRepo.GetByMetadata(ctx, core.Metadata{
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

// AddOperation adds an operation to the MetaInstance.
func (l *ledger) AddOperation(ctx context.Context, req *AddOperationRequest) (*UpdateResponse, error) {
	err := l.metaInstanceRepo.InsertOperation(ctx, req.Metadata, req.Operation)
	if err != nil {
		return nil, err
	}

	record, err := l.metaInstanceRepo.GetByMetadata(ctx, core.Metadata{
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

// UpdateOperationStatus updates the state and message of an operation in the MetaInstance.
func (l *ledger) UpdateOperationStatus(ctx context.Context, req *UpdateOperationStatusRequest) (*UpdateResponse, error) {
	err := l.metaInstanceRepo.UpdateOperationStatus(ctx, req.Metadata, req.OperationID, req.Status)
	if err != nil {
		return nil, err
	}

	record, err := l.metaInstanceRepo.GetByMetadata(ctx, core.Metadata{
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

// RemoveOperation removes an operation from the MetaInstance.
func (l *ledger) RemoveOperation(ctx context.Context, req *RemoveOperationRequest) (*UpdateResponse, error) {
	err := l.metaInstanceRepo.DeleteOperation(ctx, req.Metadata, req.OperationID)
	if err != nil {
		return nil, err
	}

	record, err := l.metaInstanceRepo.GetByMetadata(ctx, core.Metadata{
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
