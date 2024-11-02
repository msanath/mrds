package deploymentplan

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
	Insert(context.Context, DeploymentPlanRecord) error
	GetByID(context.Context, string) (DeploymentPlanRecord, error)
	GetByName(context.Context, string) (DeploymentPlanRecord, error)
	UpdateStatus(context.Context, core.Metadata, DeploymentPlanStatus) error
	Delete(context.Context, core.Metadata) error
	List(context.Context, DeploymentPlanListFilters) ([]DeploymentPlanRecord, error)

	InsertDeployment(context.Context, core.Metadata, Deployment) error
	UpdateDeploymentStatus(ctx context.Context, metadata core.Metadata, deploymentID string, status DeploymentStatus) error
}

// NewLedger creates a new Ledger instance.
func NewLedger(repo Repository) Ledger {
	return &ledger{repo: repo}
}

// Create creates a new DeploymentPlan.
func (l *ledger) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	// validate the request
	if req.Name == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"DeploymentPlan name is required",
		)
	}
	if req.Namespace == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Namespace is required",
		)
	}
	if req.ServiceName == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ServiceName is required",
		)
	}
	if len(req.Applications) == 0 {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Applications are required",
		)
	}
	for _, app := range req.Applications {
		if app.PayloadName == "" {
			return nil, ledgererrors.NewLedgerError(
				ledgererrors.ErrRequestInvalid,
				"PayloadName is required",
			)
		}
	}

	rec := DeploymentPlanRecord{
		Metadata: core.Metadata{
			ID:      uuid.New().String(),
			Version: 0,
		},
		Name:                        req.Name,
		Namespace:                   req.Namespace,
		ServiceName:                 req.ServiceName,
		MatchingComputeCapabilities: req.MatchingComputeCapabilities,
		Applications:                req.Applications,
		Status: DeploymentPlanStatus{
			State:   DeploymentPlanStateActive,
			Message: "",
		},
	}

	err := l.repo.Insert(ctx, rec)
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		Record: rec,
	}, nil
}

// GetByID retrieves a DeploymentPlan by its ID.
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

// GetByName retrieves a DeploymentPlan by its name.
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

var validStateTransitions = map[DeploymentPlanState][]DeploymentPlanState{
	DeploymentPlanStateActive:   {DeploymentPlanStateInactive},
	DeploymentPlanStateInactive: {DeploymentPlanStateActive},
}

// UpdateStatus updates the state and message of an existing DeploymentPlan.
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

	err = l.repo.UpdateStatus(ctx, req.Metadata, req.Status)
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

// List returns a list of DeploymentPlans that match the provided filters.
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

func (l *ledger) AddDeployment(ctx context.Context, req *AddDeploymentRequest) (*UpdateResponse, error) {
	if req.DeploymentID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"DeploymentID is required",
		)
	}
	existingPlan, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		return nil, err
	}

	// If there is an existing deployment InProgress, then we cannot add a new deployment
	// ??? Why can't we add a new deployment when there is an existing deployment in progress?
	for _, deployment := range existingPlan.Deployments {
		if deployment.Status.State == DeploymentStateInProgress {
			return nil, ledgererrors.NewLedgerError(
				ledgererrors.ErrRequestInvalid,
				fmt.Sprintf("Cannot add deployment when there is an existing deployment in progress. InProgress DeploymentID: %s", deployment.ID),
			)
		}
	}

	// make a map of app to bool.
	// if all apps are present, then we can add the deployment
	// if any app is missing, then we cannot add the deployment
	appMap := make(map[string]bool)
	for _, app := range existingPlan.Applications {
		appMap[app.PayloadName] = false
	}
	for _, coordinates := range req.PayloadCoordinates {
		if coordinates.Coordinates == nil {
			return nil, ledgererrors.NewLedgerError(
				ledgererrors.ErrRequestInvalid,
				fmt.Sprintf("Coordinates missing for payload %s", coordinates.PayloadName),
			)
		}
		if _, ok := appMap[coordinates.PayloadName]; ok {
			appMap[coordinates.PayloadName] = true
		}
	}
	for appName, present := range appMap {
		if !present {
			return nil, ledgererrors.NewLedgerError(
				ledgererrors.ErrRequestInvalid,
				fmt.Sprintf("Payload coordinates missing for payload %s", appName),
			)
		}
	}

	err = l.repo.InsertDeployment(ctx, req.Metadata, Deployment{
		ID:                 req.DeploymentID,
		PayloadCoordinates: req.PayloadCoordinates,
		InstanceCount:      req.InstanceCount,
		Status: DeploymentStatus{
			State:   DeploymentStatePending,
			Message: "",
		},
	})
	if err != nil {
		return nil, err
	}

	// Get the record again to return the updated record
	record, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Record: record,
	}, nil
}

var validDeploymentStateTransitions = map[DeploymentState][]DeploymentState{
	DeploymentStatePending:    {DeploymentStateInProgress, DeploymentStateCancelled},
	DeploymentStateInProgress: {DeploymentStateCancelled, DeploymentStateFailed, DeploymentStatePaused, DeploymentStateCompleted},
	DeploymentStatePaused:     {DeploymentStateInProgress, DeploymentStateCancelled},
}

func (l *ledger) UpdateDeploymentStatus(ctx context.Context, req *UpdateDeploymentStatusRequest) (*UpdateResponse, error) {
	existingPlan, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		return nil, err
	}

	// validate the deployment
	var deployment *Deployment
	for i, d := range existingPlan.Deployments {
		if d.ID == req.DeploymentID {
			deployment = &existingPlan.Deployments[i]
			break
		}
	}
	if deployment == nil {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			fmt.Sprintf("Deployment with ID %s not found", req.DeploymentID),
		)
	}

	// validate the state transition
	validStates, ok := validDeploymentStateTransitions[deployment.Status.State]
	if !ok {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			fmt.Sprintf("Invalid state transition from %s to %s", deployment.Status.State, req.Status.State),
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
			fmt.Sprintf("Invalid state transition from %s to %s", deployment.Status.State, req.Status.State),
		)
	}

	err = l.repo.UpdateDeploymentStatus(ctx, req.Metadata, req.DeploymentID, req.Status)
	if err != nil {
		return nil, err
	}
	// Get the record again to return the updated record
	record, err := l.repo.GetByID(ctx, req.Metadata.ID)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Record: record,
	}, nil
}
