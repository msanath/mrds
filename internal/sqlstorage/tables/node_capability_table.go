package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var nodeCapabilityTableMigrations = []simplesql.Migration{
	{
		Version: 8, // Update the version number sequentially.
		Up: `
			CREATE TABLE node_capability (
				node_id VARCHAR(255) NOT NULL,
				capability_id VARCHAR(255) NOT NULL,
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				PRIMARY KEY (node_id, capability_id),
				FOREIGN KEY (node_id) REFERENCES node(id) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS node_capability;
			`,
	},
}

type NodeCapabilityRow struct {
	NodeID       string `db:"node_id" orm:"op=create filter=In"`
	CapabilityID string `db:"capability_id" orm:"op=create filter=In"`
}

type NodeCapabilityTableSelectFilters struct {
	NodeIDIn       []string `db:"node_id:in"`       // IN condition
	CapabilityIDIn []string `db:"capability_id:in"` // IN condition
}

const nodeCapabilityTableName = "node_capability"

type NodeCapabilityTable struct {
	simplesql.Database
	tableName string
}

func NewNodeCapabilityTable(db simplesql.Database) *NodeCapabilityTable {
	return &NodeCapabilityTable{
		Database:  db,
		tableName: nodeCapabilityTableName,
	}
}

func (s *NodeCapabilityTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row NodeCapabilityRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *NodeCapabilityTable) Delete(ctx context.Context, execer sqlx.ExecerContext, nodeID string, capabilityID string) error {
	query := `
		DELETE FROM node_capability
		WHERE node_id = :node_id and capability_id = :capability_id
	`
	params := map[string]interface{}{
		"node_id":       nodeID,
		"capability_id": capabilityID,
	}
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = s.DB.Rebind(query)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

func (s *NodeCapabilityTable) List(ctx context.Context, filters NodeCapabilityTableSelectFilters) ([]NodeCapabilityRow, error) {
	var rows []NodeCapabilityRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
