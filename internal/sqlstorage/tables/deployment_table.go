package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

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

type DeploymentRow struct {
	ID        string `db:"id" orm:"op=create key=primary_key filter=In"`
	Version   uint64 `db:"version" orm:"op=create,update"`
	Name      string `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"`
	IsDeleted bool   `db:"is_deleted"`
	State     string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message   string `db:"message" orm:"op=create,update"`
}

type DeploymentTableUpdateFields struct {
	State   *string `db:"state"`
	Message *string `db:"message"`
}

type DeploymentTableSelectFilters struct {
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

type DeploymentTable struct {
	simplesql.Database
	tableName string
}

func NewDeploymentTable(db simplesql.Database) *DeploymentTable {
	return &DeploymentTable{
		Database:  db,
		tableName: deploymentTableName,
	}
}

func (s *DeploymentTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row DeploymentRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *DeploymentTable) GetByIDAndVersion(ctx context.Context, id string, version uint64, isDeleted bool) (DeploymentRow, error) {
	var row DeploymentRow
	err := s.Database.GetRowByID(ctx, id, version, isDeleted, s.tableName, &row)
	if err != nil {
		return DeploymentRow{}, err
	}
	return row, nil
}

func (s *DeploymentTable) GetByName(ctx context.Context, name string) (DeploymentRow, error) {
	var row DeploymentRow
	err := s.Database.GetRowByName(ctx, name, s.tableName, &row)
	if err != nil {
		return DeploymentRow{}, err
	}
	return row, nil
}

func (s *DeploymentTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, id string, version uint64, updateFields DeploymentTableUpdateFields,
) error {
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, updateFields)
}

func (s *DeploymentTable) Delete(ctx context.Context, execer sqlx.ExecerContext, id string, version uint64) error {
	return s.Database.MarkRowAsDeleted(ctx, execer, id, version, s.tableName)
}

func (s *DeploymentTable) List(ctx context.Context, filters DeploymentTableSelectFilters) ([]DeploymentRow, error) {
	var rows []DeploymentRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
