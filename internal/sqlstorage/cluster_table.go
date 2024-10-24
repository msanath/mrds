package sqlstorage

import (
	"context"

	"github.com/msanath/mrds/internal/ledger/cluster"
	"github.com/msanath/mrds/internal/ledger/core"

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
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				UNIQUE (id, name, is_deleted)
			);
		`,
		Down: `
				DROP TABLE IF EXISTS cluster;
			`,
	},
}

type clusterRow struct {
	ID        string `db:"id" orm:"op=create key=primary_key filter=In"`
	Version   uint64 `db:"version" orm:"op=create,update"`
	Name      string `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"`
	IsDeleted bool   `db:"is_deleted"`
	State     string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message   string `db:"message" orm:"op=create,update"`
}

type clusterUpdateFields struct {
	State   *string `db:"state"`
	Message *string `db:"message"`
}

type clusterSelectFilters struct {
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

func clusterModelToRow(model cluster.ClusterRecord) clusterRow {
	return clusterRow{
		ID:      model.Metadata.ID,
		Version: model.Metadata.Version,
		Name:    model.Name,
		State:   model.Status.State.ToString(),
		Message: model.Status.Message,
	}
}

func clusterRowToModel(row clusterRow) cluster.ClusterRecord {
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

// clusterStorage is a concrete implementation of ClusterRepository using sqlx
type clusterStorage struct {
	simplesql.Database
	tableName  string
	modelToRow func(cluster.ClusterRecord) clusterRow
	rowToModel func(clusterRow) cluster.ClusterRecord
}

// newClusterStorage creates a new storage instance satisfying the ClusterRepository interface
func newClusterStorage(db simplesql.Database) cluster.Repository {
	return &clusterStorage{
		Database:   db,
		tableName:  clusterTableName,
		modelToRow: clusterModelToRow,
		rowToModel: clusterRowToModel,
	}
}

func (s *clusterStorage) Insert(ctx context.Context, record cluster.ClusterRecord) error {
	row := s.modelToRow(record)
	err := s.Database.InsertRow(ctx, s.DB, s.tableName, row)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *clusterStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (cluster.ClusterRecord, error) {
	var row clusterRow
	err := s.Database.GetRowByID(ctx, metadata.ID, metadata.Version, metadata.IsDeleted, s.tableName, &row)
	if err != nil {
		return cluster.ClusterRecord{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *clusterStorage) GetByName(ctx context.Context, clusterName string) (cluster.ClusterRecord, error) {
	var row clusterRow
	err := s.Database.GetRowByName(ctx, clusterName, s.tableName, &row)
	if err != nil {
		return cluster.ClusterRecord{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *clusterStorage) UpdateState(ctx context.Context, metadata core.Metadata, status cluster.ClusterStatus) error {
	state := status.State.ToString()
	err := s.Database.UpdateRow(ctx, s.DB, metadata.ID, metadata.Version, s.tableName, clusterUpdateFields{
		State:   &state,
		Message: &status.Message,
	})
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *clusterStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	err := s.Database.MarkRowAsDeleted(ctx, s.DB, metadata.ID, metadata.Version, s.tableName)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *clusterStorage) List(ctx context.Context, filters cluster.ClusterListFilters) ([]cluster.ClusterRecord, error) {
	// Extract core filters
	dbFilters := clusterSelectFilters{
		IDIn:           append([]string{}, filters.IDIn...),
		NameIn:         append([]string{}, filters.NameIn...),
		VersionGte:     filters.VersionGte,
		VersionLte:     filters.VersionLte,
		VersionEq:      filters.VersionEq,
		IncludeDeleted: filters.IncludeDeleted,
		Limit:          filters.Limit,
	}

	// Extract cluster specific filters
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, state.ToString())
	}
	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, state.ToString())
	}

	var rows []clusterRow
	err := s.Database.SelectRows(ctx, s.tableName, dbFilters, &rows)
	if err != nil {
		return nil, err
	}

	var records []cluster.ClusterRecord
	for _, row := range rows {
		records = append(records, s.rowToModel(row))
	}

	return records, nil
}
