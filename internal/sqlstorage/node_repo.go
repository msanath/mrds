package sqlstorage

import (
	"context"
	"time"

	"github.com/msanath/gondolf/pkg/simplesql"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/node"
	"github.com/msanath/mrds/internal/sqlstorage/tables"
)

// nodeStorage is a concrete implementation of NodeRepository using sqlx
type nodeStorage struct {
	simplesql.Database
	nodeTable            *tables.NodeTable
	nodeLocalVolumeTable *tables.NodeLocalVolumeTable
	nodeDisruptionTable  *tables.NodeDisruptionTable
}

func nodeRecordToRow(record node.NodeRecord) tables.NodeRow {
	return tables.NodeRow{
		ID:                   record.Metadata.ID,
		Version:              record.Metadata.Version,
		Name:                 record.Name,
		State:                record.Status.State.ToString(),
		Message:              record.Status.Message,
		UpdateDomain:         record.UpdateDomain,
		ClusterID:            record.ClusterID,
		TotalCores:           record.TotalResources.Cores,
		TotalMemory:          record.TotalResources.Memory,
		SystemReservedCores:  record.SystemReservedResources.Cores,
		SystemReservedMemory: record.SystemReservedResources.Memory,
		RemainingCores:       record.RemainingResources.Cores,
		RemainingMemory:      record.RemainingResources.Memory,
	}
}

func nodeRowToRecord(row tables.NodeRow) node.NodeRecord {
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

func nodeDisruptionRecordToRow(nodeID string, record node.NodeDisruption) tables.NodeDisruptionRow {
	return tables.NodeDisruptionRow{
		ID:        record.ID,
		NodeID:    nodeID,
		StartTime: uint64(record.StartTime.Unix()),
		EvictNode: record.EvictNode,
		State:     string(record.Status.State),
		Message:   record.Status.Message,
	}
}

func nodeDisruptionRowToRecord(row tables.NodeDisruptionRow) node.NodeDisruption {
	return node.NodeDisruption{
		ID: row.ID,
		Status: node.NodeDisruptionStatus{
			State:   node.DisruptionState(row.State),
			Message: row.Message,
		},
		StartTime: time.Unix(int64(row.StartTime), 0), // Convert from uint64 to time.Time
		EvictNode: row.EvictNode,
	}
}

func nodeLocalVolumeRecordToRow(nodeID string, record node.NodeLocalVolume) tables.NodeLocalVolumeRow {
	return tables.NodeLocalVolumeRow{
		NodeID:          nodeID,
		MountPath:       record.MountPath,
		StorageClass:    record.StorageClass,
		StorageCapacity: record.StorageCapacity,
	}
}

func nodeLocalVolumeRowToRecord(row tables.NodeLocalVolumeRow) node.NodeLocalVolume {
	return node.NodeLocalVolume{
		MountPath:       row.MountPath,
		StorageClass:    row.StorageClass,
		StorageCapacity: row.StorageCapacity,
	}
}

// newNodeStorage creates a new storage instance satisfying the NodeRepository interface
func newNodeStorage(db simplesql.Database) node.Repository {
	return &nodeStorage{
		Database:             db,
		nodeTable:            tables.NewNodeTable(db),
		nodeLocalVolumeTable: tables.NewNodeLocalVolumeTable(db),
		nodeDisruptionTable:  tables.NewNodeDisruptionTable(db),
	}
}

func (s *nodeStorage) Insert(ctx context.Context, record node.NodeRecord) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	err = s.nodeTable.Insert(ctx, tx, nodeRecordToRow(record))
	if err != nil {
		return errHandler(err)
	}

	for _, localVolume := range record.LocalVolumes {
		err = s.nodeLocalVolumeTable.Insert(ctx, tx, nodeLocalVolumeRecordToRow(record.Metadata.ID, localVolume))
		if err != nil {
			return errHandler(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *nodeStorage) GetByMetadata(ctx context.Context, metadata core.Metadata) (node.NodeRecord, error) {
	nodeRow, err := s.nodeTable.GetByIDAndVersion(ctx, metadata.ID, metadata.Version, metadata.IsDeleted)
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	nodeRecord := nodeRowToRecord(nodeRow)

	// Get the corresponding disruptionRows
	disruptionRows, err := s.nodeDisruptionTable.List(ctx, tables.NodeDisruptionTableSelectFilters{
		NodeIDIn: []string{metadata.ID},
	})
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	for _, disruption := range disruptionRows {
		nodeRecord.Disruptions = append(nodeRecord.Disruptions, nodeDisruptionRowToRecord(disruption))
	}

	// Get the corresponding volumeRows
	volumeRows, err := s.nodeLocalVolumeTable.List(ctx, tables.NodeLocalVolumeTableSelectFilters{
		NodeIDIn: []string{metadata.ID},
	})
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	for _, volume := range volumeRows {
		nodeRecord.LocalVolumes = append(nodeRecord.LocalVolumes, nodeLocalVolumeRowToRecord(volume))
	}

	return nodeRecord, nil
}

func (s *nodeStorage) GetByName(ctx context.Context, nodeName string) (node.NodeRecord, error) {
	nodeRow, err := s.nodeTable.GetByName(ctx, nodeName)
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	nodeRecord := nodeRowToRecord(nodeRow)

	// Get the corresponding disruptionRows
	disruptionRows, err := s.nodeDisruptionTable.List(ctx, tables.NodeDisruptionTableSelectFilters{
		NodeIDIn: []string{nodeRecord.Metadata.ID},
	})
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	for _, disruption := range disruptionRows {
		nodeRecord.Disruptions = append(nodeRecord.Disruptions, nodeDisruptionRowToRecord(disruption))
	}

	// Get the corresponding volumeRows
	volumeRows, err := s.nodeLocalVolumeTable.List(ctx, tables.NodeLocalVolumeTableSelectFilters{
		NodeIDIn: []string{nodeRecord.Metadata.ID},
	})
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	for _, volume := range volumeRows {
		nodeRecord.LocalVolumes = append(nodeRecord.LocalVolumes, nodeLocalVolumeRowToRecord(volume))
	}
	return nodeRecord, nil
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

	// Get all the IDs of the nodes
	var nodeIDs []string
	for _, row := range rows {
		nodeIDs = append(nodeIDs, row.ID)
	}

	// Get the corresponding disruptionRows
	disruptionRows, err := s.nodeDisruptionTable.List(ctx, tables.NodeDisruptionTableSelectFilters{
		NodeIDIn: nodeIDs,
	})
	if err != nil {
		return nil, errHandler(err)
	}

	// Create a map of nodeID to list of disruptions
	disruptionMap := make(map[string][]node.NodeDisruption)
	for _, disruption := range disruptionRows {
		disruptionMap[disruption.NodeID] = append(disruptionMap[disruption.NodeID], nodeDisruptionRowToRecord(disruption))
	}

	volumeRows, err := s.nodeLocalVolumeTable.List(ctx, tables.NodeLocalVolumeTableSelectFilters{
		NodeIDIn: nodeIDs,
	})
	if err != nil {
		return nil, errHandler(err)
	}
	volumeMap := make(map[string][]node.NodeLocalVolume)
	for _, volume := range volumeRows {
		volumeMap[volume.NodeID] = append(volumeMap[volume.NodeID], nodeLocalVolumeRowToRecord(volume))
	}

	var records []node.NodeRecord
	for _, row := range rows {
		record := nodeRowToRecord(row)
		record.Disruptions = disruptionMap[row.ID]
		record.LocalVolumes = volumeMap[row.ID]
		records = append(records, record)
	}

	return records, nil
}

func (s *nodeStorage) InsertDisruption(ctx context.Context, nodeMetadata core.Metadata, record node.NodeDisruption) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	err = s.nodeDisruptionTable.Insert(ctx, execer, nodeDisruptionRecordToRow(nodeMetadata.ID, record))
	if err != nil {
		return errHandler(err)
	}

	// update the node record to bump the version.
	err = s.nodeTable.Update(ctx, execer, nodeMetadata.ID, nodeMetadata.Version, tables.NodeUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	err = tx.Commit()
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *nodeStorage) DeleteDisruption(ctx context.Context, nodeMetadata core.Metadata, disruptionID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	err = s.nodeDisruptionTable.Delete(ctx, execer, nodeMetadata.ID, disruptionID)
	if err != nil {
		return errHandler(err)
	}

	// update the node record to bump the version.
	err = s.nodeTable.Update(ctx, execer, nodeMetadata.ID, nodeMetadata.Version, tables.NodeUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	err = tx.Commit()
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *nodeStorage) UpdateDisruptionStatus(ctx context.Context, nodeMetadata core.Metadata, disruptionID string, status node.NodeDisruptionStatus) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	state := string(status.State)
	message := status.Message
	updateFields := tables.NodeDisruptionTableUpdateFields{
		State:   &state,
		Message: &message,
	}
	err = s.nodeDisruptionTable.Update(ctx, execer, nodeMetadata.ID, disruptionID, updateFields)
	if err != nil {
		return errHandler(err)
	}

	// update the node record to bump the version.
	err = s.nodeTable.Update(ctx, execer, nodeMetadata.ID, nodeMetadata.Version, tables.NodeUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	err = tx.Commit()
	if err != nil {
		return errHandler(err)
	}
	return nil
}
