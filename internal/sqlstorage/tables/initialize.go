package tables

import "github.com/msanath/gondolf/pkg/simplesql"

func Initialize(simpleDB simplesql.Database) error {
	var schemaMigrations = []simplesql.Migration{}

	schemaMigrations = append(schemaMigrations, clusterTableMigrations...)
	schemaMigrations = append(schemaMigrations, computeCapabilityTableMigrations...)
	schemaMigrations = append(schemaMigrations, nodeTableMigrations...)
	schemaMigrations = append(schemaMigrations, metaInstanceTableMigrations...)
	schemaMigrations = append(schemaMigrations, nodeDisruptionTableMigrations...)
	schemaMigrations = append(schemaMigrations, nodeLocalVolumeTableMigrations...)
	schemaMigrations = append(schemaMigrations, nodeCapabilityTableMigrations...)
	schemaMigrations = append(schemaMigrations, deploymentPlanTableMigrations...)
	schemaMigrations = append(schemaMigrations, deploymentPlanApplicationTableMigrations...)
	schemaMigrations = append(schemaMigrations, deploymentPlanApplicationPortTableMigrations...)
	schemaMigrations = append(schemaMigrations, deploymentPlanApplicationPersistentVolumeTableMigrations...)
	schemaMigrations = append(schemaMigrations, deploymentPlanDeploymentTableMigrations...)
	schemaMigrations = append(schemaMigrations, deploymentPlanMatchingCapabilityTableMigrations...)
	schemaMigrations = append(schemaMigrations, deploymentPlanDeploymentPayloadCoordinatesTableMigrations...)
	schemaMigrations = append(schemaMigrations, metaInstanceOperationTableMigrations...)
	schemaMigrations = append(schemaMigrations, metaInstanceRuntimeInstanceTableMigrations...)
	// ++ledgerbuilder:Migrations

	err := simpleDB.ApplyMigrations(schemaMigrations)
	if err != nil {
		return err
	}
	return nil
}
