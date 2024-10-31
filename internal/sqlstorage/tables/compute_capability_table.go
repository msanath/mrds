package tables

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var computeCapabilityTableMigrations = []simplesql.Migration{
	{
		Version: 2, // Update the version number sequentially.
		Up: `
			CREATE TABLE compute_capability (
				id VARCHAR(255) NOT NULL PRIMARY KEY,
				version BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				type VARCHAR(255) NOT NULL,
				score BIGINT NOT NULL,
				UNIQUE (name, deleted_at)
			);
		`,
		Down: `
				DROP TABLE IF EXISTS compute_capability;
			`,
	},
}

type ComputeCapabilityRow struct {
	ID        string `db:"id" orm:"op=create key=primary_key filter=In"`
	Version   uint64 `db:"version" orm:"op=create,update"`
	Name      string `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"`
	DeletedAt int64  `db:"deleted_at"`
	State     string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message   string `db:"message" orm:"op=create,update"`

	Type  string `db:"type" orm:"op=create filter=In,NotIn"`
	Score uint32 `db:"score" orm:"op=create,update"`
}

type ComputeCapabilityKeys struct {
	ID   *string `db:"id"`
	Name *string `db:"name"`
}

type ComputeCapabilityTableUpdateFields struct {
	State     *string `db:"state"`
	Message   *string `db:"message"`
	DeletedAt *int64  `db:"deleted_at"`
}

type ComputeCapabilityTableSelectFilters struct {
	IDIn       []string `db:"id:in"`        // IN condition
	NameIn     []string `db:"name:in"`      // IN condition
	StateIn    []string `db:"state:in"`     // IN condition
	StateNotIn []string `db:"state:not_in"` // NOT IN condition
	VersionGte *uint64  `db:"version:gte"`  // Greater than or equal condition
	VersionLte *uint64  `db:"version:lte"`  // Less than or equal condition
	VersionEq  *uint64  `db:"version:eq"`   // Equal condition

	IncludeDeleted bool   `db:"include_deleted"` // Special boolean handling
	Limit          uint32 `db:"limit"`

	TypeIn []string `db:"type:in"` // IN condition
}

const computeCapabilityTableName = "compute_capability"

type ComputeCapabilityTable struct {
	simplesql.Database
	tableName string
}

func NewComputeCapabilityTable(db simplesql.Database) *ComputeCapabilityTable {
	return &ComputeCapabilityTable{
		Database:  db,
		tableName: computeCapabilityTableName,
	}
}

func (s *ComputeCapabilityTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row ComputeCapabilityRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *ComputeCapabilityTable) Get(ctx context.Context, keys ComputeCapabilityKeys) (ComputeCapabilityRow, error) {
	var row ComputeCapabilityRow
	err := s.Database.GetRowByKey(ctx, s.tableName, keys, &row)
	if err != nil {
		return ComputeCapabilityRow{}, err
	}
	return row, nil
}

func (s *ComputeCapabilityTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, id string, version uint64, updateFields ComputeCapabilityTableUpdateFields,
) error {
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, updateFields)
}

func (s *ComputeCapabilityTable) Delete(ctx context.Context, execer sqlx.ExecerContext, id string, version uint64) error {
	timeNow := time.Now().Unix()
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, ComputeCapabilityTableUpdateFields{
		DeletedAt: &timeNow,
	})
}

func (s *ComputeCapabilityTable) List(ctx context.Context, filters ComputeCapabilityTableSelectFilters) ([]ComputeCapabilityRow, error) {
	var rows []ComputeCapabilityRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
