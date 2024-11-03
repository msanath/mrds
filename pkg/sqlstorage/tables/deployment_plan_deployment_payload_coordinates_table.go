package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var deploymentPlanDeploymentPayloadCoordinatesTableMigrations = []simplesql.Migration{
	{
		Version: 12, // Update the version number sequentially.
		Up: `
			CREATE TABLE deployment_plan_deployment_payload_coordinates (
				deployment_plan_id VARCHAR(255) NOT NULL,
				deployment_id VARCHAR(255) NOT NULL,
				payload_name VARCHAR(255) NOT NULL,
				coordinates TEXT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				FOREIGN KEY (deployment_plan_id, payload_name) REFERENCES deployment_plan_application(deployment_plan_id, payload_name),
				FOREIGN KEY (deployment_plan_id) REFERENCES deployment_plan(id) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS deployment_plan_deployment_payload_coordinates;
			`,
	},
}

type DeploymentPlanDeploymentPayloadCoordinatesRow struct {
	DeploymentPlanID string `db:"deployment_plan_id" orm:"op=create filter=In"`
	DeploymentID     string `db:"deployment_id" orm:"op=create filter=In"`
	PayloadName      string `db:"payload_name" orm:"op=create"`
	Coordinates      string `db:"coordinates" orm:"op=create"`
}

type DeploymentPlanDeploymentPayloadCoordinatesTableSelectFilters struct {
	DeploymentPlanIDIn []string `db:"deployment_plan_id:in"` // IN condition
	DeploymentIDIn     []string `db:"deployment_id:in"`      // IN condition
}

const deploymentPlanDeploymentPayloadCoordinatesTableName = "deployment_plan_deployment_payload_coordinates"

type DeploymentPlanDeploymentPayloadCoordinatesTable struct {
	simplesql.Database
	tableName string
}

func NewDeploymentPlanDeploymentPayloadCoordinatesTable(db simplesql.Database) *DeploymentPlanDeploymentPayloadCoordinatesTable {
	return &DeploymentPlanDeploymentPayloadCoordinatesTable{
		Database:  db,
		tableName: deploymentPlanDeploymentPayloadCoordinatesTableName,
	}
}

func (s *DeploymentPlanDeploymentPayloadCoordinatesTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row DeploymentPlanDeploymentPayloadCoordinatesRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *DeploymentPlanDeploymentPayloadCoordinatesTable) List(ctx context.Context, filters DeploymentPlanDeploymentPayloadCoordinatesTableSelectFilters) ([]DeploymentPlanDeploymentPayloadCoordinatesRow, error) {
	var rows []DeploymentPlanDeploymentPayloadCoordinatesRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
