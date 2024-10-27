package grpcservers

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/deployment"
)

type DeploymentService struct {
	ledger              deployment.Ledger
	protoToLedgerRecord func(proto *mrdspb.Deployment) deployment.DeploymentRecord
	ledgerRecordToProto func(record deployment.DeploymentRecord) *mrdspb.Deployment

	mrdspb.UnimplementedDeploymentsServer
}

func deploymentProtoToLedgerRecord(proto *mrdspb.Deployment) deployment.DeploymentRecord {
	return deployment.DeploymentRecord{
		Metadata: core.Metadata{
			ID:      proto.Metadata.Id,
			Version: proto.Metadata.Version,
		},
		Name: proto.Name,
		Status: deployment.DeploymentStatus{
			State:   deployment.DeploymentState(proto.Status.State.String()),
			Message: proto.Status.Message,
		},
	}
}

func deploymentLedgerRecordToProto(record deployment.DeploymentRecord) *mrdspb.Deployment {
	return &mrdspb.Deployment{
		Metadata: &mrdspb.Metadata{
			Id:      record.Metadata.ID,
			Version: record.Metadata.Version,
		},
		Name: record.Name,
		Status: &mrdspb.DeploymentStatus{
			State:   mrdspb.DeploymentState(mrdspb.DeploymentState_value[record.Status.State.ToString()]),
			Message: record.Status.Message,
		},
	}
}

func NewDeploymentService(ledger deployment.Ledger) *DeploymentService {
	return &DeploymentService{
		ledger:              ledger,
		protoToLedgerRecord: deploymentProtoToLedgerRecord,
		ledgerRecordToProto: deploymentLedgerRecordToProto,
	}
}

// Create creates a new Deployment
func (s *DeploymentService) Create(ctx context.Context, req *mrdspb.CreateDeploymentRequest) (*mrdspb.CreateDeploymentResponse, error) {
	createResponse, err := s.ledger.Create(ctx, &deployment.CreateRequest{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &mrdspb.CreateDeploymentResponse{Record: s.ledgerRecordToProto(createResponse.Record)}, nil
}

// GetByMetadata retrieves a Deployment by its metadata
func (s *DeploymentService) GetByMetadata(ctx context.Context, req *mrdspb.GetDeploymentByMetadataRequest) (*mrdspb.GetDeploymentResponse, error) {
	getResponse, err := s.ledger.GetByMetadata(ctx, &core.Metadata{
		ID:      req.Metadata.Id,
		Version: req.Metadata.Version,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetDeploymentResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// GetByName retrieves a Deployment by its name
func (s *DeploymentService) GetByName(ctx context.Context, req *mrdspb.GetDeploymentByNameRequest) (*mrdspb.GetDeploymentResponse, error) {
	getResponse, err := s.ledger.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetDeploymentResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// UpdateStatus updates the state and message of an existing Deployment
func (s *DeploymentService) UpdateStatus(ctx context.Context, req *mrdspb.UpdateDeploymentStatusRequest) (*mrdspb.UpdateDeploymentResponse, error) {
	updateResponse, err := s.ledger.UpdateStatus(ctx, &deployment.UpdateStatusRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		Status: deployment.DeploymentStatus{
			State:   deployment.DeploymentState(req.Status.State.String()),
			Message: req.Status.Message,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateDeploymentResponse{Record: s.ledgerRecordToProto(updateResponse.Record)}, nil
}

// List returns a list of Deployments that match the provided filters
func (s *DeploymentService) List(ctx context.Context, req *mrdspb.ListDeploymentRequest) (*mrdspb.ListDeploymentResponse, error) {
	if req == nil {
		req = &mrdspb.ListDeploymentRequest{}
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

	stateIn := make([]deployment.DeploymentState, len(req.StateIn))
	for i, state := range req.StateIn {
		stateIn[i] = deployment.DeploymentStateFromString(state.String())
	}

	stateNotIn := make([]deployment.DeploymentState, len(req.StateNotIn))
	for i, state := range req.StateNotIn {
		stateNotIn[i] = deployment.DeploymentStateFromString(state.String())
	}

	listResponse, err := s.ledger.List(ctx, &deployment.ListRequest{
		Filters: deployment.DeploymentListFilters{
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

	records := make([]*mrdspb.Deployment, len(listResponse.Records))
	for i, record := range listResponse.Records {
		records[i] = s.ledgerRecordToProto(record)
	}

	return &mrdspb.ListDeploymentResponse{Records: records}, nil
}

// Delete deletes a Deployment
func (s *DeploymentService) Delete(ctx context.Context, req *mrdspb.DeleteDeploymentRequest) (*mrdspb.DeleteDeploymentResponse, error) {
	err := s.ledger.Delete(ctx, &deployment.DeleteRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.DeleteDeploymentResponse{}, nil
}
