package sqlstorage

import (
	"context"

	"github.com/msanath/gondolf/pkg/simplesql"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/node"
	"github.com/msanath/mrds/internal/sqlstorage/tables"
)

// nodeStorage is a concrete implementation of NodeRepository using sqlx
type nodeStorage struct {
	simplesql.Database
	nodeTable *tables.NodeTable
}

func nodeModelToRow(model node.NodeRecord) tables.NodeRow {
	return tables.NodeRow{
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

func nodeRowToModel(row tables.NodeRow) node.NodeRecord {
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

// newNodeStorage creates a new storage instance satisfying the NodeRepository interface
func newNodeStorage(db simplesql.Database) node.Repository {
	return &nodeStorage{
		Database:  db,
		nodeTable: tables.NewNodeTable(db),
	}
}

func (s *nodeStorage) Insert(ctx context.Context, record node.NodeRecord) error {
	execer := s.DB
	err := s.nodeTable.Insert(ctx, execer, nodeModelToRow(record))
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *nodeStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (node.NodeRecord, error) {
	row, err := s.nodeTable.GetByIDAndVersion(ctx, metadata.ID, metadata.Version, metadata.IsDeleted)
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	return nodeRowToModel(row), nil
}

func (s *nodeStorage) GetByName(ctx context.Context, nodeName string) (node.NodeRecord, error) {
	row, err := s.nodeTable.GetByName(ctx, nodeName)
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	return nodeRowToModel(row), nil
}

func (s *nodeStorage) UpdateState(ctx context.Context, metadata core.Metadata, status node.NodeStatus) error {
	execer := s.DB
	state := status.State.ToString()
	message := status.Message
	updateFields := tables.NodeUpdateFields{
		State:   &state,
		Message: &message,
	}
	err := s.nodeTable.Update(ctx, execer, metadata.ID, metadata.Version, updateFields)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *nodeStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	execer := s.DB
	err := s.nodeTable.Delete(ctx, execer, metadata.ID, metadata.Version)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *nodeStorage) List(ctx context.Context, filters node.NodeListFilters) ([]node.NodeRecord, error) {
	// Extract core filters
	dbFilters := tables.NodeSelectFilters{
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

	rows, errs := s.nodeTable.List(ctx, dbFilters)
	if errs != nil {
		return nil, errs
	}

	var records []node.NodeRecord
	for _, row := range rows {
		records = append(records, nodeRowToModel(row))
	}

	return records, nil
}
