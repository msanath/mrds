package grpcservers

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/ledger/cluster"
	"github.com/msanath/mrds/ledger/core"
)

type ClusterService struct {
	ledger              cluster.Ledger
	protoToLedgerRecord func(proto *mrdspb.Cluster) cluster.ClusterRecord
	ledgerRecordToProto func(record cluster.ClusterRecord) *mrdspb.Cluster

	mrdspb.UnimplementedClustersServer
}

func clusterProtoToLedgerRecord(proto *mrdspb.Cluster) cluster.ClusterRecord {
	return cluster.ClusterRecord{
		Metadata: core.Metadata{
			ID:      proto.Metadata.Id,
			Version: proto.Metadata.Version,
		},
		Name: proto.Name,
		Status: cluster.ClusterStatus{
			State:   cluster.ClusterState(proto.Status.State.String()),
			Message: proto.Status.Message,
		},
	}
}

func clusterLedgerRecordToProto(record cluster.ClusterRecord) *mrdspb.Cluster {
	return &mrdspb.Cluster{
		Metadata: &mrdspb.Metadata{
			Id:      record.Metadata.ID,
			Version: record.Metadata.Version,
		},
		Name: record.Name,
		Status: &mrdspb.ClusterStatus{
			State:   mrdspb.ClusterState(mrdspb.ClusterState_value[record.Status.State.ToString()]),
			Message: record.Status.Message,
		},
	}
}

func NewClusterService(ledger cluster.Ledger) *ClusterService {
	return &ClusterService{
		ledger:              ledger,
		protoToLedgerRecord: clusterProtoToLedgerRecord,
		ledgerRecordToProto: clusterLedgerRecordToProto,
	}
}

// Create creates a new Cluster
func (s *ClusterService) Create(ctx context.Context, req *mrdspb.CreateClusterRequest) (*mrdspb.CreateClusterResponse, error) {
	createResponse, err := s.ledger.Create(ctx, &cluster.CreateRequest{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &mrdspb.CreateClusterResponse{Record: s.ledgerRecordToProto(createResponse.Record)}, nil
}

// GetByMetadata retrieves a Cluster by its metadata
func (s *ClusterService) GetByID(ctx context.Context, req *mrdspb.GetClusterByIDRequest) (*mrdspb.GetClusterResponse, error) {
	getResponse, err := s.ledger.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetClusterResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// GetByName retrieves a Cluster by its name
func (s *ClusterService) GetByName(ctx context.Context, req *mrdspb.GetClusterByNameRequest) (*mrdspb.GetClusterResponse, error) {
	getResponse, err := s.ledger.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetClusterResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// UpdateState updates the state and message of an existing Cluster
func (s *ClusterService) UpdateStatus(ctx context.Context, req *mrdspb.UpdateClusterStatusRequest) (*mrdspb.UpdateClusterResponse, error) {
	updateResponse, err := s.ledger.UpdateStatus(ctx, &cluster.UpdateStateRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		Status: cluster.ClusterStatus{
			State:   cluster.ClusterState(req.Status.State.String()),
			Message: req.Status.Message,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateClusterResponse{Record: s.ledgerRecordToProto(updateResponse.Record)}, nil
}

// List returns a list of Clusters that match the provided filters
func (s *ClusterService) List(ctx context.Context, req *mrdspb.ListClusterRequest) (*mrdspb.ListClusterResponse, error) {
	if req == nil {
		req = &mrdspb.ListClusterRequest{}
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

	stateIn := make([]cluster.ClusterState, len(req.StateIn))
	for i, state := range req.StateIn {
		stateIn[i] = cluster.ClusterStateFromString(state.String())
	}

	stateNotIn := make([]cluster.ClusterState, len(req.StateNotIn))
	for i, state := range req.StateNotIn {
		stateNotIn[i] = cluster.ClusterStateFromString(state.String())
	}

	listResponse, err := s.ledger.List(ctx, &cluster.ListRequest{
		Filters: cluster.ClusterListFilters{
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

	records := make([]*mrdspb.Cluster, len(listResponse.Records))
	for i, record := range listResponse.Records {
		records[i] = s.ledgerRecordToProto(record)
	}

	return &mrdspb.ListClusterResponse{Records: records}, nil
}

// Delete deletes a Cluster
func (s *ClusterService) Delete(ctx context.Context, req *mrdspb.DeleteClusterRequest) (*mrdspb.DeleteClusterResponse, error) {
	err := s.ledger.Delete(ctx, &cluster.DeleteRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.DeleteClusterResponse{}, nil
}
