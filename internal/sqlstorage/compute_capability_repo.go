package sqlstorage

import (
	"context"

	"github.com/msanath/mrds/internal/ledger/computecapability"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/sqlstorage/tables"

	"github.com/msanath/gondolf/pkg/simplesql"
)

// computeCapabilityStorage is a concrete implementation of ComputeCapabilityRepository using sqlx
type computeCapabilityStorage struct {
	simplesql.Database
	computeCapabilityTable *tables.ComputeCapabilityTable
}

func computeCapabilityModelToRow(model computecapability.ComputeCapabilityRecord) tables.ComputeCapabilityRow {
	return tables.ComputeCapabilityRow{
		ID:      model.Metadata.ID,
		Version: model.Metadata.Version,
		Name:    model.Name,
		State:   model.Status.State.ToString(),
		Message: model.Status.Message,
		Type:    model.Type,
		Score:   model.Score,
	}
}

func computeCapabilityRowToModel(row tables.ComputeCapabilityRow) computecapability.ComputeCapabilityRecord {
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

// newComputeCapabilityStorage creates a new storage instance satisfying the ComputeCapabilityRepository interface
func newComputeCapabilityStorage(db simplesql.Database) computecapability.Repository {
	return &computeCapabilityStorage{
		Database:               db,
		computeCapabilityTable: tables.NewComputeCapabilityTable(db),
	}
}

func (s *computeCapabilityStorage) Insert(ctx context.Context, record computecapability.ComputeCapabilityRecord) error {
	execer := s.DB
	err := s.computeCapabilityTable.Insert(ctx, execer, computeCapabilityModelToRow(record))
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *computeCapabilityStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (computecapability.ComputeCapabilityRecord, error) {
	row, err := s.computeCapabilityTable.GetByIDAndVersion(ctx, metadata.ID, metadata.Version, metadata.IsDeleted)
	if err != nil {
		return computecapability.ComputeCapabilityRecord{}, errHandler(err)
	}
	return computeCapabilityRowToModel(row), nil
}

func (s *computeCapabilityStorage) GetByName(ctx context.Context, computeCapabilityName string) (computecapability.ComputeCapabilityRecord, error) {
	row, err := s.computeCapabilityTable.GetByName(ctx, computeCapabilityName)
	if err != nil {
		return computecapability.ComputeCapabilityRecord{}, errHandler(err)
	}
	return computeCapabilityRowToModel(row), nil
}

func (s *computeCapabilityStorage) UpdateState(ctx context.Context, metadata core.Metadata, status computecapability.ComputeCapabilityStatus) error {
	execer := s.DB
	state := status.State.ToString()
	message := status.Message
	updateFields := tables.ComputeCapabilityTableUpdateFields{
		State:   &state,
		Message: &message,
	}
	return s.computeCapabilityTable.Update(ctx, execer, metadata.ID, metadata.Version, updateFields)
}

func (s *computeCapabilityStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	execer := s.DB
	return s.computeCapabilityTable.Delete(ctx, execer, metadata.ID, metadata.Version)
}

func (s *computeCapabilityStorage) List(ctx context.Context, filters computecapability.ComputeCapabilityListFilters) ([]computecapability.ComputeCapabilityRecord, error) {
	// Extract core filters
	dbFilters := tables.ComputeCapabilityTableSelectFilters{
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

	rows, err := s.computeCapabilityTable.List(ctx, dbFilters)
	if err != nil {
		return nil, err
	}

	var records []computecapability.ComputeCapabilityRecord
	for _, row := range rows {
		records = append(records, computeCapabilityRowToModel(row))
	}

	return records, nil
}
