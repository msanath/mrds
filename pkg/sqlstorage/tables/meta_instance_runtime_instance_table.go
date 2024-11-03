package tables

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var metaInstanceRuntimeInstanceTableMigrations = []simplesql.Migration{
	{
		Version: 15, // Update the version number sequentially.
		Up: `
			CREATE TABLE meta_instance_runtime_instance (
				id VARCHAR(255) NOT NULL PRIMARY KEY,
				meta_instance_id VARCHAR(255) NOT NULL,
				node_id VARCHAR(255) NOT NULL,
				is_active BOOLEAN NOT NULL DEFAULT FALSE,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				FOREIGN KEY (meta_instance_id) REFERENCES meta_instance(id) ON DELETE CASCADE,
				FOREIGN KEY (node_id) REFERENCES node(id) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS meta_instance_runtime_instance;
			`,
	},
}

type MetaInstanceRuntimeInstanceRow struct {
	ID             string `db:"id" orm:"op=create key=primary_key filter=In"`
	MetaInstanceID string `db:"meta_instance_id" orm:"op=create filter=In"`
	NodeID         string `db:"node_id" orm:"op=create filter=In"`
	IsActive       bool   `db:"is_active" orm:"op=create,update filter=In"`
	State          string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message        string `db:"message" orm:"op=create,update"`
}

type MetaInstanceRuntimeInstanceTableUpdateFields struct {
	State    *string `db:"state"`
	Message  *string `db:"message"`
	IsActive *bool   `db:"is_active"`
}

type MetaInstanceRuntimeInstanceTableSelectFilters struct {
	IDIn             []string `db:"id:in"`               // IN condition
	MetaInstanceIDIn []string `db:"meta_instance_id:in"` // IN condition
	NodeIDIn         []string `db:"node_id:in"`          // IN condition
	IsActive         *bool    `db:"is_active"`           // Equal condition
	StateIn          []string `db:"state:in"`            // IN condition
	StateNotIn       []string `db:"state:not_in"`        // NOT IN condition
}

const metaInstanceRuntimeInstanceTableName = "meta_instance_runtime_instance"

type MetaInstanceRuntimeInstanceTable struct {
	simplesql.Database
	tableName string
}

func NewMetaInstanceRuntimeInstanceTable(db simplesql.Database) *MetaInstanceRuntimeInstanceTable {
	return &MetaInstanceRuntimeInstanceTable{
		Database:  db,
		tableName: metaInstanceRuntimeInstanceTableName,
	}
}

func (s *MetaInstanceRuntimeInstanceTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row MetaInstanceRuntimeInstanceRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *MetaInstanceRuntimeInstanceTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, runtimeInstanceID string, metaInstanceID string, updateFields MetaInstanceRuntimeInstanceTableUpdateFields,
) error {
	query := `
		UPDATE meta_instance_runtime_instance
		SET
	`
	params := map[string]interface{}{
		"id":               runtimeInstanceID,
		"meta_instance_id": metaInstanceID,
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

	if updateFields.IsActive != nil {
		updates = append(updates, "is_active = :is_active")
		params["is_active"] = *updateFields.IsActive
	}

	query += strings.Join(updates, ", ") + " WHERE id = :id AND meta_instance_id = :meta_instance_id"
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = s.DB.Rebind(query)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

func (s *MetaInstanceRuntimeInstanceTable) Delete(ctx context.Context, execer sqlx.ExecerContext, runtimeInstanceID string, metaInstanceID string) error {
	query := `
		DELETE FROM meta_instance_runtime_instance
		WHERE id = :id AND meta_instance_id = :meta_instance_id
	`
	params := map[string]interface{}{
		"id":               runtimeInstanceID,
		"meta_instance_id": metaInstanceID,
	}
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = s.DB.Rebind(query)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

func (s *MetaInstanceRuntimeInstanceTable) List(ctx context.Context, filters MetaInstanceRuntimeInstanceTableSelectFilters) ([]MetaInstanceRuntimeInstanceRow, error) {
	var rows []MetaInstanceRuntimeInstanceRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
