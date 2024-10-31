package tables

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var clusterTableMigrations = []simplesql.Migration{
	{
		Version: 1, // Update the version number sequentially.
		Up: `
			CREATE TABLE cluster (
				id VARCHAR(255) NOT NULL PRIMARY KEY,
				version BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				UNIQUE (name, deleted_at)
			);
		`,
		Down: `
				DROP TABLE IF EXISTS cluster;
			`,
	},
}

type ClusterRow struct {
	ID        string `db:"id" orm:"op=create key=primary_key filter=In"`
	Version   uint64 `db:"version" orm:"op=create,update"`
	Name      string `db:"name" orm:"op=create composite_unique_key:name,deleted_at filter=In"`
	DeletedAt int64  `db:"deleted_at"`
	State     string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message   string `db:"message" orm:"op=create,update"`
}

type ClusterTableKeys struct {
	ID   *string `db:"id"`
	Name *string `db:"name"`
}

type ClusterTableUpdateFields struct {
	State     *string `db:"state"`
	Message   *string `db:"message"`
	DeletedAt *int64  `db:"deleted_at"`
}

type ClusterTableSelectFilters struct {
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

const clusterTableName = "cluster"

type ClusterTable struct {
	simplesql.Database
	tableName string
}

func NewClusterTable(db simplesql.Database) *ClusterTable {
	return &ClusterTable{
		Database:  db,
		tableName: clusterTableName,
	}
}

func (s *ClusterTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row ClusterRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *ClusterTable) Get(ctx context.Context, keys ClusterTableKeys) (ClusterRow, error) {
	var row ClusterRow
	err := s.Database.GetRowByKey(ctx, s.tableName, keys, &row)
	if err != nil {
		return ClusterRow{}, err
	}
	return row, nil
}

func (s *ClusterTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, id string, version uint64, updateFields ClusterTableUpdateFields,
) error {
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, updateFields)
}

func (s *ClusterTable) Delete(ctx context.Context, execer sqlx.ExecerContext, id string, version uint64) error {
	timeNow := time.Now().Unix()
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, ClusterTableUpdateFields{
		DeletedAt: &timeNow,
	})
}

func (s *ClusterTable) List(ctx context.Context, filters ClusterTableSelectFilters) ([]ClusterRow, error) {
	var rows []ClusterRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
