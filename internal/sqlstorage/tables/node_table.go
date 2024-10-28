package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var nodeTableMigrations = []simplesql.Migration{
	{
		Version: 3, // Update the version number sequentially.
		Up: `
			CREATE TABLE node (
				id VARCHAR(255) NOT NULL PRIMARY KEY,
				version BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				update_domain VARCHAR(255) NOT NULL,
				cluster_id VARCHAR(255) NOT NULL,
				total_cores INT NOT NULL,
				total_memory INT NOT NULL,
				system_reserved_cores INT NOT NULL,
				system_reserved_memory INT NOT NULL,
				remaning_cores INT NOT NULL,
				remaning_memory INT NOT NULL,
				UNIQUE (id, name, is_deleted)
			);
		`,
		Down: `
				DROP TABLE IF EXISTS node;
			`,
	},
}

type NodeRow struct {
	ID                   string `db:"id" orm:"op=create key=primary_key filter=In"`
	Version              uint64 `db:"version" orm:"op=create,update"`
	Name                 string `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"`
	IsDeleted            bool   `db:"is_deleted"`
	State                string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message              string `db:"message" orm:"op=create,update"`
	UpdateDomain         string `db:"update_domain" orm:"op=create filter=In"`
	ClusterID            string `db:"cluster_id" orm:"op=create filter=In"`
	TotalCores           uint32 `db:"total_cores" orm:"op=create,update"`
	TotalMemory          uint32 `db:"total_memory" orm:"op=create,update"`
	SystemReservedCores  uint32 `db:"system_reserved_cores" orm:"op=create,update"`
	SystemReservedMemory uint32 `db:"system_reserved_memory" orm:"op=create,update"`
	RemainingCores       uint32 `db:"remaning_cores" orm:"op=create,update filter=lte,gte"`
	RemainingMemory      uint32 `db:"remaning_memory" orm:"op=create,update filter=lte,gte"`
}

type NodeUpdateFields struct {
	State     *string `db:"state"`
	Message   *string `db:"message"`
	ClusterID *string `db:"cluster_id"`
}

type NodeSelectFilters struct {
	IDIn               []string `db:"id:in"`        // IN condition
	NameIn             []string `db:"name:in"`      // IN condition
	StateIn            []string `db:"state:in"`     // IN condition
	StateNotIn         []string `db:"state:not_in"` // NOT IN condition
	VersionGte         *uint64  `db:"version:gte"`  // Greater than or equal condition
	VersionLte         *uint64  `db:"version:lte"`  // Less than or equal condition
	VersionEq          *uint64  `db:"version:eq"`   // Equal condition
	ClusterIDIn        []string `db:"cluster_id:in"`
	ClusterIDNotIn     []string `db:"cluster_id:not_in"`
	RemainingCoresGte  *uint32  `db:"remaning_cores:gte"`  // Greater than or equal condition
	RemainingCoresLte  *uint32  `db:"remaning_cores:lte"`  // Less than or equal condition
	RemainingMemoryGte *uint32  `db:"remaning_memory:gte"` // Greater than or equal condition
	RemainingMemoryLte *uint32  `db:"remaning_memory:lte"` // Less than or equal condition

	IncludeDeleted bool   `db:"include_deleted"` // Special boolean handling
	Limit          uint32 `db:"limit"`
}

const nodeTableName = "node"

type NodeTable struct {
	simplesql.Database
	tableName string
}

func NewNodeTable(db simplesql.Database) *NodeTable {
	return &NodeTable{
		Database:  db,
		tableName: nodeTableName,
	}
}

func (s *NodeTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row NodeRow) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *NodeTable) GetByIDAndVersion(ctx context.Context, id string, version uint64, isDeleted bool) (NodeRow, error) {
	var row NodeRow
	err := s.Database.GetRowByID(ctx, id, version, isDeleted, s.tableName, &row)
	if err != nil {
		return NodeRow{}, err
	}
	return row, nil
}

func (s *NodeTable) GetByName(ctx context.Context, name string) (NodeRow, error) {
	var row NodeRow
	err := s.Database.GetRowByName(ctx, name, s.tableName, &row)
	if err != nil {
		return NodeRow{}, err
	}
	return row, nil
}

func (s *NodeTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, id string, version uint64, updateFields NodeUpdateFields,
) error {
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, updateFields)
}

func (s *NodeTable) Delete(ctx context.Context, execer sqlx.ExecerContext, id string, version uint64) error {
	return s.Database.MarkRowAsDeleted(ctx, execer, id, version, s.tableName)
}

func (s *NodeTable) List(ctx context.Context, filters NodeSelectFilters) ([]NodeRow, error) {
	var rows []NodeRow
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
