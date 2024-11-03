package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var nodePayloadTableMigrations = []simplesql.Migration{
	{
		Version: 17, // Update the version number sequentially.
		Up: `
			CREATE TABLE node_payload (
				node_id VARCHAR(255) NOT NULL,
				payload_name VARCHAR(255) NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				PRIMARY KEY (node_id, payload_name),
				FOREIGN KEY (node_id) REFERENCES node(id) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS node_payload;
			`,
	},
}

type NodePayloadRow struct {
	NodeID      string `db:"node_id" orm:"op=create key=primary_key filter=In"`
	PayloadName string `db:"payload_name" orm:"op=create key=primary_key filter=In"`
}

type NodePayloadTableSelectFilters struct {
	NodeIDIn         []string `db:"node_id:in"`          // IN condition
	PayloadNameIn    []string `db:"payload_name:in"`     // IN condition
	PayloadNameNotIn []string `db:"payload_name:not_in"` // NOT IN condition
}

const nodePayloadTableName = "node_payload"

type NodePayloadTable struct {
	simplesql.Database
	tableName string
}

func NewNodePayloadTable(db simplesql.Database) *NodePayloadTable {
	return &NodePayloadTable{
		Database:  db,
		tableName: nodePayloadTableName,
	}
}

func (s *NodePayloadTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row NodePayloadRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *NodePayloadTable) Delete(ctx context.Context, execer sqlx.ExecerContext, nodeID string, payloadName string) error {
	query := `
		DELETE FROM node_payload
		WHERE node_id = :node_id AND payload_name = :payload_name
	`
	params := map[string]interface{}{
		"node_id":      nodeID,
		"payload_name": payloadName,
	}
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = s.DB.Rebind(query)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

func (s *NodePayloadTable) List(ctx context.Context, filters NodePayloadTableSelectFilters) ([]NodePayloadRow, error) {
	var rows []NodePayloadRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
