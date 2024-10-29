package tables

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var deploymentPlanDeploymentTableMigrations = []simplesql.Migration{
	{
		Version: 12, // Update the version number sequentially.
		Up: `
			CREATE TABLE deployment_plan_deployment (
				id VARCHAR(255) NOT NULL PRIMARY KEY,
				deployment_plan_id VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				instance_count INT NOT NULL,
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				FOREIGN KEY (deployment_plan_id) REFERENCES deployment_plan(id) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS deployment_plan_deployment;
			`,
	},
}

type DeploymentPlanDeploymentRow struct {
	ID               string `db:"id" orm:"op=create key=primary_key filter=In"`
	DeploymentPlanID string `db:"deployment_plan_id" orm:"op=create filter=In"`
	InstanceCount    uint32 `db:"instance_count" orm:"op=create"`
	State            string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message          string `db:"message" orm:"op=create,update"`
}

type DeploymentPlanDeploymentTableUpdateFields struct {
	State   *string `db:"state"`
	Message *string `db:"message"`
}

type DeploymentPlanDeploymentTableSelectFilters struct {
	IDIn               []string `db:"id:in"`                 // IN condition
	StateIn            []string `db:"state:in"`              // IN condition
	StateNotIn         []string `db:"state:not_in"`          // NOT IN condition
	DeploymentPlanIDIn []string `db:"deployment_plan_id:in"` // IN condition

}

const deploymentPlanDeploymentTableName = "deployment_plan_deployment"

type DeploymentPlanDeploymentTable struct {
	simplesql.Database
	tableName string
}

func NewDeploymentPlanDeploymentTable(db simplesql.Database) *DeploymentPlanDeploymentTable {
	return &DeploymentPlanDeploymentTable{
		Database:  db,
		tableName: deploymentPlanDeploymentTableName,
	}
}

func (s *DeploymentPlanDeploymentTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row DeploymentPlanDeploymentRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *DeploymentPlanDeploymentTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, deploymentID string, deploymentPlanID string, updateFields DeploymentPlanDeploymentTableUpdateFields,
) error {
	query := `
		UPDATE deployment_plan_deployment
		SET
	`
	params := map[string]interface{}{
		"id":                 deploymentID,
		"deployment_plan_id": deploymentPlanID,
	}

	var updates []string
	if updateFields.State != nil {
		updates = append(updates, "state = :state")
		params["state"] = *updateFields.State
	}

	if updateFields.Message != nil {
		updates = append(updates, "message = :message")
		params["message"] = *updateFields.Message
	}

	query += strings.Join(updates, ", ") + " WHERE id = :id AND deployment_plan_id = :deployment_plan_id"
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = s.DB.Rebind(query)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

func (s *DeploymentPlanDeploymentTable) List(ctx context.Context, filters DeploymentPlanDeploymentTableSelectFilters) ([]DeploymentPlanDeploymentRow, error) {
	var rows []DeploymentPlanDeploymentRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
