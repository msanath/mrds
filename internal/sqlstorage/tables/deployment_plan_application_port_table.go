package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var deploymentPlanApplicationPortTableMigrations = []simplesql.Migration{
	{
		Version: 9, // Update the version number sequentially.
		Up: `
			CREATE TABLE deployment_plan_application_port (
				deployment_plan_id VARCHAR(255) NOT NULL,
				payload_name VARCHAR(255) NOT NULL,
				protocol VARCHAR(255) NOT NULL,
				port INT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				FOREIGN KEY (deployment_plan_id, payload_name) REFERENCES deployment_plan_application(deployment_plan_id, payload_name) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS deployment_plan_application_port;
			`,
	},
}

type DeploymentPlanApplicationPortRow struct {
	DeploymentPlanID string `db:"deployment_plan_id" orm:"op=create filter=In"`
	PayloadName      string `db:"payload_name" orm:"op=create filter=In"`
	Protocol         string `db:"protocol" orm:"op=create filter=In"`
	Port             uint32 `db:"port" orm:"op=create"`
}

type DeploymentPlanApplicationPortTableSelectFilters struct {
	DeploymentPlanIDIn []string `db:"deployment_plan_id:in"` // IN condition
	PayloadNameIn      []string `db:"payload_name:in"`       // IN condition
}

const deploymentPlanApplicationPortTableName = "deployment_plan_application_port"

type DeploymentPlanApplicationPortTable struct {
	simplesql.Database
	tableName string
}

func NewDeploymentPlanApplicationPortTable(db simplesql.Database) *DeploymentPlanApplicationPortTable {
	return &DeploymentPlanApplicationPortTable{
		Database:  db,
		tableName: deploymentPlanApplicationPortTableName,
	}
}

func (s *DeploymentPlanApplicationPortTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row DeploymentPlanApplicationPortRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *DeploymentPlanApplicationPortTable) List(ctx context.Context, filters DeploymentPlanApplicationPortTableSelectFilters) ([]DeploymentPlanApplicationPortRow, error) {
	var rows []DeploymentPlanApplicationPortRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
