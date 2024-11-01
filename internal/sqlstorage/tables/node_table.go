package tables

import (
	"context"
	"strings"
	"time"

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
				deleted_at BIGINT NOT NULL DEFAULT 0,
				update_domain VARCHAR(255) NOT NULL,
				cluster_id VARCHAR(255) NOT NULL,
				total_cores INT NOT NULL,
				total_memory INT NOT NULL,
				system_reserved_cores INT NOT NULL,
				system_reserved_memory INT NOT NULL,
				remaining_cores INT NOT NULL,
				remaining_memory INT NOT NULL,
				UNIQUE (name, deleted_at)
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
	DeletedAt            int64  `db:"deleted_at"`
	State                string `db:"state" orm:"op=create,update filter=In,NotIn"`
	Message              string `db:"message" orm:"op=create,update"`
	UpdateDomain         string `db:"update_domain" orm:"op=create filter=In"`
	ClusterID            string `db:"cluster_id" orm:"op=create filter=In"`
	TotalCores           uint32 `db:"total_cores" orm:"op=create,update"`
	TotalMemory          uint32 `db:"total_memory" orm:"op=create,update"`
	SystemReservedCores  uint32 `db:"system_reserved_cores" orm:"op=create,update"`
	SystemReservedMemory uint32 `db:"system_reserved_memory" orm:"op=create,update"`
	RemainingCores       uint32 `db:"remaining_cores" orm:"op=create,update filter=lte,gte"`
	RemainingMemory      uint32 `db:"remaining_memory" orm:"op=create,update filter=lte,gte"`
}

type NodeKeys struct {
	ID   *string `db:"id"`
	Name *string `db:"name"`
}

type NodeUpdateFields struct {
	State           *string `db:"state"`
	Message         *string `db:"message"`
	ClusterID       *string `db:"cluster_id"`
	DeletedAt       *int64  `db:"deleted_at"`
	RemainingCores  *uint32 `db:"remaining_cores"`
	RemainingMemory *uint32 `db:"remaining_memory"`
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
	RemainingCoresGte  *uint32  `db:"remaining_memory:gte"` // Greater than or equal condition
	RemainingCoresLte  *uint32  `db:"remaining_memory:lte"` // Less than or equal condition
	RemainingMemoryGte *uint32  `db:"remaining_memory:gte"` // Greater than or equal condition
	RemainingMemoryLte *uint32  `db:"remaining_memory:lte"` // Less than or equal condition
	PayloadNameIn      []string `db:"payload_name:in"`      // IN condition
	PayloadNameNotIn   []string `db:"payload_name:not_in"`  // NOT IN condition
	UpdateDomainIn     []string `db:"update_domain:in"`

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

func (s *NodeTable) Get(ctx context.Context, keys NodeKeys) (NodeRow, error) {
	var row NodeRow
	err := s.Database.GetRowByKey(ctx, s.tableName, keys, &row)
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
	timeNow := time.Now().Unix()
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, NodeUpdateFields{
		DeletedAt: &timeNow,
	})
}

func (s *NodeTable) List(ctx context.Context, filters NodeSelectFilters) ([]NodeRow, error) {
	// TODO: Optimize this function to reduce the number of queries.
	uniqueIDs, err := s.getUniqueIDs(ctx, filters)
	if err != nil {
		return nil, err
	}

	if len(uniqueIDs) == 0 {
		return nil, nil
	}
	var rows []NodeRow
	err = s.Database.SelectRows(ctx, s.tableName, NodeSelectFilters{
		IDIn:           uniqueIDs,
		IncludeDeleted: filters.IncludeDeleted,
	}, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Get retrieves a NodeRecord by its ID, along with its capabilities and local volumes.
// TODO: Optimize this function to reduce the number of queries.
func (s *NodeTable) getUniqueIDs(ctx context.Context, filters NodeSelectFilters) ([]string, error) {
	query := `
		SELECT DISTINCT n.id
		FROM node n
		LEFT JOIN node_capability nc ON n.id = nc.node_id
		LEFT JOIN node_local_volume lv ON n.id = lv.node_id
		LEFT JOIN node_payload np ON n.id = np.node_id
		WHERE 1=1
	`

	// Prepare dynamic query filters
	var filtersQuery []string
	params := map[string]interface{}{}

	// Apply filters
	if len(filters.IDIn) > 0 {
		filtersQuery = append(filtersQuery, "n.id IN (:id_in)")
		params["id_in"] = filters.IDIn
	}
	if len(filters.NameIn) > 0 {
		filtersQuery = append(filtersQuery, "n.name IN (:name_in)")
		params["name_in"] = filters.NameIn
	}
	if filters.VersionGte != nil {
		filtersQuery = append(filtersQuery, "n.version >= :version_gte")
		params["version_gte"] = *filters.VersionGte
	}
	if filters.VersionLte != nil {
		filtersQuery = append(filtersQuery, "n.version <= :version_lte")
		params["version_lte"] = *filters.VersionLte
	}
	if filters.VersionEq != nil {
		filtersQuery = append(filtersQuery, "n.version = :version_eq")
		params["version_eq"] = *filters.VersionEq
	}
	if filters.IncludeDeleted {
		filtersQuery = append(filtersQuery, "n.deleted_at >= 0")
	} else {
		filtersQuery = append(filtersQuery, "n.deleted_at = 0")
	}
	if len(filters.UpdateDomainIn) > 0 {
		filtersQuery = append(filtersQuery, "n.update_domain IN (:update_domain_in)")
		params["update_domain_in"] = filters.UpdateDomainIn
	}
	if len(filters.ClusterIDIn) > 0 {
		filtersQuery = append(filtersQuery, "n.cluster_id IN (:cluster_id_in)")
		params["cluster_id_in"] = filters.ClusterIDIn
	}
	if len(filters.StateIn) > 0 {
		filtersQuery = append(filtersQuery, "n.state IN (:state_in)")
		states := []string{}
		for _, state := range filters.StateIn {
			states = append(states, string(state))
		}
		params["state_in"] = states
	}
	if len(filters.StateNotIn) > 0 {
		filtersQuery = append(filtersQuery, "n.state NOT IN (:state_not_in)")
		states := []string{}
		for _, state := range filters.StateNotIn {
			states = append(states, string(state))
		}
		params["state_not_in"] = states
	}
	if filters.RemainingCoresGte != nil {
		filtersQuery = append(filtersQuery, "n.remaining_cores >= :remaining_cores_gte")
		params["remaining_cores_gte"] = *filters.RemainingCoresGte
	}
	if filters.RemainingCoresLte != nil {
		filtersQuery = append(filtersQuery, "n.remaining_cores <= :remaining_cores_lte")
		params["remaining_cores_lte"] = *filters.RemainingCoresLte
	}
	if filters.RemainingMemoryGte != nil {
		filtersQuery = append(filtersQuery, "n.remaining_memory >= :remaining_memory_gte")
		params["remaining_memory_gte"] = *filters.RemainingMemoryGte
	}
	if filters.RemainingMemoryLte != nil {
		filtersQuery = append(filtersQuery, "n.remaining_memory <= :remaining_memory_lte")
		params["remaining_memory_lte"] = *filters.RemainingMemoryLte
	}
	if len(filters.PayloadNameIn) > 0 {
		filtersQuery = append(filtersQuery, "np.payload_name IN (:payload_name_in)")
		params["payload_name_in"] = filters.PayloadNameIn
	}
	if len(filters.PayloadNameNotIn) > 0 {
		filtersQuery = append(filtersQuery, "np.payload_name NOT IN (:payload_name_not_in) or np.payload_name IS NULL")
		params["payload_name_not_in"] = filters.PayloadNameNotIn
	}

	if len(filtersQuery) > 0 {
		query += " AND " + strings.Join(filtersQuery, " AND ")
	}

	if filters.Limit > 0 {
		query += " LIMIT :limit"
		params["limit"] = filters.Limit
	}

	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return nil, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = s.DB.Rebind(query)
	var ids []struct {
		ID string `db:"id"`
	}
	err = s.DB.SelectContext(ctx, &ids, query, args...)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	uniqueIDs := make([]string, len(ids))
	for i, id := range ids {
		uniqueIDs[i] = id.ID
	}
	return uniqueIDs, nil
}
