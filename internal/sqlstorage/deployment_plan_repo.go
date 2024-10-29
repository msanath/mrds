
package sqlstorage

import (
	"context"

	"github.com/msanath/gondolf/pkg/simplesql"
	"github.com/msanath/mrds/internal/ledger/deploymentplan"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/sqlstorage/tables"
)

// deploymentPlanStorage is a concrete implementation of DeploymentPlanRepository using sqlx
type deploymentPlanStorage struct {
	simplesql.Database
	deploymentPlanTable *tables.DeploymentPlanTable
}

// newDeploymentPlanStorage creates a new storage instance satisfying the DeploymentPlanRepository interface
func newDeploymentPlanStorage(db simplesql.Database) deploymentplan.Repository {
	return &deploymentPlanStorage{
		Database:     db,
		deploymentPlanTable: tables.NewDeploymentPlanTable(db),
	}
}

func deploymentPlanModelToRow(model deploymentplan.DeploymentPlanRecord) tables.DeploymentPlanRow {
	return tables.DeploymentPlanRow{
		ID:      model.Metadata.ID,
		Version: model.Metadata.Version,
		Name:    model.Name,
		State:   model.Status.State.ToString(),
		Message: model.Status.Message,
	}
}

func deploymentPlanRowToModel(row tables.DeploymentPlanRow) deploymentplan.DeploymentPlanRecord {
	return deploymentplan.DeploymentPlanRecord{
		Metadata: core.Metadata{
			ID:      row.ID,
			Version: row.Version,
		},
		Name: row.Name,
		Status: deploymentplan.DeploymentPlanStatus{
			State:   deploymentplan.DeploymentPlanStateFromString(row.State),
			Message: row.Message,
		},
	}
}

func (s *deploymentPlanStorage) Insert(ctx context.Context, record deploymentplan.DeploymentPlanRecord) error {
	execer := s.DB
	err := s.deploymentPlanTable.Insert(ctx, execer, deploymentPlanModelToRow(record))
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *deploymentPlanStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (deploymentplan.DeploymentPlanRecord, error) {
	row, err := s.deploymentPlanTable.GetByIDAndVersion(ctx, metadata.ID, metadata.Version, metadata.IsDeleted)
	if err != nil {
		return deploymentplan.DeploymentPlanRecord{}, errHandler(err)
	}
	return deploymentPlanRowToModel(row), nil
}

func (s *deploymentPlanStorage) GetByName(ctx context.Context, name string) (deploymentplan.DeploymentPlanRecord, error) {
	row, err := s.deploymentPlanTable.GetByName(ctx, name)
	if err != nil {
		return deploymentplan.DeploymentPlanRecord{}, errHandler(err)
	}
	return deploymentPlanRowToModel(row), nil
}

func (s *deploymentPlanStorage) UpdateState(ctx context.Context, metadata core.Metadata, status deploymentplan.DeploymentPlanStatus) error {
	execer := s.DB
	state := status.State.ToString()
	message := status.Message
	updateFields := tables.DeploymentPlanTableUpdateFields{
		State:   &state,
		Message: &message,
	}
	return s.deploymentPlanTable.Update(ctx, execer, metadata.ID, metadata.Version, updateFields)
}

func (s *deploymentPlanStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	execer := s.DB
	return s.deploymentPlanTable.Delete(ctx, execer, metadata.ID, metadata.Version)
}

func (s *deploymentPlanStorage) List(ctx context.Context, filters deploymentplan.DeploymentPlanListFilters) ([]deploymentplan.DeploymentPlanRecord, error) {
	dbFilters := tables.DeploymentPlanTableSelectFilters{
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

	rows, err := s.deploymentPlanTable.List(ctx, dbFilters)
	if err != nil {
		return nil, err
	}
	var records []deploymentplan.DeploymentPlanRecord
	for _, row := range rows {
		records = append(records, deploymentPlanRowToModel(row))
	}
	return records, nil
}
