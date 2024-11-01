package sqlstorage

import (
	"context"
	"encoding/json"

	"github.com/msanath/gondolf/pkg/simplesql"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/deploymentplan"
	"github.com/msanath/mrds/internal/sqlstorage/tables"
)

// deploymentPlanStorage is a concrete implementation of DeploymentPlanRepository using sqlx
type deploymentPlanStorage struct {
	simplesql.Database
	deploymentPlanTable                             *tables.DeploymentPlanTable
	deploymentPlanApplicationTable                  *tables.DeploymentPlanApplicationTable
	deploymentPlanApplicationPortTable              *tables.DeploymentPlanApplicationPortTable
	deploymentPlanApplicationPersistentVolumeTable  *tables.DeploymentPlanApplicationPersistentVolumeTable
	deploymentPlanDeploymentTable                   *tables.DeploymentPlanDeploymentTable
	deploymentPlanDeploymentPayloadCoordinatesTable *tables.DeploymentPlanDeploymentPayloadCoordinatesTable
	deploymentPlanMatchingCapabilityTable           *tables.DeploymentMatchingCapabilityTable
}

// newDeploymentPlanStorage creates a new storage instance satisfying the DeploymentPlanRepository interface
func newDeploymentPlanStorage(db simplesql.Database) deploymentplan.Repository {
	return &deploymentPlanStorage{
		Database:                                        db,
		deploymentPlanTable:                             tables.NewDeploymentPlanTable(db),
		deploymentPlanApplicationTable:                  tables.NewDeploymentPlanApplicationTable(db),
		deploymentPlanApplicationPortTable:              tables.NewDeploymentPlanApplicationPortTable(db),
		deploymentPlanApplicationPersistentVolumeTable:  tables.NewDeploymentPlanApplicationPersistentVolumeTable(db),
		deploymentPlanDeploymentTable:                   tables.NewDeploymentPlanDeploymentTable(db),
		deploymentPlanDeploymentPayloadCoordinatesTable: tables.NewDeploymentPlanDeploymentPayloadCoordinatesTable(db),
		deploymentPlanMatchingCapabilityTable:           tables.NewDeploymentMatchingCapabilityTable(db),
	}
}

func deploymentPlanRecordToRow(record deploymentplan.DeploymentPlanRecord) tables.DeploymentPlanRow {
	return tables.DeploymentPlanRow{
		ID:          record.Metadata.ID,
		Version:     record.Metadata.Version,
		Name:        record.Name,
		State:       string(record.Status.State),
		Message:     record.Status.Message,
		Namespace:   record.Namespace,
		ServiceName: record.ServiceName,
	}
}

func deploymentPlanRowToRecord(row tables.DeploymentPlanRow) deploymentplan.DeploymentPlanRecord {
	return deploymentplan.DeploymentPlanRecord{
		Metadata: core.Metadata{
			ID:      row.ID,
			Version: row.Version,
		},
		Name: row.Name,
		Status: deploymentplan.DeploymentPlanStatus{
			State:   deploymentplan.DeploymentPlanState(row.State),
			Message: row.Message,
		},
		Namespace:   row.Namespace,
		ServiceName: row.ServiceName,
	}
}

func deploymentApplicationRecordToRow(deploymentPlanID string, record deploymentplan.Application) tables.DeploymentPlanApplicationRow {
	return tables.DeploymentPlanApplicationRow{
		DeploymentPlanID: deploymentPlanID,
		PayloadName:      record.PayloadName,
		Cores:            record.Resources.Cores,
		Memory:           record.Resources.Memory,
	}
}

func deploymentApplicationRowToRecord(row tables.DeploymentPlanApplicationRow) deploymentplan.Application {
	return deploymentplan.Application{
		PayloadName: row.PayloadName,
		Resources: deploymentplan.ApplicationResources{
			Cores:  row.Cores,
			Memory: row.Memory,
		},
	}
}

func deploymentApplicationPortRecordToRow(deploymentPlanID string, payloadName string, record deploymentplan.ApplicationPort) tables.DeploymentPlanApplicationPortRow {
	return tables.DeploymentPlanApplicationPortRow{
		DeploymentPlanID: deploymentPlanID,
		PayloadName:      payloadName,
		Protocol:         record.Protocol,
		Port:             record.Port,
	}
}

func deploymentApplicationPortRowToRecord(row tables.DeploymentPlanApplicationPortRow) deploymentplan.ApplicationPort {
	return deploymentplan.ApplicationPort{
		Protocol: row.Protocol,
		Port:     row.Port,
	}
}

func deploymentApplicationPersistentVolumeRecordToRow(deploymentPlanID string, payloadName string, record deploymentplan.ApplicationPersistentVolume) tables.DeploymentPlanApplicationPersistentVolumeRow {
	return tables.DeploymentPlanApplicationPersistentVolumeRow{
		DeploymentPlanID: deploymentPlanID,
		PayloadName:      payloadName,
		StorageClass:     record.StorageClass,
		Capacity:         record.Capacity,
		MountPath:        record.MountPath,
	}
}

func deploymentApplicationPersistentVolumeRowToRecord(row tables.DeploymentPlanApplicationPersistentVolumeRow) deploymentplan.ApplicationPersistentVolume {
	return deploymentplan.ApplicationPersistentVolume{
		StorageClass: row.StorageClass,
		MountPath:    row.MountPath,
		Capacity:     row.Capacity,
	}
}

func deploymentRecordToRow(deploymentPlanID string, record deploymentplan.Deployment) tables.DeploymentPlanDeploymentRow {
	return tables.DeploymentPlanDeploymentRow{
		DeploymentPlanID: deploymentPlanID,
		ID:               record.ID,
		State:            string(record.Status.State),
		Message:          record.Status.Message,
		InstanceCount:    record.InstanceCount,
	}
}

func deploymentRowToRecord(row tables.DeploymentPlanDeploymentRow) deploymentplan.Deployment {
	return deploymentplan.Deployment{
		ID:            row.ID,
		Status:        deploymentplan.DeploymentStatus{State: deploymentplan.DeploymentState(row.State), Message: row.Message},
		InstanceCount: row.InstanceCount,
	}
}

func deploymentPlanDeploymentPayloadCoordinatesRecordToRow(
	deploymentPlanID string, deploymentID string, coordinates deploymentplan.PayloadCoordinates,
) tables.DeploymentPlanDeploymentPayloadCoordinatesRow {
	jsonStr, _ := json.Marshal(coordinates.Coordinates)

	return tables.DeploymentPlanDeploymentPayloadCoordinatesRow{
		DeploymentPlanID: deploymentPlanID,
		DeploymentID:     deploymentID,
		PayloadName:      coordinates.PayloadName,
		Coordinates:      string(jsonStr),
	}
}

func deploymentPlanDeploymentPayloadCoordinatesRowToRecord(row tables.DeploymentPlanDeploymentPayloadCoordinatesRow) deploymentplan.PayloadCoordinates {
	var coordinates map[string]string
	_ = json.Unmarshal([]byte(row.Coordinates), &coordinates)

	return deploymentplan.PayloadCoordinates{
		PayloadName: row.PayloadName,
		Coordinates: coordinates,
	}
}

func deploymentPlanMatchingCapabilityRecordToRow(deploymentPlanID string, record deploymentplan.MatchingComputeCapability) tables.DeploymentPlanMatchingCapabilityRow {
	jsonStr, _ := json.Marshal(record.CapabilityNames)

	return tables.DeploymentPlanMatchingCapabilityRow{
		DeploymentPlanID: deploymentPlanID,
		CapabilityType:   record.CapabilityType,
		Comparator:       string(record.Comparator),
		CapabilityNames:  string(jsonStr),
	}
}

func deploymentPlanMatchingCapabilityRowToRecord(row tables.DeploymentPlanMatchingCapabilityRow) deploymentplan.MatchingComputeCapability {
	var capabilityNames []string
	_ = json.Unmarshal([]byte(row.CapabilityNames), &capabilityNames)

	return deploymentplan.MatchingComputeCapability{
		CapabilityType:  row.CapabilityType,
		Comparator:      deploymentplan.ComparatorType(row.Comparator),
		CapabilityNames: capabilityNames,
	}
}

func (s *deploymentPlanStorage) Insert(ctx context.Context, record deploymentplan.DeploymentPlanRecord) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	err = s.deploymentPlanTable.Insert(ctx, execer, deploymentPlanRecordToRow(record))
	if err != nil {
		return errHandler(err)
	}

	for _, app := range record.Applications {
		err = s.deploymentPlanApplicationTable.Insert(ctx, execer, deploymentApplicationRecordToRow(record.Metadata.ID, app))
		if err != nil {
			return errHandler(err)
		}

		for _, port := range app.Ports {
			err = s.deploymentPlanApplicationPortTable.Insert(ctx, execer, deploymentApplicationPortRecordToRow(record.Metadata.ID, app.PayloadName, port))
			if err != nil {
				return errHandler(err)
			}
		}

		for _, pv := range app.PersistentVolumes {
			err = s.deploymentPlanApplicationPersistentVolumeTable.Insert(ctx, execer, deploymentApplicationPersistentVolumeRecordToRow(record.Metadata.ID, app.PayloadName, pv))
			if err != nil {
				return errHandler(err)
			}
		}
	}
	for _, capability := range record.MatchingComputeCapabilities {
		err = s.deploymentPlanMatchingCapabilityTable.Insert(ctx, execer, deploymentPlanMatchingCapabilityRecordToRow(record.Metadata.ID, capability))
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

func (s *deploymentPlanStorage) GetByID(ctx context.Context, id string) (deploymentplan.DeploymentPlanRecord, error) {
	row, err := s.deploymentPlanTable.Get(ctx, tables.DeploymentPlanKeys{
		ID: &id,
	})
	if err != nil {
		return deploymentplan.DeploymentPlanRecord{}, errHandler(err)
	}
	return s.getByPartialRecord(ctx, deploymentPlanRowToRecord(row))
}

func (s *deploymentPlanStorage) GetByName(ctx context.Context, name string) (deploymentplan.DeploymentPlanRecord, error) {
	row, err := s.deploymentPlanTable.Get(ctx, tables.DeploymentPlanKeys{
		Name: &name,
	})
	if err != nil {
		return deploymentplan.DeploymentPlanRecord{}, errHandler(err)
	}
	return s.getByPartialRecord(ctx, deploymentPlanRowToRecord(row))
}

func (s *deploymentPlanStorage) getByPartialRecord(ctx context.Context, deploymentPlanRecord deploymentplan.DeploymentPlanRecord) (deploymentplan.DeploymentPlanRecord, error) {
	record := deploymentPlanRecord
	applicationRows, err := s.deploymentPlanApplicationTable.List(ctx, tables.DeploymentPlanApplicationTableSelectFilters{
		DeploymentPlanIDIn: []string{record.Metadata.ID},
	})
	if err != nil {
		return deploymentplan.DeploymentPlanRecord{}, errHandler(err)
	}
	for _, row := range applicationRows {
		application := deploymentApplicationRowToRecord(row)

		portRows, err := s.deploymentPlanApplicationPortTable.List(ctx, tables.DeploymentPlanApplicationPortTableSelectFilters{
			DeploymentPlanIDIn: []string{record.Metadata.ID},
			PayloadNameIn:      []string{application.PayloadName},
		})
		if err != nil {
			return deploymentplan.DeploymentPlanRecord{}, errHandler(err)
		}
		for _, portRow := range portRows {
			application.Ports = append(application.Ports, deploymentApplicationPortRowToRecord(portRow))
		}

		pvRows, err := s.deploymentPlanApplicationPersistentVolumeTable.List(ctx, tables.DeploymentPlanApplicationPersistentVolumeTableSelectFilters{
			DeploymentPlanIDIn: []string{record.Metadata.ID},
			PayloadNameIn:      []string{application.PayloadName},
		})
		if err != nil {
			return deploymentplan.DeploymentPlanRecord{}, errHandler(err)
		}
		for _, pvRow := range pvRows {
			application.PersistentVolumes = append(application.PersistentVolumes, deploymentApplicationPersistentVolumeRowToRecord(pvRow))
		}

		record.Applications = append(record.Applications, application)
	}

	matchingCapabilityRows, err := s.deploymentPlanMatchingCapabilityTable.List(ctx, tables.DeploymentPlanMatchingCapabilityTableSelectFilters{
		DeploymentPlanIDIn: []string{record.Metadata.ID},
	})
	if err != nil {
		return deploymentplan.DeploymentPlanRecord{}, errHandler(err)
	}
	for _, row := range matchingCapabilityRows {
		record.MatchingComputeCapabilities = append(record.MatchingComputeCapabilities, deploymentPlanMatchingCapabilityRowToRecord(row))
	}

	deploymentRows, err := s.deploymentPlanDeploymentTable.List(ctx, tables.DeploymentPlanDeploymentTableSelectFilters{
		DeploymentPlanIDIn: []string{record.Metadata.ID},
	})
	if err != nil {
		return deploymentplan.DeploymentPlanRecord{}, errHandler(err)
	}
	for _, row := range deploymentRows {
		deployment := deploymentRowToRecord(row)
		rows, err := s.deploymentPlanDeploymentPayloadCoordinatesTable.List(ctx, tables.DeploymentPlanDeploymentPayloadCoordinatesTableSelectFilters{
			DeploymentPlanIDIn: []string{record.Metadata.ID},
			DeploymentIDIn:     []string{row.ID},
		})
		if err != nil {
			return deploymentplan.DeploymentPlanRecord{}, errHandler(err)
		}
		for _, row := range rows {
			deployment.PayloadCoordinates = append(deployment.PayloadCoordinates, deploymentPlanDeploymentPayloadCoordinatesRowToRecord(row))
		}
		record.Deployments = append(record.Deployments, deployment)
	}

	return record, nil
}

func (s *deploymentPlanStorage) UpdateStatus(ctx context.Context, metadata core.Metadata, status deploymentplan.DeploymentPlanStatus) error {
	execer := s.DB
	state := string(status.State)
	message := status.Message
	updateFields := tables.DeploymentPlanTableUpdateFields{
		State:   &state,
		Message: &message,
	}
	err := s.deploymentPlanTable.Update(ctx, execer, metadata.ID, metadata.Version, updateFields)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *deploymentPlanStorage) Delete(ctx context.Context, metadata core.Metadata) error {
	execer := s.DB
	err := s.deploymentPlanTable.Delete(ctx, execer, metadata.ID, metadata.Version)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *deploymentPlanStorage) List(ctx context.Context, filters deploymentplan.DeploymentPlanListFilters) ([]deploymentplan.DeploymentPlanRecord, error) {
	dbFilters := tables.DeploymentPlanTableSelectFilters{
		IDIn:           append([]string{}, filters.IDIn...),
		NameIn:         append([]string{}, filters.NameIn...),
		VersionGte:     filters.VersionGte,
		VersionLte:     filters.VersionLte,
		VersionEq:      filters.VersionEq,
		IncludeDeleted: filters.IncludeDeleted,
		Limit:          filters.Limit,
	}
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, string(state))
	}
	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, string(state))
	}

	rows, err := s.deploymentPlanTable.List(ctx, dbFilters)
	if err != nil {
		return nil, err
	}
	var records []deploymentplan.DeploymentPlanRecord
	for _, row := range rows {
		records = append(records, deploymentPlanRowToRecord(row))
	}
	deploymentPlanIDs := make([]string, len(records))
	for i, record := range records {
		deploymentPlanIDs[i] = record.Metadata.ID
	}

	applicationRows, err := s.deploymentPlanApplicationTable.List(ctx, tables.DeploymentPlanApplicationTableSelectFilters{
		DeploymentPlanIDIn: deploymentPlanIDs,
	})
	if err != nil {
		return nil, errHandler(err)
	}
	deploymentPlanIDToApplications := make(map[string][]deploymentplan.Application)
	for _, row := range applicationRows {
		application := deploymentApplicationRowToRecord(row)

		portRows, err := s.deploymentPlanApplicationPortTable.List(ctx, tables.DeploymentPlanApplicationPortTableSelectFilters{
			DeploymentPlanIDIn: []string{row.DeploymentPlanID},
			PayloadNameIn:      []string{application.PayloadName},
		})
		if err != nil {
			return nil, errHandler(err)
		}
		for _, portRow := range portRows {
			application.Ports = append(application.Ports, deploymentApplicationPortRowToRecord(portRow))
		}

		pvRows, err := s.deploymentPlanApplicationPersistentVolumeTable.List(ctx, tables.DeploymentPlanApplicationPersistentVolumeTableSelectFilters{
			DeploymentPlanIDIn: []string{row.DeploymentPlanID},
			PayloadNameIn:      []string{application.PayloadName},
		})
		if err != nil {
			return nil, errHandler(err)
		}
		for _, pvRow := range pvRows {
			application.PersistentVolumes = append(application.PersistentVolumes, deploymentApplicationPersistentVolumeRowToRecord(pvRow))
		}

		deploymentPlanIDToApplications[row.DeploymentPlanID] = append(deploymentPlanIDToApplications[row.DeploymentPlanID], deploymentApplicationRowToRecord(row))
	}
	for i, record := range records {
		records[i].Applications = deploymentPlanIDToApplications[record.Metadata.ID]
	}

	deploymentRows, err := s.deploymentPlanDeploymentTable.List(ctx, tables.DeploymentPlanDeploymentTableSelectFilters{
		DeploymentPlanIDIn: deploymentPlanIDs,
	})
	if err != nil {
		return nil, errHandler(err)
	}
	deploymentPlanIDToDeployments := make(map[string][]deploymentplan.Deployment)
	for _, row := range deploymentRows {
		deployment := deploymentRowToRecord(row)
		rows, err := s.deploymentPlanDeploymentPayloadCoordinatesTable.List(ctx, tables.DeploymentPlanDeploymentPayloadCoordinatesTableSelectFilters{
			DeploymentPlanIDIn: []string{row.DeploymentPlanID},
			DeploymentIDIn:     []string{row.ID},
		})
		if err != nil {
			return nil, errHandler(err)
		}
		for _, row := range rows {
			deployment.PayloadCoordinates = append(deployment.PayloadCoordinates, deploymentPlanDeploymentPayloadCoordinatesRowToRecord(row))
		}

		deploymentPlanIDToDeployments[row.DeploymentPlanID] = append(deploymentPlanIDToDeployments[row.DeploymentPlanID], deployment)
	}
	for i, record := range records {
		records[i].Deployments = deploymentPlanIDToDeployments[record.Metadata.ID]
	}

	matchingCapabilityRows, err := s.deploymentPlanMatchingCapabilityTable.List(ctx, tables.DeploymentPlanMatchingCapabilityTableSelectFilters{
		DeploymentPlanIDIn: deploymentPlanIDs,
	})
	if err != nil {
		return nil, errHandler(err)
	}
	deploymentPlanIDToMatchingCapabilities := make(map[string][]deploymentplan.MatchingComputeCapability)
	for _, row := range matchingCapabilityRows {
		deploymentPlanIDToMatchingCapabilities[row.DeploymentPlanID] = append(deploymentPlanIDToMatchingCapabilities[row.DeploymentPlanID], deploymentPlanMatchingCapabilityRowToRecord(row))
	}
	for i, record := range records {
		records[i].MatchingComputeCapabilities = deploymentPlanIDToMatchingCapabilities[record.Metadata.ID]
	}
	return records, nil
}

func (s *deploymentPlanStorage) InsertDeployment(ctx context.Context, metadata core.Metadata, deployment deploymentplan.Deployment) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	err = s.deploymentPlanDeploymentTable.Insert(ctx, execer, deploymentRecordToRow(metadata.ID, deployment))
	if err != nil {
		return err
	}

	for _, coordinates := range deployment.PayloadCoordinates {
		err = s.deploymentPlanDeploymentPayloadCoordinatesTable.Insert(
			ctx, execer, deploymentPlanDeploymentPayloadCoordinatesRecordToRow(metadata.ID, deployment.ID, coordinates))
		if err != nil {
			return err
		}
	}

	// bump the version of the deployment plan
	err = s.deploymentPlanTable.Update(ctx, execer, metadata.ID, metadata.Version, tables.DeploymentPlanTableUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	return tx.Commit()
}

func (s *deploymentPlanStorage) UpdateDeploymentStatus(ctx context.Context, metadata core.Metadata, deploymentID string, status deploymentplan.DeploymentStatus) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errHandler(err)
	}
	defer tx.Rollback()

	execer := tx
	state := string(status.State)
	message := status.Message
	err = s.deploymentPlanDeploymentTable.Update(ctx, execer, deploymentID, metadata.ID, tables.DeploymentPlanDeploymentTableUpdateFields{
		State:   &state,
		Message: &message,
	})
	if err != nil {
		return errHandler(err)
	}

	// bump the version of the deployment plan
	err = s.deploymentPlanTable.Update(ctx, execer, metadata.ID, metadata.Version, tables.DeploymentPlanTableUpdateFields{})
	if err != nil {
		return errHandler(err)
	}

	return tx.Commit()
}
