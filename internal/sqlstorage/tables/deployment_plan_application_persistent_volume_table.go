package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var deploymentPlanApplicationPersistentVolumeTableMigrations = []simplesql.Migration{
	{
		Version: 11, // Update the version number sequentially.
		Up: `
			CREATE TABLE deployment_plan_application_persistent_volume (
				deployment_plan_id VARCHAR(255) NOT NULL,
				payload_name VARCHAR(255) NOT NULL,
				storage_class VARCHAR(255) NOT NULL,
				capacity INT NOT NULL,
				mount_path VARCHAR(255) NOT NULL,
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				PRIMARY KEY (deployment_plan_id, payload_name, storage_class, mount_path),
				FOREIGN KEY (deployment_plan_id, payload_name) REFERENCES deployment_plan_application(deployment_plan_id, payload_name) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS deployment_plan_application_persistent_volume;
			`,
	},
}

type DeploymentPlanApplicationPersistentVolumeRow struct {
	DeploymentPlanID string `db:"deployment_plan_id" orm:"op=create filter=In"`
	PayloadName      string `db:"payload_name" orm:"op=create filter=In"`
	StorageClass     string `db:"storage_class" orm:"op=create"`
	Capacity         uint32 `db:"capacity" orm:"op=create"`
	MountPath        string `db:"mount_path" orm:"op=create"`
}

type DeploymentPlanApplicationPersistentVolumeTableSelectFilters struct {
	DeploymentPlanIDIn []string `db:"deployment_plan_id:in"` // IN condition
	PayloadNameIn      []string `db:"payload_name:in"`       // IN condition
}

const deploymentPlanApplicationPersistentVolumeTableName = "deployment_plan_application_persistent_volume"

type DeploymentPlanApplicationPersistentVolumeTable struct {
	simplesql.Database
	tableName string
}

func NewDeploymentPlanApplicationPersistentVolumeTable(db simplesql.Database) *DeploymentPlanApplicationPersistentVolumeTable {
	return &DeploymentPlanApplicationPersistentVolumeTable{
		Database:  db,
		tableName: deploymentPlanApplicationPersistentVolumeTableName,
	}
}

func (s *DeploymentPlanApplicationPersistentVolumeTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row DeploymentPlanApplicationPersistentVolumeRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *DeploymentPlanApplicationPersistentVolumeTable) List(ctx context.Context, filters DeploymentPlanApplicationPersistentVolumeTableSelectFilters) ([]DeploymentPlanApplicationPersistentVolumeRow, error) {
	var rows []DeploymentPlanApplicationPersistentVolumeRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
