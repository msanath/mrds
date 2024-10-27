package sqlstorage

import (
	"context"

	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/deployment"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var deploymentTableMigrations = []simplesql.Migration{
	{
		Version: 4, // Update the version number sequentially.
		Up: `
			CREATE TABLE deployment (
				id VARCHAR(255) NOT NULL PRIMARY KEY,
				version BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				UNIQUE (id, name, is_deleted)
			);
		`,
		Down: `
				DROP TABLE IF EXISTS deployment;
			`,
	},
}

type deploymentRow struct {
	ID        string `db:"id" orm:"op=create key=primary_key filter=In"`
	Version   uint64 `db:"version" orm:"op=create,update"`
	Name      string `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"`
	IsDeleted bool   `db:"is_deleted"`
	State     string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message   string `db:"message" orm:"op=create,update"`
}

type deploymentUpdateFields struct {
	State   *string `db:"state"`
	Message *string `db:"message"`
}

type deploymentSelectFilters struct {
	IDIn       []string `db:"id:in"`        // IN condition
	NameIn     []string `db:"name:in"`      // IN condition
	StateIn    []string `db:"state:in"`     // IN condition
	StateNotIn []string `db:"state:not_in"` // NOT IN condition
	VersionGte *uint64  `db:"version:gte"`  // Greater than or equal condition
	VersionLte *uint64  `db:"version:lte"`  // Less than or equal condition
	VersionEq  *uint64  `db:"version:eq"`   // Equal condition

	IncludeDeleted bool   `db:"include_deleted"` // Special boolean handling
	Limit          uint32 `db:"limit"`
}

const deploymentTableName = "deployment"

func deploymentModelToRow(model deployment.DeploymentRecord) deploymentRow {
	return deploymentRow{
		ID:      model.Metadata.ID,
		Version: model.Metadata.Version,
		Name:    model.Name,
		State:   model.Status.State.ToString(),
		Message: model.Status.Message,
	}
}

func deploymentRowToModel(row deploymentRow) deployment.DeploymentRecord {
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

// deploymentStorage is a concrete implementation of DeploymentRepository using sqlx
type deploymentStorage struct {
	simplesql.Database
	tableName  string
	modelToRow func(deployment.DeploymentRecord) deploymentRow
	rowToModel func(deploymentRow) deployment.DeploymentRecord
}

// newDeploymentStorage creates a new storage instance satisfying the DeploymentRepository interface
func newDeploymentStorage(db simplesql.Database) deployment.Repository {
	return &deploymentStorage{
		Database:   db,
		tableName:  deploymentTableName,
		modelToRow: deploymentModelToRow,
		rowToModel: deploymentRowToModel,
	}
}

func (s *deploymentStorage) Insert(ctx context.Context, record deployment.DeploymentRecord) error {
	row := s.modelToRow(record)
	err := s.Database.InsertRow(ctx, s.DB, s.tableName, row)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *deploymentStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (deployment.DeploymentRecord, error) {
	var row deploymentRow
	err := s.Database.GetRowByID(ctx, metadata.ID, metadata.Version, metadata.IsDeleted, s.tableName, &row)
	if err != nil {
		return deployment.DeploymentRecord{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *deploymentStorage) GetByName(ctx context.Context, deploymentName string) (deployment.DeploymentRecord, error) {
	var row deploymentRow
	err := s.Database.GetRowByName(ctx, deploymentName, s.tableName, &row)
	if err != nil {
		return deployment.DeploymentRecord{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *deploymentStorage) UpdateState(ctx context.Context, metadata core.Metadata, status deployment.DeploymentStatus) error {
	state := status.State.ToString()
	err := s.Database.UpdateRow(ctx, s.DB, metadata.ID, metadata.Version, s.tableName, deploymentUpdateFields{
		State:   &state,
		Message: &status.Message,
	})
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *deploymentStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	err := s.Database.MarkRowAsDeleted(ctx, s.DB, metadata.ID, metadata.Version, s.tableName)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *deploymentStorage) List(ctx context.Context, filters deployment.DeploymentListFilters) ([]deployment.DeploymentRecord, error) {
	// Extract core filters
	dbFilters := deploymentSelectFilters{
		IDIn:           append([]string{}, filters.IDIn...),
		NameIn:         append([]string{}, filters.NameIn...),
		VersionGte:     filters.VersionGte,
		VersionLte:     filters.VersionLte,
		VersionEq:      filters.VersionEq,
		IncludeDeleted: filters.IncludeDeleted,
		Limit:          filters.Limit,
	}

	// Extract deployment specific filters
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, state.ToString())
	}
	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, state.ToString())
	}

	var rows []deploymentRow
	err := s.Database.SelectRows(ctx, s.tableName, dbFilters, &rows)
	if err != nil {
		return nil, err
	}

	var records []deployment.DeploymentRecord
	for _, row := range rows {
		records = append(records, s.rowToModel(row))
	}

	return records, nil
}
