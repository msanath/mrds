package grpcservers

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/metainstance"
)

type MetaInstanceService struct {
	ledger              metainstance.Ledger
	protoToLedgerRecord func(proto *mrdspb.MetaInstance) metainstance.MetaInstanceRecord
	ledgerRecordToProto func(record metainstance.MetaInstanceRecord) *mrdspb.MetaInstance

	mrdspb.UnimplementedMetaInstancesServer
}

func metaInstanceProtoToLedgerRecord(proto *mrdspb.MetaInstance) metainstance.MetaInstanceRecord {
	return metainstance.MetaInstanceRecord{
		Metadata: core.Metadata{
			ID:      proto.Metadata.Id,
			Version: proto.Metadata.Version,
		},
		Name: proto.Name,
		Status: metainstance.MetaInstanceStatus{
			State:   metainstance.MetaInstanceState(proto.Status.State.String()),
			Message: proto.Status.Message,
		},
	}
}

func metaInstanceLedgerRecordToProto(record metainstance.MetaInstanceRecord) *mrdspb.MetaInstance {
	return &mrdspb.MetaInstance{
		Metadata: &mrdspb.Metadata{
			Id:      record.Metadata.ID,
			Version: record.Metadata.Version,
		},
		Name: record.Name,
		Status: &mrdspb.MetaInstanceStatus{
			State:   mrdspb.MetaInstanceState(mrdspb.MetaInstanceState_value[record.Status.State.ToString()]),
			Message: record.Status.Message,
		},
	}
}

func NewMetaInstanceService(ledger metainstance.Ledger) *MetaInstanceService {
	return &MetaInstanceService{
		ledger:              ledger,
		protoToLedgerRecord: metaInstanceProtoToLedgerRecord,
		ledgerRecordToProto: metaInstanceLedgerRecordToProto,
	}
}

// Create creates a new MetaInstance
func (s *MetaInstanceService) Create(ctx context.Context, req *mrdspb.CreateMetaInstanceRequest) (*mrdspb.CreateMetaInstanceResponse, error) {
	createResponse, err := s.ledger.Create(ctx, &metainstance.CreateRequest{
		Name: req.Name,
	})
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
		stateIn[i] = metainstance.MetaInstanceStateFromString(state.String())
	}

	stateNotIn := make([]metainstance.MetaInstanceState, len(req.StateNotIn))
	for i, state := range req.StateNotIn {
		stateNotIn[i] = metainstance.MetaInstanceStateFromString(state.String())
	}

	listResponse, err := s.ledger.List(ctx, &metainstance.ListRequest{
		Filters: metainstance.MetaInstanceListFilters{
			IDIn:       req.IdIn,
			NameIn:     req.NameIn,
			VersionGte: gte,
			VersionLte: lte,
			VersionEq:  eq,
			StateIn:    stateIn,
			StateNotIn: stateNotIn,
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
