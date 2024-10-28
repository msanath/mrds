package sqlstorage

import (
	"context"

	"github.com/msanath/gondolf/pkg/simplesql"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/deployment"
	"github.com/msanath/mrds/internal/sqlstorage/tables"
)

// deploymentStorage is a concrete implementation of DeploymentRepository using sqlx
type deploymentStorage struct {
	simplesql.Database
	deploymentTable *tables.DeploymentTable
}

// newDeploymentStorage creates a new storage instance satisfying the DeploymentRepository interface
func newDeploymentStorage(db simplesql.Database) deployment.Repository {
	return &deploymentStorage{
		Database:        db,
		deploymentTable: tables.NewDeploymentTable(db),
	}
}

func deploymentModelToRow(model deployment.DeploymentRecord) tables.DeploymentRow {
	return tables.DeploymentRow{
		ID:      model.Metadata.ID,
		Version: model.Metadata.Version,
		Name:    model.Name,
		State:   model.Status.State.ToString(),
		Message: model.Status.Message,
	}
}

func deploymentRowToModel(row tables.DeploymentRow) deployment.DeploymentRecord {
	return deployment.DeploymentRecord{
		Metadata: core.Metadata{
			ID:      row.ID,
			Version: row.Version,
		},
		Name: row.Name,
		Status: deployment.DeploymentStatus{
			State:   deployment.DeploymentStateFromString(row.State),
			Message: row.Message,
		},
	}
}

func (s *deploymentStorage) Insert(ctx context.Context, record deployment.DeploymentRecord) error {
	execer := s.DB
	err := s.deploymentTable.Insert(ctx, execer, deploymentModelToRow(record))
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *deploymentStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (deployment.DeploymentRecord, error) {
	row, err := s.deploymentTable.GetByIDAndVersion(ctx, metadata.ID, metadata.Version, metadata.IsDeleted)
	if err != nil {
		return deployment.DeploymentRecord{}, errHandler(err)
	}
	return deploymentRowToModel(row), nil
}

func (s *deploymentStorage) GetByName(ctx context.Context, name string) (deployment.DeploymentRecord, error) {
	row, err := s.deploymentTable.GetByName(ctx, name)
	if err != nil {
		return deployment.DeploymentRecord{}, errHandler(err)
	}
	return deploymentRowToModel(row), nil
}

func (s *deploymentStorage) UpdateState(ctx context.Context, metadata core.Metadata, status deployment.DeploymentStatus) error {
	execer := s.DB
	state := status.State.ToString()
	message := status.Message
	updateFields := tables.DeploymentTableUpdateFields{
		State:   &state,
		Message: &message,
	}
	return s.deploymentTable.Update(ctx, execer, metadata.ID, metadata.Version, updateFields)
}

func (s *deploymentStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	execer := s.DB
	return s.deploymentTable.Delete(ctx, execer, metadata.ID, metadata.Version)
}

func (s *deploymentStorage) List(ctx context.Context, filters deployment.DeploymentListFilters) ([]deployment.DeploymentRecord, error) {
	dbFilters := tables.DeploymentTableSelectFilters{
		IDIn:           append([]string{}, filters.IDIn...),
		NameIn:         append([]string{}, filters.NameIn...),
		VersionGte:     filters.VersionGte,
		VersionLte:     filters.VersionLte,
		VersionEq:      filters.VersionEq,
		IncludeDeleted: filters.IncludeDeleted,
		Limit:          filters.Limit,
	}
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, state.ToString())
	}
	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, state.ToString())
	}

	rows, err := s.deploymentTable.List(ctx, dbFilters)
	if err != nil {
		return nil, err
	}
	var records []deployment.DeploymentRecord
	for _, row := range rows {
		records = append(records, deploymentRowToModel(row))
	}
	return records, nil
}
