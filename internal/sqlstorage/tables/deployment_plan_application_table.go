package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var deploymentPlanApplicationTableMigrations = []simplesql.Migration{
	{
		Version: 8, // Update the version number sequentially.
		Up: `
			CREATE TABLE deployment_plan_application (
				deployment_plan_id VARCHAR(255) NOT NULL,
				payload_name VARCHAR(255) NOT NULL,
				cores INT NOT NULL,
				memory INT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				PRIMARY KEY (deployment_plan_id, payload_name),
				FOREIGN KEY (deployment_plan_id) REFERENCES deployment_plan(id) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS deployment_plan_application;
			`,
	},
}

type DeploymentPlanApplicationRow struct {
	DeploymentPlanID string `db:"deployment_plan_id" orm:"op=create filter=In"`
	PayloadName      string `db:"payload_name" orm:"op=create filter=In"`
	Cores            uint32 `db:"cores" orm:"op=create"`
	Memory           uint32 `db:"memory" orm:"op=create"`
}

type DeploymentPlanApplicationTableSelectFilters struct {
	DeploymentPlanIDIn []string `db:"deployment_plan_id:in"` // IN condition
	PayloadNameIn      []string `db:"payload_name:in"`       // IN condition
}

const deploymentPlanApplicationTableName = "deployment_plan_application"

type DeploymentPlanApplicationTable struct {
	simplesql.Database
	tableName string
}

func NewDeploymentPlanApplicationTable(db simplesql.Database) *DeploymentPlanApplicationTable {
	return &DeploymentPlanApplicationTable{
		Database:  db,
		tableName: deploymentPlanApplicationTableName,
	}
}

func (s *DeploymentPlanApplicationTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row DeploymentPlanApplicationRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *DeploymentPlanApplicationTable) List(ctx context.Context, filters DeploymentPlanApplicationTableSelectFilters) ([]DeploymentPlanApplicationRow, error) {
	var rows []DeploymentPlanApplicationRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
