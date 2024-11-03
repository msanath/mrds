package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var nodeLocalVolumeTableMigrations = []simplesql.Migration{
	{
		Version: 4, // Update the version number sequentially.
		Up: `
			CREATE TABLE node_local_volume (
				node_id VARCHAR(255) NOT NULL,
				mount_path VARCHAR(255) NOT NULL,
				storage_class VARCHAR(255) NOT NULL,
				storage_capacity INT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				PRIMARY KEY (node_id, mount_path),
				FOREIGN KEY (node_id) REFERENCES node(id) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS node_local_volume;
			`,
	},
}

type NodeLocalVolumeRow struct {
	NodeID          string `db:"node_id" orm:"op=create filter=In"`
	MountPath       string `db:"mount_path" orm:"op=create"`
	StorageClass    string `db:"storage_class" orm:"op=create filter=In"`
	StorageCapacity uint32 `db:"storage_capacity" orm:"op=create"`
}

type NodeLocalVolumeTableSelectFilters struct {
	NodeIDIn       []string `db:"node_id:in"`       // IN condition
	StorageClassIn []string `db:"storage_class:in"` // IN condition
}

const nodeLocalVolumeTableName = "node_local_volume"

type NodeLocalVolumeTable struct {
	simplesql.Database
	tableName string
}

func NewNodeLocalVolumeTable(db simplesql.Database) *NodeLocalVolumeTable {
	return &NodeLocalVolumeTable{
		Database:  db,
		tableName: nodeLocalVolumeTableName,
	}
}

func (s *NodeLocalVolumeTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row NodeLocalVolumeRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *NodeLocalVolumeTable) Delete(ctx context.Context, execer sqlx.ExecerContext, nodeID string) error {
	query := `
		DELETE FROM node_local_volume
		WHERE node_id = :node_id
	`
	params := map[string]interface{}{
		"node_id": nodeID,
	}
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = s.DB.Rebind(query)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

// TODO: This should also be auto-generated.
func (s *NodeLocalVolumeTable) List(ctx context.Context, filters NodeLocalVolumeTableSelectFilters) ([]NodeLocalVolumeRow, error) {
	var rows []NodeLocalVolumeRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
