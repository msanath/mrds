package tables

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var deploymentPlanTableMigrations = []simplesql.Migration{
	{
		Version: 7, // Update the version number sequentially.
		Up: `
			CREATE TABLE deployment_plan (
				id VARCHAR(255) NOT NULL PRIMARY KEY,
				version BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				namespace VARCHAR(255) NOT NULL,
				service_name VARCHAR(255) NOT NULL,
				UNIQUE (name, deleted_at)
			);
		`,
		Down: `
				DROP TABLE IF EXISTS deployment_plan;
			`,
	},
}

type DeploymentPlanRow struct {
	ID          string `db:"id" orm:"op=create key=primary_key filter=In"`
	Version     uint64 `db:"version" orm:"op=create,update"`
	Name        string `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"`
	DeletedAt   int64  `db:"deleted_at"`
	State       string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message     string `db:"message" orm:"op=create,update"`
	Namespace   string `db:"namespace" orm:"op=create filter=In"`
	ServiceName string `db:"service_name" orm:"op=create filter=In"`
}

type DeploymentPlanKeys struct {
	ID   *string `db:"id"`
	Name *string `db:"name"`
}

type DeploymentPlanTableUpdateFields struct {
	State     *string `db:"state"`
	Message   *string `db:"message"`
	DeletedAt *int64  `db:"deleted_at"`
}

type DeploymentPlanTableSelectFilters struct {
	IDIn          []string `db:"id:in"`           // IN condition
	NameIn        []string `db:"name:in"`         // IN condition
	StateIn       []string `db:"state:in"`        // IN condition
	StateNotIn    []string `db:"state:not_in"`    // NOT IN condition
	VersionGte    *uint64  `db:"version:gte"`     // Greater than or equal condition
	VersionLte    *uint64  `db:"version:lte"`     // Less than or equal condition
	VersionEq     *uint64  `db:"version:eq"`      // Equal condition
	ServiceNameIn []string `db:"service_name:in"` // IN condition
	NamespaceIn   []string `db:"namespace:in"`    // IN condition

	IncludeDeleted bool   `db:"include_deleted"` // Special boolean handling
	Limit          uint32 `db:"limit"`
}

const deploymentPlanTableName = "deployment_plan"

type DeploymentPlanTable struct {
	simplesql.Database
	tableName string
}

func NewDeploymentPlanTable(db simplesql.Database) *DeploymentPlanTable {
	return &DeploymentPlanTable{
		Database:  db,
		tableName: deploymentPlanTableName,
	}
}

func (s *DeploymentPlanTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row DeploymentPlanRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *DeploymentPlanTable) Get(ctx context.Context, keys DeploymentPlanKeys) (DeploymentPlanRow, error) {
	var row DeploymentPlanRow
	err := s.Database.GetRowByKey(ctx, s.tableName, keys, &row)
	if err != nil {
		return DeploymentPlanRow{}, err
	}
	return row, nil
}

func (s *DeploymentPlanTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, id string, version uint64, updateFields DeploymentPlanTableUpdateFields,
) error {
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, updateFields)
}

func (s *DeploymentPlanTable) Delete(ctx context.Context, execer sqlx.ExecerContext, id string, version uint64) error {
	timeNow := time.Now().Unix()
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, DeploymentPlanTableUpdateFields{
		DeletedAt: &timeNow,
	})
}

func (s *DeploymentPlanTable) List(ctx context.Context, filters DeploymentPlanTableSelectFilters) ([]DeploymentPlanRow, error) {
	var rows []DeploymentPlanRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
