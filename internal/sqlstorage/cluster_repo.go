package sqlstorage

import (
	"context"

	"github.com/msanath/gondolf/pkg/simplesql"
	"github.com/msanath/mrds/internal/ledger/cluster"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/sqlstorage/tables"
)

// clusterStorage is a concrete implementation of ClusterRepository using sqlx
type clusterStorage struct {
	simplesql.Database
	clusterTable *tables.ClusterTable
}

// newClusterStorage creates a new storage instance satisfying the ClusterRepository interface
func newClusterStorage(db simplesql.Database) cluster.Repository {
	return &clusterStorage{
		Database:     db,
		clusterTable: tables.NewClusterTable(db),
	}
}

func clusterModelToRow(model cluster.ClusterRecord) tables.ClusterRow {
	return tables.ClusterRow{
		ID:      model.Metadata.ID,
		Version: model.Metadata.Version,
		Name:    model.Name,
		State:   model.Status.State.ToString(),
		Message: model.Status.Message,
	}
}

func clusterRowToModel(row tables.ClusterRow) cluster.ClusterRecord {
	return cluster.ClusterRecord{
		Metadata: core.Metadata{
			ID:      row.ID,
			Version: row.Version,
		},
		Name: row.Name,
		Status: cluster.ClusterStatus{
			State:   cluster.ClusterStateFromString(row.State),
			Message: row.Message,
		},
	}
}

func (s *clusterStorage) Insert(ctx context.Context, record cluster.ClusterRecord) error {
	execer := s.DB
	err := s.clusterTable.Insert(ctx, execer, clusterModelToRow(record))
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *clusterStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (cluster.ClusterRecord, error) {
	row, err := s.clusterTable.GetByIDAndVersion(ctx, metadata.ID, metadata.Version, metadata.IsDeleted)
	if err != nil {
		return cluster.ClusterRecord{}, errHandler(err)
	}
	return clusterRowToModel(row), nil
}

func (s *clusterStorage) GetByName(ctx context.Context, name string) (cluster.ClusterRecord, error) {
	row, err := s.clusterTable.GetByName(ctx, name)
	if err != nil {
		return cluster.ClusterRecord{}, errHandler(err)
	}
	return clusterRowToModel(row), nil
}

func (s *clusterStorage) UpdateState(ctx context.Context, metadata core.Metadata, status cluster.ClusterStatus) error {
	execer := s.DB
	state := status.State.ToString()
	message := status.Message
	updateFields := tables.ClusterTableUpdateFields{
		State:   &state,
		Message: &message,
	}
	return s.clusterTable.Update(ctx, execer, metadata.ID, metadata.Version, updateFields)
}

func (s *clusterStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	execer := s.DB
	return s.clusterTable.Delete(ctx, execer, metadata.ID, metadata.Version)
}

func (s *clusterStorage) List(ctx context.Context, filters cluster.ClusterListFilters) ([]cluster.ClusterRecord, error) {
	dbFilters := tables.ClusterTableSelectFilters{
		IDIn:           append([]string{}, filters.IDIn...),
		NameIn:         append([]string{}, filters.NameIn...),
		VersionGte:     filters.VersionGte,
		VersionLte:     filters.VersionLte,
		VersionEq:      filters.VersionEq,
		IncludeDeleted: filters.IncludeDeleted,
		Limit:          filters.Limit,
	}
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, state.ToString())
	}
	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, state.ToString())
	}

	rows, err := s.clusterTable.List(ctx, dbFilters)
	if err != nil {
		return nil, err
	}
	var records []cluster.ClusterRecord
	for _, row := range rows {
		records = append(records, clusterRowToModel(row))
	}
	return records, nil
}
