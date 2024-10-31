package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var deploymentPlanMatchingCapabilityTableMigrations = []simplesql.Migration{
	{
		Version: 13, // Update the version number sequentially.
		Up: `
			CREATE TABLE deployment_plan_matching_capability (
				deployment_plan_id VARCHAR(255) NOT NULL,
				capability_type VARCHAR(255) NOT NULL,
				comparator VARCHAR(255) NOT NULL,
				capability_names TEXT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				PRIMARY KEY (deployment_plan_id, capability_type),
				FOREIGN KEY (deployment_plan_id) REFERENCES deployment_plan(id) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS deployment_plan_matching_capability;
			`,
	},
}

type DeploymentPlanMatchingCapabilityRow struct {
	DeploymentPlanID string `db:"deployment_plan_id" orm:"op=create filter=In"`
	CapabilityType   string `db:"capability_type" orm:"op=create"`
	Comparator       string `db:"comparator" orm:"op=create"`
	CapabilityNames  string `db:"capability_names" orm:"op=create"`
}

type DeploymentPlanMatchingCapabilityTableSelectFilters struct {
	DeploymentPlanIDIn []string `db:"deployment_plan_id:in"` // IN condition
}

const deploymentPlanMatchingCapabilityTableName = "deployment_plan_matching_capability"

type DeploymentMatchingCapabilityTable struct {
	simplesql.Database
	tableName string
}

func NewDeploymentMatchingCapabilityTable(db simplesql.Database) *DeploymentMatchingCapabilityTable {
	return &DeploymentMatchingCapabilityTable{
		Database:  db,
		tableName: deploymentPlanMatchingCapabilityTableName,
	}
}

func (s *DeploymentMatchingCapabilityTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row DeploymentPlanMatchingCapabilityRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *DeploymentMatchingCapabilityTable) List(ctx context.Context, filters DeploymentPlanMatchingCapabilityTableSelectFilters) ([]DeploymentPlanMatchingCapabilityRow, error) {
	var rows []DeploymentPlanMatchingCapabilityRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
