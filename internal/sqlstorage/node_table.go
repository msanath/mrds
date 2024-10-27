package sqlstorage

import (
	"context"

	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/node"

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

type nodeRow struct {
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

type nodeUpdateFields struct {
	State     *string `db:"state"`
	Message   *string `db:"message"`
	ClusterID *string `db:"cluster_id"`
}

type nodeSelectFilters struct {
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

func nodeModelToRow(model node.NodeRecord) nodeRow {
	return nodeRow{
		ID:                   model.Metadata.ID,
		Version:              model.Metadata.Version,
		Name:                 model.Name,
		State:                model.Status.State.ToString(),
		Message:              model.Status.Message,
		UpdateDomain:         model.UpdateDomain,
		ClusterID:            model.ClusterID,
		TotalCores:           model.TotalResources.Cores,
		TotalMemory:          model.TotalResources.Memory,
		SystemReservedCores:  model.SystemReservedResources.Cores,
		SystemReservedMemory: model.SystemReservedResources.Memory,
		RemainingCores:       model.RemainingResources.Cores,
		RemainingMemory:      model.RemainingResources.Memory,
	}
}

func nodeRowToModel(row nodeRow) node.NodeRecord {
	return node.NodeRecord{
		Metadata: core.Metadata{
			ID:      row.ID,
			Version: row.Version,
		},
		Name: row.Name,
		Status: node.NodeStatus{
			State:   node.NodeStateFromString(row.State),
			Message: row.Message,
		},
		UpdateDomain: row.UpdateDomain,
		ClusterID:    row.ClusterID,
		TotalResources: node.Resources{
			Cores:  row.TotalCores,
			Memory: row.TotalMemory,
		},
		SystemReservedResources: node.Resources{
			Cores:  row.SystemReservedCores,
			Memory: row.SystemReservedMemory,
		},
		RemainingResources: node.Resources{
			Cores:  row.RemainingCores,
			Memory: row.RemainingMemory,
		},
	}
}

// nodeStorage is a concrete implementation of NodeRepository using sqlx
type nodeStorage struct {
	simplesql.Database
	tableName  string
	modelToRow func(node.NodeRecord) nodeRow
	rowToModel func(nodeRow) node.NodeRecord
}

// newNodeStorage creates a new storage instance satisfying the NodeRepository interface
func newNodeStorage(db simplesql.Database) node.Repository {
	return &nodeStorage{
		Database:   db,
		tableName:  nodeTableName,
		modelToRow: nodeModelToRow,
		rowToModel: nodeRowToModel,
	}
}

func (s *nodeStorage) Insert(ctx context.Context, record node.NodeRecord) error {
	row := s.modelToRow(record)
	err := s.Database.InsertRow(ctx, s.DB, s.tableName, row)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *nodeStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (node.NodeRecord, error) {
	var row nodeRow
	err := s.Database.GetRowByID(ctx, metadata.ID, metadata.Version, metadata.IsDeleted, s.tableName, &row)
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *nodeStorage) GetByName(ctx context.Context, nodeName string) (node.NodeRecord, error) {
	var row nodeRow
	err := s.Database.GetRowByName(ctx, nodeName, s.tableName, &row)
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *nodeStorage) UpdateState(ctx context.Context, metadata core.Metadata, status node.NodeStatus) error {
	state := status.State.ToString()
	err := s.Database.UpdateRow(ctx, s.DB, metadata.ID, metadata.Version, s.tableName, nodeUpdateFields{
		State:   &state,
		Message: &status.Message,
	})
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *nodeStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	err := s.Database.MarkRowAsDeleted(ctx, s.DB, metadata.ID, metadata.Version, s.tableName)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *nodeStorage) List(ctx context.Context, filters node.NodeListFilters) ([]node.NodeRecord, error) {
	// Extract core filters
	dbFilters := nodeSelectFilters{
		IDIn:               append([]string{}, filters.IDIn...),
		NameIn:             append([]string{}, filters.NameIn...),
		VersionGte:         filters.VersionGte,
		VersionLte:         filters.VersionLte,
		VersionEq:          filters.VersionEq,
		IncludeDeleted:     filters.IncludeDeleted,
		Limit:              filters.Limit,
		RemainingCoresGte:  filters.RemainingCoresGte,
		RemainingCoresLte:  filters.RemainingCoresLte,
		RemainingMemoryGte: filters.RemainingMemoryGte,
		RemainingMemoryLte: filters.RemainingMemoryLte,
		ClusterIDIn:        append([]string{}, filters.ClusterIDIn...),
	}

	// Extract node specific filters
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, state.ToString())
	}
	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, state.ToString())
	}

	var rows []nodeRow
	err := s.Database.SelectRows(ctx, s.tableName, dbFilters, &rows)
	if err != nil {
		return nil, err
	}

	var records []node.NodeRecord
	for _, row := range rows {
		records = append(records, s.rowToModel(row))
	}

	return records, nil
}
