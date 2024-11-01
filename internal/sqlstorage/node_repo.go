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
	nodeCapabilityTable  *tables.NodeCapabilityTable
	nodeDisruptionTable  *tables.NodeDisruptionTable
}

func nodeRecordToRow(record node.NodeRecord) tables.NodeRow {
	return tables.NodeRow{
		ID:                   record.Metadata.ID,
		Version:              record.Metadata.Version,
		Name:                 record.Name,
		State:                string(record.Status.State),
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
			State:   node.NodeState(row.State),
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

func nodeDisruptionRecordToRow(nodeID string, record node.Disruption) tables.NodeDisruptionRow {
	startTime := uint64(0)
	if !record.StartTime.IsZero() {
		startTime = uint64(record.StartTime.Unix())
	}
	return tables.NodeDisruptionRow{
		ID:        record.ID,
		NodeID:    nodeID,
		StartTime: startTime,
		EvictNode: record.ShouldEvict,
		State:     string(record.Status.State),
		Message:   record.Status.Message,
	}
}

func nodeDisruptionRowToRecord(row tables.NodeDisruptionRow) node.Disruption {
	return node.Disruption{
		ID: row.ID,
		Status: node.DisruptionStatus{
			State:   node.DisruptionState(row.State),
			Message: row.Message,
		},
		StartTime:   time.Unix(int64(row.StartTime), 0), // Convert from uint64 to time.Time
		ShouldEvict: row.EvictNode,
	}
}

func nodeLocalVolumeRecordToRow(nodeID string, record node.LocalVolume) tables.NodeLocalVolumeRow {
	return tables.NodeLocalVolumeRow{
		NodeID:          nodeID,
		MountPath:       record.MountPath,
		StorageClass:    record.StorageClass,
		StorageCapacity: record.StorageCapacity,
	}
}

func nodeLocalVolumeRowToRecord(row tables.NodeLocalVolumeRow) node.LocalVolume {
	return node.LocalVolume{
		MountPath:       row.MountPath,
		StorageClass:    row.StorageClass,
		StorageCapacity: row.StorageCapacity,
	}
}

func nodeCapabilityRecordToRow(nodeID string, id string) tables.NodeCapabilityRow {
	return tables.NodeCapabilityRow{
		NodeID:       nodeID,
		CapabilityID: id,
	}
}

func nodeCapabilityRowToRecord(row tables.NodeCapabilityRow) string {
	return row.CapabilityID
}

// newNodeStorage creates a new storage instance satisfying the NodeRepository interface
func newNodeStorage(db simplesql.Database) node.Repository {
	return &nodeStorage{
		Database:             db,
		nodeTable:            tables.NewNodeTable(db),
		nodeLocalVolumeTable: tables.NewNodeLocalVolumeTable(db),
		nodeCapabilityTable:  tables.NewNodeCapabilityTable(db),
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

	for _, capabilityID := range record.CapabilityIDs {
		err = s.nodeCapabilityTable.Insert(ctx, tx, nodeCapabilityRecordToRow(record.Metadata.ID, capabilityID))
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

func (s *nodeStorage) GetByID(ctx context.Context, id string) (node.NodeRecord, error) {
	nodeRow, err := s.nodeTable.Get(ctx, tables.NodeKeys{
		ID: &id,
	})
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	nodeRecord := nodeRowToRecord(nodeRow)
	return s.getByPartialRecord(ctx, nodeRecord)
}

func (s *nodeStorage) GetByName(ctx context.Context, nodeName string) (node.NodeRecord, error) {
	nodeRow, err := s.nodeTable.Get(ctx, tables.NodeKeys{
		Name: &nodeName,
	})
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	nodeRecord := nodeRowToRecord(nodeRow)
	return s.getByPartialRecord(ctx, nodeRecord)
}

func (s *nodeStorage) getByPartialRecord(ctx context.Context, record node.NodeRecord) (node.NodeRecord, error) {
	// Get the corresponding disruptionRows
	disruptionRows, err := s.nodeDisruptionTable.List(ctx, tables.NodeDisruptionTableSelectFilters{
		NodeIDIn: []string{record.Metadata.ID},
	})
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	for _, disruption := range disruptionRows {
		record.Disruptions = append(record.Disruptions, nodeDisruptionRowToRecord(disruption))
	}

	// Get the corresponding volumeRows
	volumeRows, err := s.nodeLocalVolumeTable.List(ctx, tables.NodeLocalVolumeTableSelectFilters{
		NodeIDIn: []string{record.Metadata.ID},
	})
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	for _, volume := range volumeRows {
		record.LocalVolumes = append(record.LocalVolumes, nodeLocalVolumeRowToRecord(volume))
	}

	// Get the corresponding capabilityRows
	capabilityRows, err := s.nodeCapabilityTable.List(ctx, tables.NodeCapabilityTableSelectFilters{
		NodeIDIn: []string{record.Metadata.ID},
	})
	if err != nil {
		return node.NodeRecord{}, errHandler(err)
	}
	for _, capability := range capabilityRows {
		record.CapabilityIDs = append(record.CapabilityIDs, nodeCapabilityRowToRecord(capability))
	}

	return record, nil
}

func (s *nodeStorage) UpdateStatus(ctx context.Context, metadata core.Metadata, status node.NodeStatus, clusterID string) error {
	execer := s.DB
	state := string(status.State)
	message := status.Message
	updateFields := tables.NodeUpdateFields{
		State:     &state,
		Message:   &message,
		ClusterID: &clusterID,
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
		PayloadNameIn:      append([]string{}, filters.PayloadNameIn...),
		PayloadNameNotIn:   append([]string{}, filters.PayloadNameNotIn...),
		UpdateDomainIn:     append([]string{}, filters.UpdateDomainIn...),
	}

	// Extract node specific filters
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, string(state))
	}

	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, string(state))
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
	disruptionMap := make(map[string][]node.Disruption)
	for _, disruption := range disruptionRows {
		disruptionMap[disruption.NodeID] = append(disruptionMap[disruption.NodeID], nodeDisruptionRowToRecord(disruption))
	}

	volumeRows, err := s.nodeLocalVolumeTable.List(ctx, tables.NodeLocalVolumeTableSelectFilters{
		NodeIDIn: nodeIDs,
	})
	if err != nil {
		return nil, errHandler(err)
	}
	volumeMap := make(map[string][]node.LocalVolume)
	for _, volume := range volumeRows {
		volumeMap[volume.NodeID] = append(volumeMap[volume.NodeID], nodeLocalVolumeRowToRecord(volume))
	}

	// Get the corresponding capabilityRows
	capabilityRows, err := s.nodeCapabilityTable.List(ctx, tables.NodeCapabilityTableSelectFilters{
		NodeIDIn: nodeIDs,
	})
	if err != nil {
		return nil, errHandler(err)
	}
	capabilityMap := make(map[string][]string)
	for _, capability := range capabilityRows {
		capabilityMap[capability.NodeID] = append(capabilityMap[capability.NodeID], nodeCapabilityRowToRecord(capability))
	}

	var records []node.NodeRecord
	for _, row := range rows {
		record := nodeRowToRecord(row)
		record.Disruptions = disruptionMap[row.ID]
		record.LocalVolumes = volumeMap[row.ID]
		record.CapabilityIDs = capabilityMap[row.ID]
		records = append(records, record)
	}

	return records, nil
}

func (s *nodeStorage) InsertDisruption(ctx context.Context, nodeMetadata core.Metadata, record node.Disruption) error {
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

func (s *nodeStorage) UpdateDisruptionStatus(ctx context.Context, nodeMetadata core.Metadata, disruptionID string, status node.DisruptionStatus) error {
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

func (s *nodeStorage) InsertCapability(ctx context.Context, nodeMetadata core.Metadata, capabilityID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	err = s.nodeCapabilityTable.Insert(ctx, execer, nodeCapabilityRecordToRow(nodeMetadata.ID, capabilityID))
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

func (s *nodeStorage) DeleteCapability(ctx context.Context, nodeMetadata core.Metadata, capabilityID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	err = s.nodeCapabilityTable.Delete(ctx, execer, nodeMetadata.ID, capabilityID)
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
