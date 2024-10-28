package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var metaInstanceTableMigrations = []simplesql.Migration{
	{
		Version: 5, // Update the version number sequentially.
		Up: `
			CREATE TABLE meta_instance (
				id VARCHAR(255) NOT NULL PRIMARY KEY,
				version BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				UNIQUE (id, name, is_deleted)
			);
		`,
		Down: `
				DROP TABLE IF EXISTS meta_instance;
			`,
	},
}

type MetaInstanceRow struct {
	ID        string `db:"id" orm:"op=create key=primary_key filter=In"`
	Version   uint64 `db:"version" orm:"op=create,update"`
	Name      string `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"`
	IsDeleted bool   `db:"is_deleted"`
	State     string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message   string `db:"message" orm:"op=create,update"`
}

type MetaInstanceTableUpdateFields struct {
	State   *string `db:"state"`
	Message *string `db:"message"`
}

type MetaInstanceTableSelectFilters struct {
	IDIn       []string `db:"id:in"`        // IN condition
	NameIn     []string `db:"name:in"`      // IN condition
	StateIn    []string `db:"state:in"`     // IN condition
	StateNotIn []string `db:"state:not_in"` // NOT IN condition
	VersionGte *uint64  `db:"version:gte"`  // Greater than or equal condition
	VersionLte *uint64  `db:"version:lte"`  // Less than or equal condition
	VersionEq  *uint64  `db:"version:eq"`   // Equal condition

	IncludeDeleted bool   `db:"include_deleted"` // Special boolean handling
	Limit          uint32 `db:"limit"`
}

const metaInstanceTableName = "meta_instance"

type MetaInstanceTable struct {
	simplesql.Database
	tableName string
}

func NewMetaInstanceTable(db simplesql.Database) *MetaInstanceTable {
	return &MetaInstanceTable{
		Database:  db,
		tableName: metaInstanceTableName,
	}
}

func (s *MetaInstanceTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row MetaInstanceRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *MetaInstanceTable) GetByIDAndVersion(ctx context.Context, id string, version uint64, isDeleted bool) (MetaInstanceRow, error) {
	var row MetaInstanceRow
	err := s.Database.GetRowByID(ctx, id, version, isDeleted, s.tableName, &row)
	if err != nil {
		return MetaInstanceRow{}, err
	}
	return row, nil
}

func (s *MetaInstanceTable) GetByName(ctx context.Context, name string) (MetaInstanceRow, error) {
	var row MetaInstanceRow
	err := s.Database.GetRowByName(ctx, name, s.tableName, &row)
	if err != nil {
		return MetaInstanceRow{}, err
	}
	return row, nil
}

func (s *MetaInstanceTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, id string, version uint64, updateFields MetaInstanceTableUpdateFields,
) error {
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, updateFields)
}

func (s *MetaInstanceTable) Delete(ctx context.Context, execer sqlx.ExecerContext, id string, version uint64) error {
	return s.Database.MarkRowAsDeleted(ctx, execer, id, version, s.tableName)
}

func (s *MetaInstanceTable) List(ctx context.Context, filters MetaInstanceTableSelectFilters) ([]MetaInstanceRow, error) {
	var rows []MetaInstanceRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
