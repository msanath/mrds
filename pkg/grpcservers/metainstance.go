package grpcservers

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/metainstance"
)

type MetaInstanceService struct {
	ledger              metainstance.Ledger
	ledgerRecordToProto func(record metainstance.MetaInstanceRecord) *mrdspb.MetaInstance

	mrdspb.UnimplementedMetaInstancesServer
}

func metaInstanceLedgerRecordToProto(record metainstance.MetaInstanceRecord) *mrdspb.MetaInstance {
	metaInstance := &mrdspb.MetaInstance{
		Metadata: &mrdspb.Metadata{
			Id:      record.Metadata.ID,
			Version: record.Metadata.Version,
		},
		Name: record.Name,
		Status: &mrdspb.MetaInstanceStatus{
			State:   mrdspb.MetaInstanceState(mrdspb.MetaInstanceState_value[string(record.Status.State)]),
			Message: record.Status.Message,
		},
		DeploymentPlanId: record.DeploymentPlanID,
		DeploymentId:     record.DeploymentID,
	}

	for _, runtimeInstance := range record.RuntimeInstances {
		metaInstance.RuntimeInstances = append(metaInstance.RuntimeInstances, &mrdspb.RuntimeInstance{
			Id:       runtimeInstance.ID,
			NodeId:   runtimeInstance.NodeID,
			IsActive: runtimeInstance.IsActive,
			Status: &mrdspb.RuntimeInstanceStatus{
				State:   mrdspb.RuntimeInstanceState(mrdspb.RuntimeInstanceState_value[string(runtimeInstance.Status.State)]),
				Message: runtimeInstance.Status.Message,
			},
		})
	}

	for _, operation := range record.Operations {
		metaInstance.Operations = append(metaInstance.Operations, &mrdspb.Operation{
			Id:       operation.ID,
			Type:     operation.Type,
			IntentId: operation.IntentID,
			Status: &mrdspb.OperationStatus{
				State:   mrdspb.OperationState(mrdspb.OperationState_value[string(operation.Status.State)]),
				Message: operation.Status.Message,
			},
		})
	}

	return metaInstance
}

func NewMetaInstanceService(ledger metainstance.Ledger) *MetaInstanceService {
	return &MetaInstanceService{
		ledger:              ledger,
		ledgerRecordToProto: metaInstanceLedgerRecordToProto,
	}
}

// Create creates a new MetaInstance
func (s *MetaInstanceService) Create(ctx context.Context, req *mrdspb.CreateMetaInstanceRequest) (*mrdspb.CreateMetaInstanceResponse, error) {
	cr := &metainstance.CreateRequest{
		Name:             req.Name,
		DeploymentPlanID: req.DeploymentPlanId,
		DeploymentID:     req.DeploymentId,
	}

	createResponse, err := s.ledger.Create(ctx, cr)
	if err != nil {
		return nil, err
	}

	return &mrdspb.CreateMetaInstanceResponse{Record: s.ledgerRecordToProto(createResponse.Record)}, nil
}

// GetByMetadata retrieves a MetaInstance by its metadata
func (s *MetaInstanceService) GetByMetadata(ctx context.Context, req *mrdspb.GetMetaInstanceByMetadataRequest) (*mrdspb.GetMetaInstanceResponse, error) {
	getResponse, err := s.ledger.GetByMetadata(ctx, &core.Metadata{
		ID:      req.Metadata.Id,
		Version: req.Metadata.Version,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetMetaInstanceResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// GetByName retrieves a MetaInstance by its name
func (s *MetaInstanceService) GetByName(ctx context.Context, req *mrdspb.GetMetaInstanceByNameRequest) (*mrdspb.GetMetaInstanceResponse, error) {
	getResponse, err := s.ledger.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetMetaInstanceResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// UpdateStatus updates the state and message of an existing MetaInstance
func (s *MetaInstanceService) UpdateStatus(ctx context.Context, req *mrdspb.UpdateMetaInstanceStatusRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	updateResponse, err := s.ledger.UpdateStatus(ctx, &metainstance.UpdateStatusRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		Status: metainstance.MetaInstanceStatus{
			State:   metainstance.MetaInstanceState(req.Status.State.String()),
			Message: req.Status.Message,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateMetaInstanceResponse{Record: s.ledgerRecordToProto(updateResponse.Record)}, nil
}

// UpdateDeploymentID updates the DeploymentID of an existing MetaInstance
func (s *MetaInstanceService) UpdateDeploymentID(ctx context.Context, req *mrdspb.UpdateDeploymentIDRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	updateResponse, err := s.ledger.UpdateDeploymentID(ctx, &metainstance.UpdateDeploymentIDRequest{
		Metadata:     core.Metadata{ID: req.Metadata.Id, Version: req.Metadata.Version},
		DeploymentID: req.DeploymentId,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateMetaInstanceResponse{Record: s.ledgerRecordToProto(updateResponse.Record)}, nil
}

// List returns a list of MetaInstances that match the provided filters
func (s *MetaInstanceService) List(ctx context.Context, req *mrdspb.ListMetaInstanceRequest) (*mrdspb.ListMetaInstanceResponse, error) {
	if req == nil {
		req = &mrdspb.ListMetaInstanceRequest{}
	}
	var gte, lte, eq *uint64
	if req.VersionGte != 0 {
		gte = &req.VersionGte
	}
	if req.VersionLte != 0 {
		lte = &req.VersionLte
	}
	if req.VersionEq != 0 {
		eq = &req.VersionEq
	}

	stateIn := make([]metainstance.MetaInstanceState, len(req.StateIn))
	for i, state := range req.StateIn {
		stateIn[i] = metainstance.MetaInstanceState(state.String())
	}

	listResponse, err := s.ledger.List(ctx, &metainstance.ListRequest{
		Filters: metainstance.MetaInstanceListFilters{
			IDIn:               req.IdIn,
			NameIn:             req.NameIn,
			VersionGte:         gte,
			VersionLte:         lte,
			VersionEq:          eq,
			StateIn:            stateIn,
			DeploymentIDIn:     req.DeploymentIdIn,
			DeploymentPlanIDIn: req.DeploymentPlanIdIn,
		},
	})
	if err != nil {
		return nil, err
	}

	records := make([]*mrdspb.MetaInstance, len(listResponse.Records))
	for i, record := range listResponse.Records {
		records[i] = s.ledgerRecordToProto(record)
	}

	return &mrdspb.ListMetaInstanceResponse{Records: records}, nil
}

// Delete deletes a MetaInstance
func (s *MetaInstanceService) Delete(ctx context.Context, req *mrdspb.DeleteMetaInstanceRequest) (*mrdspb.DeleteMetaInstanceResponse, error) {
	err := s.ledger.Delete(ctx, &metainstance.DeleteRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.DeleteMetaInstanceResponse{}, nil
}

// AddRuntimeInstance adds a runtime instance to a MetaInstance
func (s *MetaInstanceService) AddRuntimeInstance(ctx context.Context, req *mrdspb.AddRuntimeInstanceRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	addRuntimeResponse, err := s.ledger.AddRuntimeInstance(ctx, &metainstance.AddRuntimeInstanceRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		RuntimeInstance: metainstance.RuntimeInstance{
			ID:     req.RuntimeInstance.Id,
			NodeID: req.RuntimeInstance.NodeId,
			Status: metainstance.RuntimeInstanceStatus{
				State:   metainstance.RuntimeInstanceState(req.RuntimeInstance.Status.State.String()),
				Message: req.RuntimeInstance.Status.Message,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateMetaInstanceResponse{Record: s.ledgerRecordToProto(addRuntimeResponse.Record)}, nil
}

// UpdateRuntimeStatus updates the status of a runtime instance on a MetaInstance
func (s *MetaInstanceService) UpdateRuntimeStatus(ctx context.Context, req *mrdspb.UpdateRuntimeStatusRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	updateRuntimeStatusResponse, err := s.ledger.UpdateRuntimeStatus(ctx, &metainstance.UpdateRuntimeStatusRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		RuntimeInstanceID: req.RuntimeInstanceId,
		Status: metainstance.RuntimeInstanceStatus{
			State:   metainstance.RuntimeInstanceState(req.Status.State.String()),
			Message: req.Status.Message,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateMetaInstanceResponse{Record: s.ledgerRecordToProto(updateRuntimeStatusResponse.Record)}, nil
}

// RemoveRuntimeInstance removes a runtime instance from a MetaInstance
func (s *MetaInstanceService) RemoveRuntimeInstance(ctx context.Context, req *mrdspb.RemoveRuntimeInstanceRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	removeRuntimeResponse, err := s.ledger.RemoveRuntimeInstance(ctx, &metainstance.RemoveRuntimeInstanceRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		RuntimeInstanceID: req.RuntimeInstanceId,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateMetaInstanceResponse{Record: s.ledgerRecordToProto(removeRuntimeResponse.Record)}, nil
}

// AddOperation adds an operation to a MetaInstance
func (s *MetaInstanceService) AddOperation(ctx context.Context, req *mrdspb.AddOperationRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	addOperationResponse, err := s.ledger.AddOperation(ctx, &metainstance.AddOperationRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		Operation: metainstance.Operation{
			ID:       req.Operation.Id,
			Type:     req.Operation.Type,
			IntentID: req.Operation.IntentId,
			Status: metainstance.OperationStatus{
				State:   metainstance.OperationState(req.Operation.Status.State.String()),
				Message: req.Operation.Status.Message,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateMetaInstanceResponse{Record: s.ledgerRecordToProto(addOperationResponse.Record)}, nil
}

// UpdateOperationStatus updates the status of an operation on a MetaInstance
func (s *MetaInstanceService) UpdateOperationStatus(ctx context.Context, req *mrdspb.UpdateOperationStatusRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	updateOperationStatusResponse, err := s.ledger.UpdateOperationStatus(ctx, &metainstance.UpdateOperationStatusRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		OperationID: req.OperationId,
		Status: metainstance.OperationStatus{
			State:   metainstance.OperationState(req.Status.State.String()),
			Message: req.Status.Message,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateMetaInstanceResponse{Record: s.ledgerRecordToProto(updateOperationStatusResponse.Record)}, nil
}

// RemoveOperation removes an operation from a MetaInstance
func (s *MetaInstanceService) RemoveOperation(ctx context.Context, req *mrdspb.RemoveOperationRequest) (*mrdspb.UpdateMetaInstanceResponse, error) {
	removeOperationResponse, err := s.ledger.RemoveOperation(ctx, &metainstance.RemoveOperationRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		OperationID: req.OperationId,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateMetaInstanceResponse{Record: s.ledgerRecordToProto(removeOperationResponse.Record)}, nil
}
