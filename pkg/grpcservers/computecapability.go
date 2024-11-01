package grpcservers

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/internal/ledger/computecapability"
	"github.com/msanath/mrds/internal/ledger/core"
)

type ComputeCapabilityService struct {
	ledger              computecapability.Ledger
	protoToLedgerRecord func(proto *mrdspb.ComputeCapability) computecapability.ComputeCapabilityRecord
	ledgerRecordToProto func(record computecapability.ComputeCapabilityRecord) *mrdspb.ComputeCapability

	mrdspb.UnimplementedComputeCapabilitiesServer
}

func computeCapabilityProtoToLedgerRecord(proto *mrdspb.ComputeCapability) computecapability.ComputeCapabilityRecord {
	return computecapability.ComputeCapabilityRecord{
		Metadata: core.Metadata{
			ID:      proto.Metadata.Id,
			Version: proto.Metadata.Version,
		},
		Name: proto.Name,
		Status: computecapability.ComputeCapabilityStatus{
			State:   computecapability.ComputeCapabilityState(proto.Status.State.String()),
			Message: proto.Status.Message,
		},
		Type:  proto.Type,
		Score: proto.Score,
	}
}

func computeCapabilityLedgerRecordToProto(record computecapability.ComputeCapabilityRecord) *mrdspb.ComputeCapability {
	return &mrdspb.ComputeCapability{
		Metadata: &mrdspb.Metadata{
			Id:      record.Metadata.ID,
			Version: record.Metadata.Version,
		},
		Name: record.Name,
		Status: &mrdspb.ComputeCapabilityStatus{
			State:   mrdspb.ComputeCapabilityState(mrdspb.ComputeCapabilityState_value[record.Status.State.ToString()]),
			Message: record.Status.Message,
		},
		Type:  record.Type,
		Score: record.Score,
	}
}

func NewComputeCapabilityService(ledger computecapability.Ledger) *ComputeCapabilityService {
	return &ComputeCapabilityService{
		ledger:              ledger,
		protoToLedgerRecord: computeCapabilityProtoToLedgerRecord,
		ledgerRecordToProto: computeCapabilityLedgerRecordToProto,
	}
}

// Create creates a new ComputeCapability
func (s *ComputeCapabilityService) Create(ctx context.Context, req *mrdspb.CreateComputeCapabilityRequest) (*mrdspb.CreateComputeCapabilityResponse, error) {
	createResponse, err := s.ledger.Create(ctx, &computecapability.CreateRequest{
		Name:  req.Name,
		Type:  req.Type,
		Score: req.Score,
	})
	if err != nil {
		return nil, err
	}

	return &mrdspb.CreateComputeCapabilityResponse{Record: s.ledgerRecordToProto(createResponse.Record)}, nil
}

// GetByName retrieves a ComputeCapability by its name
func (s *ComputeCapabilityService) GetByID(ctx context.Context, req *mrdspb.GetComputeCapabilityByIDRequest) (*mrdspb.GetComputeCapabilityResponse, error) {
	getResponse, err := s.ledger.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetComputeCapabilityResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// GetByName retrieves a ComputeCapability by its name
func (s *ComputeCapabilityService) GetByName(ctx context.Context, req *mrdspb.GetComputeCapabilityByNameRequest) (*mrdspb.GetComputeCapabilityResponse, error) {
	getResponse, err := s.ledger.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetComputeCapabilityResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// UpdateState updates the state and message of an existing ComputeCapability
func (s *ComputeCapabilityService) UpdateStatus(ctx context.Context, req *mrdspb.UpdateComputeCapabilityStatusRequest) (*mrdspb.UpdateComputeCapabilityResponse, error) {
	updateResponse, err := s.ledger.UpdateStatus(ctx, &computecapability.UpdateStateRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		Status: computecapability.ComputeCapabilityStatus{
			State:   computecapability.ComputeCapabilityState(req.Status.State.String()),
			Message: req.Status.Message,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateComputeCapabilityResponse{Record: s.ledgerRecordToProto(updateResponse.Record)}, nil
}

// List returns a list of ComputeCapabilitys that match the provided filters
func (s *ComputeCapabilityService) List(ctx context.Context, req *mrdspb.ListComputeCapabilityRequest) (*mrdspb.ListComputeCapabilityResponse, error) {
	if req == nil {
		req = &mrdspb.ListComputeCapabilityRequest{}
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

	stateIn := make([]computecapability.ComputeCapabilityState, len(req.StateIn))
	for i, state := range req.StateIn {
		stateIn[i] = computecapability.ComputeCapabilityStateFromString(state.String())
	}

	stateNotIn := make([]computecapability.ComputeCapabilityState, len(req.StateNotIn))
	for i, state := range req.StateNotIn {
		stateNotIn[i] = computecapability.ComputeCapabilityStateFromString(state.String())
	}

	listResponse, err := s.ledger.List(ctx, &computecapability.ListRequest{
		Filters: computecapability.ComputeCapabilityListFilters{
			IDIn:       req.IdIn,
			NameIn:     req.NameIn,
			VersionGte: gte,
			VersionLte: lte,
			VersionEq:  eq,
			StateIn:    stateIn,
			StateNotIn: stateNotIn,
			TypeIn:     req.TypeIn,
		},
	})
	if err != nil {
		return nil, err
	}

	records := make([]*mrdspb.ComputeCapability, len(listResponse.Records))
	for i, record := range listResponse.Records {
		records[i] = s.ledgerRecordToProto(record)
	}

	return &mrdspb.ListComputeCapabilityResponse{Records: records}, nil
}

// Delete deletes a ComputeCapability
func (s *ComputeCapabilityService) Delete(ctx context.Context, req *mrdspb.DeleteComputeCapabilityRequest) (*mrdspb.DeleteComputeCapabilityResponse, error) {
	err := s.ledger.Delete(ctx, &computecapability.DeleteRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.DeleteComputeCapabilityResponse{}, nil
}
