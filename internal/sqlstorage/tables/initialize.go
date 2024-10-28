package tables

import "github.com/msanath/gondolf/pkg/simplesql"

func Initialize(simpleDB simplesql.Database) error {
	var schemaMigrations = []simplesql.Migration{}

	schemaMigrations = append(schemaMigrations, clusterTableMigrations...)
	schemaMigrations = append(schemaMigrations, computeCapabilityTableMigrations...)
	schemaMigrations = append(schemaMigrations, nodeTableMigrations...)
	schemaMigrations = append(schemaMigrations, deploymentTableMigrations...)
	schemaMigrations = append(schemaMigrations, metaInstanceTableMigrations...)
	// ++ledgerbuilder:Migrations

	err := simpleDB.ApplyMigrations(schemaMigrations)
	if err != nil {
		return err
	}
	return nil
}
