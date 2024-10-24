package sqlstorage

import (
	"context"

	"github.com/msanath/mrds/internal/ledger/computecapability"
	"github.com/msanath/mrds/internal/ledger/core"

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
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				type VARCHAR(255) NOT NULL,
				score BIGINT NOT NULL,
				UNIQUE (id, name, is_deleted)
			);
		`,
		Down: `
				DROP TABLE IF EXISTS compute_capability;
			`,
	},
}

type computeCapabilityRow struct {
	ID        string `db:"id" orm:"op=create key=primary_key filter=In"`
	Version   uint64 `db:"version" orm:"op=create,update"`
	Name      string `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"`
	IsDeleted bool   `db:"is_deleted"`
	State     string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message   string `db:"message" orm:"op=create,update"`

	Type  string `db:"type" orm:"op=create filter=In,NotIn"`
	Score uint32 `db:"score" orm:"op=create,update"`
}

type computeCapabilityUpdateFields struct {
	State   *string `db:"state"`
	Message *string `db:"message"`
}

type computeCapabilitySelectFilters struct {
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

func computeCapabilityModelToRow(model computecapability.ComputeCapabilityRecord) computeCapabilityRow {
	return computeCapabilityRow{
		ID:      model.Metadata.ID,
		Version: model.Metadata.Version,
		Name:    model.Name,
		State:   model.Status.State.ToString(),
		Message: model.Status.Message,
		Type:    model.Type,
		Score:   model.Score,
	}
}

func computeCapabilityRowToModel(row computeCapabilityRow) computecapability.ComputeCapabilityRecord {
	return computecapability.ComputeCapabilityRecord{
		Metadata: core.Metadata{
			ID:      row.ID,
			Version: row.Version,
		},
		Name: row.Name,
		Status: computecapability.ComputeCapabilityStatus{
			State:   computecapability.ComputeCapabilityStateFromString(row.State),
			Message: row.Message,
		},
		Type:  row.Type,
		Score: row.Score,
	}
}

// computeCapabilityStorage is a concrete implementation of ComputeCapabilityRepository using sqlx
type computeCapabilityStorage struct {
	simplesql.Database
	tableName  string
	modelToRow func(computecapability.ComputeCapabilityRecord) computeCapabilityRow
	rowToModel func(computeCapabilityRow) computecapability.ComputeCapabilityRecord
}

// newComputeCapabilityStorage creates a new storage instance satisfying the ComputeCapabilityRepository interface
func newComputeCapabilityStorage(db simplesql.Database) computecapability.Repository {
	return &computeCapabilityStorage{
		Database:   db,
		tableName:  computeCapabilityTableName,
		modelToRow: computeCapabilityModelToRow,
		rowToModel: computeCapabilityRowToModel,
	}
}

func (s *computeCapabilityStorage) Insert(ctx context.Context, record computecapability.ComputeCapabilityRecord) error {
	row := s.modelToRow(record)
	err := s.Database.InsertRow(ctx, s.DB, s.tableName, row)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *computeCapabilityStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (computecapability.ComputeCapabilityRecord, error) {
	var row computeCapabilityRow
	err := s.Database.GetRowByID(ctx, metadata.ID, metadata.Version, metadata.IsDeleted, s.tableName, &row)
	if err != nil {
		return computecapability.ComputeCapabilityRecord{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *computeCapabilityStorage) GetByName(ctx context.Context, computeCapabilityName string) (computecapability.ComputeCapabilityRecord, error) {
	var row computeCapabilityRow
	err := s.Database.GetRowByName(ctx, computeCapabilityName, s.tableName, &row)
	if err != nil {
		return computecapability.ComputeCapabilityRecord{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *computeCapabilityStorage) UpdateState(ctx context.Context, metadata core.Metadata, status computecapability.ComputeCapabilityStatus) error {
	state := status.State.ToString()
	err := s.Database.UpdateRow(ctx, s.DB, metadata.ID, metadata.Version, s.tableName, computeCapabilityUpdateFields{
		State:   &state,
		Message: &status.Message,
	})
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *computeCapabilityStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	err := s.Database.MarkRowAsDeleted(ctx, s.DB, metadata.ID, metadata.Version, s.tableName)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *computeCapabilityStorage) List(ctx context.Context, filters computecapability.ComputeCapabilityListFilters) ([]computecapability.ComputeCapabilityRecord, error) {
	// Extract core filters
	dbFilters := computeCapabilitySelectFilters{
		IDIn:           append([]string{}, filters.IDIn...),
		NameIn:         append([]string{}, filters.NameIn...),
		VersionGte:     filters.VersionGte,
		VersionLte:     filters.VersionLte,
		VersionEq:      filters.VersionEq,
		IncludeDeleted: filters.IncludeDeleted,
		Limit:          filters.Limit,
		TypeIn:         append([]string{}, filters.TypeIn...),
	}

	// Extract computeCapability specific filters
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, state.ToString())
	}
	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, state.ToString())
	}

	var rows []computeCapabilityRow
	err := s.Database.SelectRows(ctx, s.tableName, dbFilters, &rows)
	if err != nil {
		return nil, err
	}

	var records []computecapability.ComputeCapabilityRecord
	for _, row := range rows {
		records = append(records, s.rowToModel(row))
	}

	return records, nil
}
