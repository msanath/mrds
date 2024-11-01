package sqlstorage

import (
	"context"
	"fmt"

	"github.com/msanath/gondolf/pkg/simplesql"
	"github.com/msanath/mrds/internal/ledger/core"
	ledgererrors "github.com/msanath/mrds/internal/ledger/errors"
	"github.com/msanath/mrds/internal/ledger/metainstance"
	"github.com/msanath/mrds/internal/sqlstorage/tables"
)

// metaInstanceStorage is a concrete implementation of MetaInstanceRepository using sqlx
type metaInstanceStorage struct {
	simplesql.Database
	metaInstanceTable                *tables.MetaInstanceTable
	metaInstanceOperationTable       *tables.MetaInstanceOperationTable
	metaInstanceRuntimeInstanceTable *tables.MetaInstanceRuntimeInstanceTable
	deploymentPlanApplicationTable   *tables.DeploymentPlanApplicationTable
	nodeTable                        *tables.NodeTable
	nodePayloadTable                 *tables.NodePayloadTable
}

// newMetaInstanceStorage creates a new storage instance satisfying the MetaInstanceRepository interface
func newMetaInstanceStorage(db simplesql.Database) metainstance.Repository {
	return &metaInstanceStorage{
		Database:                         db,
		metaInstanceTable:                tables.NewMetaInstanceTable(db),
		metaInstanceOperationTable:       tables.NewMetaInstanceOperationTable(db),
		metaInstanceRuntimeInstanceTable: tables.NewMetaInstanceRuntimeInstanceTable(db),
		deploymentPlanApplicationTable:   tables.NewDeploymentPlanApplicationTable(db),
		nodeTable:                        tables.NewNodeTable(db),
		nodePayloadTable:                 tables.NewNodePayloadTable(db),
	}
}

func metaInstanceRecordToRow(record metainstance.MetaInstanceRecord) tables.MetaInstanceRow {
	return tables.MetaInstanceRow{
		ID:               record.Metadata.ID,
		Version:          record.Metadata.Version,
		Name:             record.Name,
		State:            string(record.Status.State),
		Message:          record.Status.Message,
		DeploymentPlanID: record.DeploymentPlanID,
		DeploymentID:     record.DeploymentID,
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
			State:   metainstance.MetaInstanceState(row.State),
			Message: row.Message,
		},
		DeploymentPlanID: row.DeploymentPlanID,
		DeploymentID:     row.DeploymentID,
	}
}

func metaInstanceOperationRecordToRow(metaInstanceID string, record metainstance.Operation) tables.MetaInstanceOperationRow {
	return tables.MetaInstanceOperationRow{
		ID:             record.ID,
		MetaInstanceID: metaInstanceID,
		Type:           record.Type,
		IntentID:       record.IntentID,
		State:          string(record.Status.State),
		Message:        record.Status.Message,
	}
}

func metaInstanceOperationRowToModel(row tables.MetaInstanceOperationRow) metainstance.Operation {
	return metainstance.Operation{
		ID:       row.ID,
		Type:     row.Type,
		IntentID: row.IntentID,
		Status: metainstance.OperationStatus{
			State:   metainstance.OperationState(row.State),
			Message: row.Message,
		},
	}
}

func metaInstanceRuntimeInstanceRecordToRow(metaInstanceID string, record metainstance.RuntimeInstance) tables.MetaInstanceRuntimeInstanceRow {
	return tables.MetaInstanceRuntimeInstanceRow{
		ID:             record.ID,
		MetaInstanceID: metaInstanceID,
		NodeID:         record.NodeID,
		IsActive:       record.IsActive,
		State:          string(record.Status.State),
		Message:        record.Status.Message,
	}
}

func metaInstanceRuntimeInstanceRowToModel(row tables.MetaInstanceRuntimeInstanceRow) metainstance.RuntimeInstance {
	return metainstance.RuntimeInstance{
		ID:       row.ID,
		NodeID:   row.NodeID,
		IsActive: row.IsActive,
		Status: metainstance.RuntimeInstanceStatus{
			State:   metainstance.RuntimeInstanceState(row.State),
			Message: row.Message,
		},
	}
}

func (s *metaInstanceStorage) Insert(ctx context.Context, record metainstance.MetaInstanceRecord) error {
	execer := s.DB
	err := s.metaInstanceTable.Insert(ctx, execer, metaInstanceRecordToRow(record))
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) GetByID(ctx context.Context, id string) (metainstance.MetaInstanceRecord, error) {
	row, err := s.metaInstanceTable.Get(ctx, tables.MetaInstanceKeys{
		ID: &id,
	})
	if err != nil {
		return metainstance.MetaInstanceRecord{}, errHandler(err)
	}
	return s.getByPartialRecord(ctx, metaInstanceRowToModel(row))
}

func (s *metaInstanceStorage) GetByName(ctx context.Context, name string) (metainstance.MetaInstanceRecord, error) {
	row, err := s.metaInstanceTable.Get(ctx, tables.MetaInstanceKeys{
		Name: &name,
	})
	if err != nil {
		return metainstance.MetaInstanceRecord{}, errHandler(err)
	}
	return s.getByPartialRecord(ctx, metaInstanceRowToModel(row))
}

func (s *metaInstanceStorage) getByPartialRecord(ctx context.Context, metaInstanceRecord metainstance.MetaInstanceRecord) (metainstance.MetaInstanceRecord, error) {
	record := metaInstanceRecord

	operationRows, err := s.metaInstanceOperationTable.List(ctx, tables.MetaInstanceOperationTableSelectFilters{
		MetaInstanceIDIn: []string{record.Metadata.ID},
	})
	if err != nil {
		return record, errHandler(err)
	}
	for _, row := range operationRows {
		record.Operations = append(record.Operations, metaInstanceOperationRowToModel(row))
	}

	runtimeInstanceRows, err := s.metaInstanceRuntimeInstanceTable.List(ctx, tables.MetaInstanceRuntimeInstanceTableSelectFilters{
		MetaInstanceIDIn: []string{record.Metadata.ID},
	})
	if err != nil {
		return record, errHandler(err)
	}
	for _, row := range runtimeInstanceRows {
		record.RuntimeInstances = append(record.RuntimeInstances, metaInstanceRuntimeInstanceRowToModel(row))
	}

	return record, nil
}

func (s *metaInstanceStorage) UpdateStatus(ctx context.Context, metadata core.Metadata, status metainstance.MetaInstanceStatus) error {
	execer := s.DB
	state := string(status.State)
	message := status.Message
	updateFields := tables.MetaInstanceTableUpdateFields{
		State:   &state,
		Message: &message,
	}
	err := s.metaInstanceTable.Update(ctx, execer, metadata.ID, metadata.Version, updateFields)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) UpdateDeploymentID(ctx context.Context, metadata core.Metadata, deploymentID string) error {
	execer := s.DB
	updateFields := tables.MetaInstanceTableUpdateFields{
		DeploymentID: &deploymentID,
	}
	err := s.metaInstanceTable.Update(ctx, execer, metadata.ID, metadata.Version, updateFields)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	execer := s.DB
	err := s.metaInstanceTable.Delete(ctx, execer, metadata.ID, metadata.Version)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) List(ctx context.Context, filters metainstance.MetaInstanceListFilters) ([]metainstance.MetaInstanceRecord, error) {
	dbFilters := tables.MetaInstanceTableSelectFilters{
		IDIn:               append([]string{}, filters.IDIn...),
		NameIn:             append([]string{}, filters.NameIn...),
		VersionGte:         filters.VersionGte,
		VersionLte:         filters.VersionLte,
		VersionEq:          filters.VersionEq,
		DeploymentIDIn:     append([]string{}, filters.DeploymentIDIn...),
		DeploymentPlanIDIn: append([]string{}, filters.DeploymentPlanIDIn...),
		IncludeDeleted:     filters.IncludeDeleted,
		Limit:              filters.Limit,
	}
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, string(state))
	}
	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, string(state))
	}

	rows, err := s.metaInstanceTable.List(ctx, dbFilters)
	if err != nil {
		return nil, err
	}
	metaInstanceIDs := make([]string, len(rows))
	for i, row := range rows {
		metaInstanceIDs[i] = row.ID
	}

	operationRows, err := s.metaInstanceOperationTable.List(ctx, tables.MetaInstanceOperationTableSelectFilters{
		MetaInstanceIDIn: metaInstanceIDs,
	})
	if err != nil {
		return nil, err
	}
	operationMap := make(map[string][]metainstance.Operation)
	for _, row := range operationRows {
		operationMap[row.MetaInstanceID] = append(operationMap[row.MetaInstanceID], metaInstanceOperationRowToModel(row))
	}

	runtimeInstanceRows, err := s.metaInstanceRuntimeInstanceTable.List(ctx, tables.MetaInstanceRuntimeInstanceTableSelectFilters{
		MetaInstanceIDIn: metaInstanceIDs,
	})
	if err != nil {
		return nil, err
	}
	runtimeInstanceMap := make(map[string][]metainstance.RuntimeInstance)
	for _, row := range runtimeInstanceRows {
		runtimeInstanceMap[row.MetaInstanceID] = append(runtimeInstanceMap[row.MetaInstanceID], metaInstanceRuntimeInstanceRowToModel(row))
	}

	var records []metainstance.MetaInstanceRecord
	for _, row := range rows {
		record := metaInstanceRowToModel(row)
		record.Operations = operationMap[row.ID]
		record.RuntimeInstances = runtimeInstanceMap[row.ID]
		records = append(records, record)
	}
	return records, nil
}

func (s *metaInstanceStorage) InsertOperation(ctx context.Context, metadata core.Metadata, operation metainstance.Operation) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	err = s.metaInstanceOperationTable.Insert(ctx, execer, metaInstanceOperationRecordToRow(metadata.ID, operation))
	if err != nil {
		return errHandler(err)
	}

	// update the meta instance state version
	err = s.metaInstanceTable.Update(ctx, execer, metadata.ID, metadata.Version, tables.MetaInstanceTableUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	err = tx.Commit()
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) UpdateOperationStatus(ctx context.Context, metadata core.Metadata, operationID string, status metainstance.OperationStatus) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx

	state := string(status.State)
	message := status.Message
	updateFields := tables.MetaInstanceOperationTableUpdateFields{
		State:   &state,
		Message: &message,
	}
	err = s.metaInstanceOperationTable.Update(ctx, execer, operationID, metadata.ID, updateFields)
	if err != nil {
		return errHandler(err)
	}

	// update the meta instance state version
	err = s.metaInstanceTable.Update(ctx, execer, metadata.ID, metadata.Version, tables.MetaInstanceTableUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	err = tx.Commit()
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) DeleteOperation(ctx context.Context, metadata core.Metadata, operationID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	err = s.metaInstanceOperationTable.Delete(ctx, execer, operationID, metadata.ID)
	if err != nil {
		return errHandler(err)
	}

	// update the meta instance state version
	err = s.metaInstanceTable.Update(ctx, execer, metadata.ID, metadata.Version, tables.MetaInstanceTableUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	err = tx.Commit()
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) InsertRuntimeInstance(ctx context.Context, metadata core.Metadata, runtimeInstance metainstance.RuntimeInstance) error {
	// Get the associated metaInstance Record.
	// Now get the sum of all the cores and memory for all applications in the deployment plan.
	// This will be used to check if the node has enough resources to run the application.
	// TODO: Optimize this
	metaInstanceRow, err := s.metaInstanceTable.Get(ctx, tables.MetaInstanceKeys{
		ID: &metadata.ID,
	})
	if err != nil {
		return errHandler(err)
	}
	applicationRows, err := s.deploymentPlanApplicationTable.List(ctx, tables.DeploymentPlanApplicationTableSelectFilters{
		DeploymentPlanIDIn: []string{metaInstanceRow.DeploymentPlanID},
	})
	if err != nil {
		return errHandler(err)
	}
	requestedCores := uint32(0)
	requestedMemory := uint32(0)
	payloadNames := []string{}
	for _, app := range applicationRows {
		requestedCores += app.Cores
		requestedMemory += app.Memory
		payloadNames = append(payloadNames, app.PayloadName)
	}

	nodeRow, err := s.nodeTable.Get(ctx, tables.NodeKeys{
		ID: &runtimeInstance.NodeID,
	})
	if err != nil {
		return errHandler(err)
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	err = s.metaInstanceRuntimeInstanceTable.Insert(ctx, execer, metaInstanceRuntimeInstanceRecordToRow(metadata.ID, runtimeInstance))
	if err != nil {
		return errHandler(err)
	}

	if nodeRow.RemainingCores < requestedCores {
		return ledgererrors.NewLedgerError(
			ledgererrors.ErrRecordInsertConflict,
			fmt.Sprintf("Node does not have enough cores to run the application. Requested: %d, Available: %d", requestedCores, nodeRow.RemainingCores),
		)
	}

	if nodeRow.RemainingMemory < requestedMemory {
		return ledgererrors.NewLedgerError(
			ledgererrors.ErrRecordInsertConflict,
			fmt.Sprintf("Node does not have enough memory to run the application. Requested: %d, Available: %d", requestedMemory, nodeRow.RemainingMemory),
		)
	}
	newRemainingCores := nodeRow.RemainingCores - requestedCores
	newRemainingMemory := nodeRow.RemainingMemory - requestedMemory

	updateFields := tables.NodeUpdateFields{
		RemainingCores:  &newRemainingCores,
		RemainingMemory: &newRemainingMemory,
	}
	err = s.nodeTable.Update(ctx, execer, nodeRow.ID, nodeRow.Version, updateFields)
	if err != nil {
		return errHandler(err)
	}
	for _, payloadName := range payloadNames {
		err = s.nodePayloadTable.Insert(ctx, execer, tables.NodePayloadRow{
			NodeID:      runtimeInstance.NodeID,
			PayloadName: payloadName,
		})
		if err != nil {
			return errHandler(err)
		}
	}

	// update the meta instance state version
	err = s.metaInstanceTable.Update(ctx, execer, metadata.ID, metadata.Version, tables.MetaInstanceTableUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	err = tx.Commit()
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) UpdateRuntimeInstanceStatus(ctx context.Context, metadata core.Metadata, runtimeInstanceID string, status metainstance.RuntimeInstanceStatus) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx

	state := string(status.State)
	message := status.Message
	updateFields := tables.MetaInstanceRuntimeInstanceTableUpdateFields{
		State:   &state,
		Message: &message,
	}
	err = s.metaInstanceRuntimeInstanceTable.Update(ctx, execer, runtimeInstanceID, metadata.ID, updateFields)
	if err != nil {
		return errHandler(err)
	}

	// update the meta instance state version
	err = s.metaInstanceTable.Update(ctx, execer, metadata.ID, metadata.Version, tables.MetaInstanceTableUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	err = tx.Commit()
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *metaInstanceStorage) DeleteRuntimeInstance(ctx context.Context, metadata core.Metadata, runtimeInstanceID string) error {
	// Get the associated metaInstance Record.
	// Now get the sum of all the cores and memory for all applications in the deployment plan.
	// This will be used to check if the node has enough resources to run the application.
	// TODO: Optimize this
	metaInstanceRow, err := s.metaInstanceTable.Get(ctx, tables.MetaInstanceKeys{
		ID: &metadata.ID,
	})
	if err != nil {
		return errHandler(err)
	}
	applicationRows, err := s.deploymentPlanApplicationTable.List(ctx, tables.DeploymentPlanApplicationTableSelectFilters{
		DeploymentPlanIDIn: []string{metaInstanceRow.DeploymentPlanID},
	})
	if err != nil {
		return errHandler(err)
	}
	requestedCores := uint32(0)
	requestedMemory := uint32(0)
	payloadNames := []string{}
	for _, app := range applicationRows {
		requestedCores += app.Cores
		requestedMemory += app.Memory
		payloadNames = append(payloadNames, app.PayloadName)
	}

	// Get the corresponding runtimee instance
	runtimeInstanceRows, err := s.metaInstanceRuntimeInstanceTable.List(ctx, tables.MetaInstanceRuntimeInstanceTableSelectFilters{
		MetaInstanceIDIn: []string{metaInstanceRow.ID},
	})
	if err != nil {
		return errHandler(err)
	}
	found := false
	var runtimeInstanceRow tables.MetaInstanceRuntimeInstanceRow
	for _, row := range runtimeInstanceRows {
		if row.ID == runtimeInstanceID {
			found = true
			runtimeInstanceRow = row
			break
		}
	}
	if !found {
		return ledgererrors.NewLedgerError(
			ledgererrors.ErrRecordNotFound,
			"Runtime instance not found.",
		)
	}
	nodeRow, err := s.nodeTable.Get(ctx, tables.NodeKeys{
		ID: &runtimeInstanceRow.NodeID,
	})
	if err != nil {
		return errHandler(err)
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	err = s.metaInstanceRuntimeInstanceTable.Delete(ctx, execer, runtimeInstanceID, metadata.ID)
	if err != nil {
		return errHandler(err)
	}

	// Add the resources back to the node
	newRemainingCores := nodeRow.RemainingCores + requestedCores
	newRemainingMemory := nodeRow.RemainingMemory + requestedMemory

	updateFields := tables.NodeUpdateFields{
		RemainingCores:  &newRemainingCores,
		RemainingMemory: &newRemainingMemory,
	}
	err = s.nodeTable.Update(ctx, execer, nodeRow.ID, nodeRow.Version, updateFields)
	if err != nil {
		return errHandler(err)
	}

	for _, payloadName := range payloadNames {
		err = s.nodePayloadTable.Delete(ctx, execer, nodeRow.ID, payloadName)
		if err != nil {
			return errHandler(err)
		}
	}

	// update the meta instance state version
	err = s.metaInstanceTable.Update(ctx, execer, metadata.ID, metadata.Version, tables.MetaInstanceTableUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	err = tx.Commit()
	if err != nil {
		return errHandler(err)
	}
	return nil
}
