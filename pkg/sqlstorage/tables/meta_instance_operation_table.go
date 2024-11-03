package tables

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var metaInstanceOperationTableMigrations = []simplesql.Migration{
	{
		Version: 16, // Update the version number sequentially.
		Up: `
			CREATE TABLE meta_instance_operation (
				id VARCHAR(255) NOT NULL,
				meta_instance_id VARCHAR(255) NOT NULL,
				type VARCHAR(255) NOT NULL,
				intent_id VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				PRIMARY KEY (id, meta_instance_id),
				FOREIGN KEY (meta_instance_id) REFERENCES meta_instance(id) ON DELETE CASCADE
			);
		`,
		Down: `
				DROP TABLE IF EXISTS meta_instance_operation;
			`,
	},
}

type MetaInstanceOperationRow struct {
	ID             string `db:"id" orm:"op=create key=primary_key filter=In"`
	MetaInstanceID string `db:"meta_instance_id" orm:"op=create filter=In"`
	Type           string `db:"type" orm:"op=create"`
	IntentID       string `db:"intent_id" orm:"op=create"`
	State          string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message        string `db:"message" orm:"op=create,update"`
}

type MetaInstanceOperationTableUpdateFields struct {
	State   *string `db:"state"`
	Message *string `db:"message"`
}

type MetaInstanceOperationTableSelectFilters struct {
	IDIn             []string `db:"id:in"`               // IN condition
	MetaInstanceIDIn []string `db:"meta_instance_id:in"` // IN condition
	StateIn          []string `db:"state:in"`            // IN condition
	StateNotIn       []string `db:"state:not_in"`        // NOT IN condition
}

const metaInstanceOperationTableName = "meta_instance_operation"

type MetaInstanceOperationTable struct {
	simplesql.Database
	tableName string
}

func NewMetaInstanceOperationTable(db simplesql.Database) *MetaInstanceOperationTable {
	return &MetaInstanceOperationTable{
		Database:  db,
		tableName: metaInstanceOperationTableName,
	}
}

func (s *MetaInstanceOperationTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row MetaInstanceOperationRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

// TODO: Auto-generate this function
func (s *MetaInstanceOperationTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, operationID string, metaInstanceID string, updateFields MetaInstanceOperationTableUpdateFields,
) error {
	query := `
		UPDATE meta_instance_operation
		SET
	`
	params := map[string]interface{}{
		"id":               operationID,
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

	query += strings.Join(updates, ", ") + " WHERE id = :id AND meta_instance_id = :meta_instance_id"
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = s.DB.Rebind(query)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

func (s *MetaInstanceOperationTable) Delete(ctx context.Context, execer sqlx.ExecerContext, operationID string, metaInstanceID string) error {
	query := `
		DELETE FROM meta_instance_operation
		WHERE id = :id AND meta_instance_id = :meta_instance_id
	`
	params := map[string]interface{}{
		"id":               operationID,
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

func (s *MetaInstanceOperationTable) List(ctx context.Context, filters MetaInstanceOperationTableSelectFilters) ([]MetaInstanceOperationRow, error) {
	var rows []MetaInstanceOperationRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
