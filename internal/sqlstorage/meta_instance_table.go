package sqlstorage

import (
	"context"

	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/metainstance"

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

type metaInstanceRow struct {
	ID        string `db:"id" orm:"op=create key=primary_key filter=In"`
	Version   uint64 `db:"version" orm:"op=create,update"`
	Name      string `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"`
	IsDeleted bool   `db:"is_deleted"`
	State     string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message   string `db:"message" orm:"op=create,update"`
}

type metaInstanceUpdateFields struct {
	State   *string `db:"state"`
	Message *string `db:"message"`
}

type metaInstanceSelectFilters struct {
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

func metaInstanceModelToRow(model metainstance.MetaInstanceRecord) metaInstanceRow {
	return metaInstanceRow{
		ID:      model.Metadata.ID,
		Version: model.Metadata.Version,
		Name:    model.Name,
		State:   model.Status.State.ToString(),
		Message: model.Status.Message,
	}
}

func metaInstanceRowToModel(row metaInstanceRow) metainstance.MetaInstanceRecord {
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

// metaInstanceStorage is a concrete implementation of MetaInstanceRepository using sqlx
type metaInstanceStorage struct {
	simplesql.Database
	tableName  string
	modelToRow func(metainstance.MetaInstanceRecord) metaInstanceRow
	rowToModel func(metaInstanceRow) metainstance.MetaInstanceRecord
}

// newMetaInstanceStorage creates a new storage instance satisfying the MetaInstanceRepository interface
func newMetaInstanceStorage(db simplesql.Database) metainstance.Repository {
	return &metaInstanceStorage{
		Database:   db,
		tableName:  metaInstanceTableName,
		modelToRow: metaInstanceModelToRow,
		rowToModel: metaInstanceRowToModel,
	}
}

func (s *metaInstanceStorage) Insert(ctx context.Context, record metainstance.MetaInstanceRecord) error {
	row := s.modelToRow(record)
	err := s.Database.InsertRow(ctx, s.DB, s.tableName, row)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (metainstance.MetaInstanceRecord, error) {
	var row metaInstanceRow
	err := s.Database.GetRowByID(ctx, metadata.ID, metadata.Version, metadata.IsDeleted, s.tableName, &row)
	if err != nil {
		return metainstance.MetaInstanceRecord{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *metaInstanceStorage) GetByName(ctx context.Context, metaInstanceName string) (metainstance.MetaInstanceRecord, error) {
	var row metaInstanceRow
	err := s.Database.GetRowByName(ctx, metaInstanceName, s.tableName, &row)
	if err != nil {
		return metainstance.MetaInstanceRecord{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *metaInstanceStorage) UpdateState(ctx context.Context, metadata core.Metadata, status metainstance.MetaInstanceStatus) error {
	state := status.State.ToString()
	err := s.Database.UpdateRow(ctx, s.DB, metadata.ID, metadata.Version, s.tableName, metaInstanceUpdateFields{
		State:   &state,
		Message: &status.Message,
	})
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	err := s.Database.MarkRowAsDeleted(ctx, s.DB, metadata.ID, metadata.Version, s.tableName)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) List(ctx context.Context, filters metainstance.MetaInstanceListFilters) ([]metainstance.MetaInstanceRecord, error) {
	// Extract core filters
	dbFilters := metaInstanceSelectFilters{
		IDIn:           append([]string{}, filters.IDIn...),
		NameIn:         append([]string{}, filters.NameIn...),
		VersionGte:     filters.VersionGte,
		VersionLte:     filters.VersionLte,
		VersionEq:      filters.VersionEq,
		IncludeDeleted: filters.IncludeDeleted,
		Limit:          filters.Limit,
	}

	// Extract metaInstance specific filters
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, state.ToString())
	}
	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, state.ToString())
	}

	var rows []metaInstanceRow
	err := s.Database.SelectRows(ctx, s.tableName, dbFilters, &rows)
	if err != nil {
		return nil, err
	}

	var records []metainstance.MetaInstanceRecord
	for _, row := range rows {
		records = append(records, s.rowToModel(row))
	}

	return records, nil
}
