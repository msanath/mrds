package tables

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var nodeDisruptionTableMigrations = []simplesql.Migration{
	{
		Version: 6, // Update the version number sequentially.
		Up: `
			CREATE TABLE node_disruption (
				id VARCHAR(255) NOT NULL,
				node_id VARCHAR(255) NOT NULL,
				evict_node BOOLEAN NOT NULL,
				start_time BIGINT NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				PRIMARY KEY (id, node_id),
				FOREIGN KEY (node_id) REFERENCES node(id) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS node_disruption;
			`,
	},
}

type NodeDisruptionRow struct {
	ID        string `db:"id" orm:"op=create key=primary_key filter=In"`
	NodeID    string `db:"node_id" orm:"op=create filter=In"`
	EvictNode bool   `db:"evict_node" orm:"op=create,update filter=In"`
	StartTime uint64 `db:"start_time" orm:"op=create,update filter=gte,lte"`
	State     string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message   string `db:"message" orm:"op=create,update"`
}

type NodeDisruptionTableUpdateFields struct {
	State   *string `db:"state"`
	Message *string `db:"message"`
}

type NodeDisruptionTableSelectFilters struct {
	DisruptionIDIn []string `db:"disruption_id:in"` // IN condition
	NodeIDIn       []string `db:"node_id:in"`       // IN condition
	StateIn        []string `db:"state:in"`         // IN condition
	StateNotIn     []string `db:"state:not_in"`     // NOT IN condition
	StartTimeGte   *uint64  `db:"start_time:gte"`   // Greater than or equal condition
	StartTimeLte   *uint64  `db:"start_time:lte"`   // Less than or equal condition

	IncludeDeleted bool   `db:"include_deleted"` // Special boolean handling
	Limit          uint32 `db:"limit"`
}

const nodeDisruptionTableName = "node_disruption"

type NodeDisruptionTable struct {
	simplesql.Database
	tableName string
}

func NewNodeDisruptionTable(db simplesql.Database) *NodeDisruptionTable {
	return &NodeDisruptionTable{
		Database:  db,
		tableName: nodeDisruptionTableName,
	}
}

func (s *NodeDisruptionTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row NodeDisruptionRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

// TODO: This should also be auto-generated.
func (s *NodeDisruptionTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, nodeID string, id string, updateFields NodeDisruptionTableUpdateFields,
) error {
	query := `
		UPDATE node_disruption
		SET
	`
	params := map[string]interface{}{
		"id":      id,
		"node_id": nodeID,
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

	query += strings.Join(updates, ", ") + " WHERE id = :id AND node_id = :node_id"
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = s.DB.Rebind(query)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

// TODO: This should also be auto-generated.
func (s *NodeDisruptionTable) Delete(ctx context.Context, execer sqlx.ExecerContext, nodeID string, id string) error {
	query := `
		DELETE FROM node_disruption
		WHERE id = :id AND node_id = :node_id
	`
	params := map[string]interface{}{
		"id":      id,
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

func (s *NodeDisruptionTable) List(ctx context.Context, filters NodeDisruptionTableSelectFilters) ([]NodeDisruptionRow, error) {
	var rows []NodeDisruptionRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
