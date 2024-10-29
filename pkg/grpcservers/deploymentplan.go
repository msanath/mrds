
package grpcservers

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/internal/ledger/deploymentplan"
	"github.com/msanath/mrds/internal/ledger/core"
)

type DeploymentPlanService struct {
	ledger              deploymentplan.Ledger
	protoToLedgerRecord func(proto *mrdspb.DeploymentPlan) deploymentplan.DeploymentPlanRecord
	ledgerRecordToProto func(record deploymentplan.DeploymentPlanRecord) *mrdspb.DeploymentPlan

	mrdspb.UnimplementedDeploymentPlansServer
}

func deploymentPlanProtoToLedgerRecord(proto *mrdspb.DeploymentPlan) deploymentplan.DeploymentPlanRecord {
	return deploymentplan.DeploymentPlanRecord{
		Metadata: core.Metadata{
			ID:      proto.Metadata.Id,
			Version: proto.Metadata.Version,
		},
		Name: proto.Name,
		Status: deploymentplan.DeploymentPlanStatus{
			State:   deploymentplan.DeploymentPlanState(proto.Status.State.String()),
			Message: proto.Status.Message,
		},
	}
}

func deploymentPlanLedgerRecordToProto(record deploymentplan.DeploymentPlanRecord) *mrdspb.DeploymentPlan {
	return &mrdspb.DeploymentPlan{
		Metadata: &mrdspb.Metadata{
			Id:      record.Metadata.ID,
			Version: record.Metadata.Version,
		},
		Name: record.Name,
		Status: &mrdspb.DeploymentPlanStatus{
			State:   mrdspb.DeploymentPlanState(mrdspb.DeploymentPlanState_value[record.Status.State.ToString()]),
			Message: record.Status.Message,
		},
	}
}

func NewDeploymentPlanService(ledger deploymentplan.Ledger) *DeploymentPlanService {
	return &DeploymentPlanService{
		ledger:              ledger,
		protoToLedgerRecord: deploymentPlanProtoToLedgerRecord,
		ledgerRecordToProto: deploymentPlanLedgerRecordToProto,
	}
}

// Create creates a new DeploymentPlan
func (s *DeploymentPlanService) Create(ctx context.Context, req *mrdspb.CreateDeploymentPlanRequest) (*mrdspb.CreateDeploymentPlanResponse, error) {
	createResponse, err := s.ledger.Create(ctx, &deploymentplan.CreateRequest{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &mrdspb.CreateDeploymentPlanResponse{Record: s.ledgerRecordToProto(createResponse.Record)}, nil
}

// GetByMetadata retrieves a DeploymentPlan by its metadata
func (s *DeploymentPlanService) GetByMetadata(ctx context.Context, req *mrdspb.GetDeploymentPlanByMetadataRequest) (*mrdspb.GetDeploymentPlanResponse, error) {
	getResponse, err := s.ledger.GetByMetadata(ctx, &core.Metadata{
		ID:      req.Metadata.Id,
		Version: req.Metadata.Version,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetDeploymentPlanResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// GetByName retrieves a DeploymentPlan by its name
func (s *DeploymentPlanService) GetByName(ctx context.Context, req *mrdspb.GetDeploymentPlanByNameRequest) (*mrdspb.GetDeploymentPlanResponse, error) {
	getResponse, err := s.ledger.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetDeploymentPlanResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// UpdateStatus updates the state and message of an existing DeploymentPlan
func (s *DeploymentPlanService) UpdateStatus(ctx context.Context, req *mrdspb.UpdateDeploymentPlanStatusRequest) (*mrdspb.UpdateDeploymentPlanResponse, error) {
	updateResponse, err := s.ledger.UpdateStatus(ctx, &deploymentplan.UpdateStatusRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		Status: deploymentplan.DeploymentPlanStatus{
			State:   deploymentplan.DeploymentPlanState(req.Status.State.String()),
			Message: req.Status.Message,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateDeploymentPlanResponse{Record: s.ledgerRecordToProto(updateResponse.Record)}, nil
}

// List returns a list of DeploymentPlans that match the provided filters
func (s *DeploymentPlanService) List(ctx context.Context, req *mrdspb.ListDeploymentPlanRequest) (*mrdspb.ListDeploymentPlanResponse, error) {
	if req == nil {
		req = &mrdspb.ListDeploymentPlanRequest{}
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

	stateIn := make([]deploymentplan.DeploymentPlanState, len(req.StateIn))
	for i, state := range req.StateIn {
		stateIn[i] = deploymentplan.DeploymentPlanStateFromString(state.String())
	}

	stateNotIn := make([]deploymentplan.DeploymentPlanState, len(req.StateNotIn))
	for i, state := range req.StateNotIn {
		stateNotIn[i] = deploymentplan.DeploymentPlanStateFromString(state.String())
	}

	listResponse, err := s.ledger.List(ctx, &deploymentplan.ListRequest{
		Filters: deploymentplan.DeploymentPlanListFilters{
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

	records := make([]*mrdspb.DeploymentPlan, len(listResponse.Records))
	for i, record := range listResponse.Records {
		records[i] = s.ledgerRecordToProto(record)
	}

	return &mrdspb.ListDeploymentPlanResponse{Records: records}, nil
}

// Delete deletes a DeploymentPlan
func (s *DeploymentPlanService) Delete(ctx context.Context, req *mrdspb.DeleteDeploymentPlanRequest) (*mrdspb.DeleteDeploymentPlanResponse, error) {
	err := s.ledger.Delete(ctx, &deploymentplan.DeleteRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.DeleteDeploymentPlanResponse{}, nil
}
