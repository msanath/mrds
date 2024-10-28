package sqlstorage

import (
	"context"

	"github.com/msanath/gondolf/pkg/simplesql"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/metainstance"
	"github.com/msanath/mrds/internal/sqlstorage/tables"
)

// metaInstanceStorage is a concrete implementation of MetaInstanceRepository using sqlx
type metaInstanceStorage struct {
	simplesql.Database
	metaInstanceTable *tables.MetaInstanceTable
}

// newMetaInstanceStorage creates a new storage instance satisfying the MetaInstanceRepository interface
func newMetaInstanceStorage(db simplesql.Database) metainstance.Repository {
	return &metaInstanceStorage{
		Database:          db,
		metaInstanceTable: tables.NewMetaInstanceTable(db),
	}
}

func metaInstanceModelToRow(model metainstance.MetaInstanceRecord) tables.MetaInstanceRow {
	return tables.MetaInstanceRow{
		ID:      model.Metadata.ID,
		Version: model.Metadata.Version,
		Name:    model.Name,
		State:   model.Status.State.ToString(),
		Message: model.Status.Message,
	}
}

func metaInstanceRowToModel(row tables.MetaInstanceRow) metainstance.MetaInstanceRecord {
	return metainstance.MetaInstanceRecord{
		Metadata: core.Metadata{
			ID:      row.ID,
			Version: row.Version,
		},
		Name: row.Name,
		Status: metainstance.MetaInstanceStatus{
			State:   metainstance.MetaInstanceStateFromString(row.State),
			Message: row.Message,
		},
	}
}

func (s *metaInstanceStorage) Insert(ctx context.Context, record metainstance.MetaInstanceRecord) error {
	execer := s.DB
	err := s.metaInstanceTable.Insert(ctx, execer, metaInstanceModelToRow(record))
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (metainstance.MetaInstanceRecord, error) {
	row, err := s.metaInstanceTable.GetByIDAndVersion(ctx, metadata.ID, metadata.Version, metadata.IsDeleted)
	if err != nil {
		return metainstance.MetaInstanceRecord{}, errHandler(err)
	}
	return metaInstanceRowToModel(row), nil
}

func (s *metaInstanceStorage) GetByName(ctx context.Context, name string) (metainstance.MetaInstanceRecord, error) {
	row, err := s.metaInstanceTable.GetByName(ctx, name)
	if err != nil {
		return metainstance.MetaInstanceRecord{}, errHandler(err)
	}
	return metaInstanceRowToModel(row), nil
}

func (s *metaInstanceStorage) UpdateState(ctx context.Context, metadata core.Metadata, status metainstance.MetaInstanceStatus) error {
	execer := s.DB
	state := status.State.ToString()
	message := status.Message
	updateFields := tables.MetaInstanceTableUpdateFields{
		State:   &state,
		Message: &message,
	}
	return s.metaInstanceTable.Update(ctx, execer, metadata.ID, metadata.Version, updateFields)
}

func (s *metaInstanceStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	execer := s.DB
	return s.metaInstanceTable.Delete(ctx, execer, metadata.ID, metadata.Version)
}

func (s *metaInstanceStorage) List(ctx context.Context, filters metainstance.MetaInstanceListFilters) ([]metainstance.MetaInstanceRecord, error) {
	dbFilters := tables.MetaInstanceTableSelectFilters{
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

	rows, err := s.metaInstanceTable.List(ctx, dbFilters)
	if err != nil {
		return nil, err
	}
	var records []metainstance.MetaInstanceRecord
	for _, row := range rows {
		records = append(records, metaInstanceRowToModel(row))
	}
	return records, nil
}
