package grpcservers

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/node"
)

type NodeService struct {
	ledger              node.Ledger
	protoToLedgerRecord func(proto *mrdspb.Node) node.NodeRecord
	ledgerRecordToProto func(record node.NodeRecord) *mrdspb.Node

	mrdspb.UnimplementedNodesServer
}

func nodeProtoToLedgerRecord(proto *mrdspb.Node) node.NodeRecord {
	return node.NodeRecord{
		Metadata: core.Metadata{
			ID:      proto.Metadata.Id,
			Version: proto.Metadata.Version,
		},
		Name: proto.Name,
		Status: node.NodeStatus{
			State:   node.NodeState(proto.Status.State.String()),
			Message: proto.Status.Message,
		},
		ClusterID:    proto.ClusterId,
		UpdateDomain: proto.UpdateDomain,
		TotalResources: node.Resources{
			Cores:  proto.TotalResources.Cores,
			Memory: proto.TotalResources.Memory,
		},
		SystemReservedResources: node.Resources{
			Cores:  proto.SystemReservedResources.Cores,
			Memory: proto.SystemReservedResources.Memory,
		},
		RemainingResources: node.Resources{
			Cores:  proto.RemainingResources.Cores,
			Memory: proto.RemainingResources.Memory,
		},
	}
}

func nodeLedgerRecordToProto(record node.NodeRecord) *mrdspb.Node {
	return &mrdspb.Node{
		Metadata: &mrdspb.Metadata{
			Id:      record.Metadata.ID,
			Version: record.Metadata.Version,
		},
		Name: record.Name,
		Status: &mrdspb.NodeStatus{
			State:   mrdspb.NodeState(mrdspb.NodeState_value[record.Status.State.ToString()]),
			Message: record.Status.Message,
		},
		ClusterId:    record.ClusterID,
		UpdateDomain: record.UpdateDomain,
		TotalResources: &mrdspb.Resources{
			Cores:  record.TotalResources.Cores,
			Memory: record.TotalResources.Memory,
		},
		SystemReservedResources: &mrdspb.Resources{
			Cores:  record.SystemReservedResources.Cores,
			Memory: record.SystemReservedResources.Memory,
		},
		RemainingResources: &mrdspb.Resources{
			Cores:  record.RemainingResources.Cores,
			Memory: record.RemainingResources.Memory,
		},
	}
}

func NewNodeService(ledger node.Ledger) *NodeService {
	return &NodeService{
		ledger:              ledger,
		protoToLedgerRecord: nodeProtoToLedgerRecord,
		ledgerRecordToProto: nodeLedgerRecordToProto,
	}
}

// Create creates a new Node
func (s *NodeService) Create(ctx context.Context, req *mrdspb.CreateNodeRequest) (*mrdspb.CreateNodeResponse, error) {
	if req.TotalResources == nil || req.SystemReservedResources == nil {
		return nil, fmt.Errorf("TotalResources and SystemReservedResources are required")
	}
	createResponse, err := s.ledger.Create(ctx, &node.CreateRequest{
		Name: req.Name,
		TotalResources: node.Resources{
			Cores:  req.TotalResources.Cores,
			Memory: req.TotalResources.Memory,
		},
		SystemReservedResources: node.Resources{
			Cores:  req.SystemReservedResources.Cores,
			Memory: req.SystemReservedResources.Memory,
		},
		UpdateDomain: req.UpdateDomain,
	})
	if err != nil {
		return nil, err
	}

	return &mrdspb.CreateNodeResponse{Record: s.ledgerRecordToProto(createResponse.Record)}, nil
}

// GetByMetadata retrieves a Node by its metadata
func (s *NodeService) GetByMetadata(ctx context.Context, req *mrdspb.GetNodeByMetadataRequest) (*mrdspb.GetNodeResponse, error) {
	getResponse, err := s.ledger.GetByMetadata(ctx, &core.Metadata{
		ID:      req.Metadata.Id,
		Version: req.Metadata.Version,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetNodeResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// GetByName retrieves a Node by its name
func (s *NodeService) GetByName(ctx context.Context, req *mrdspb.GetNodeByNameRequest) (*mrdspb.GetNodeResponse, error) {
	getResponse, err := s.ledger.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &mrdspb.GetNodeResponse{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// UpdateState updates the state and message of an existing Node
func (s *NodeService) UpdateStatus(ctx context.Context, req *mrdspb.UpdateNodeStatusRequest) (*mrdspb.UpdateNodeResponse, error) {
	updateResponse, err := s.ledger.UpdateStatus(ctx, &node.UpdateStatusRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		Status: node.NodeStatus{
			State:   node.NodeState(req.Status.State.String()),
			Message: req.Status.Message,
		},
		ClusterID: req.ClusterId,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateNodeResponse{Record: s.ledgerRecordToProto(updateResponse.Record)}, nil
}

// List returns a list of Nodes that match the provided filters
func (s *NodeService) List(ctx context.Context, req *mrdspb.ListNodeRequest) (*mrdspb.ListNodeResponse, error) {
	if req == nil {
		req = &mrdspb.ListNodeRequest{}
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

	var remainingCoresGte, remainingCoresLte *uint32
	if req.RemainingCoresGte != 0 {
		remainingCoresGte = &req.RemainingCoresGte
	}
	if req.RemainingCoresLte != 0 {
		remainingCoresLte = &req.RemainingCoresLte
	}

	var remainingMemoryGte, remainingMemoryLte *uint32
	if req.RemainingMemoryGte != 0 {
		remainingMemoryGte = &req.RemainingMemoryGte
	}
	if req.RemainingMemoryLte != 0 {
		remainingMemoryLte = &req.RemainingMemoryLte
	}

	stateIn := make([]node.NodeState, len(req.StateIn))
	for i, state := range req.StateIn {
		stateIn[i] = node.NodeStateFromString(state.String())
	}

	stateNotIn := make([]node.NodeState, len(req.StateNotIn))
	for i, state := range req.StateNotIn {
		stateNotIn[i] = node.NodeStateFromString(state.String())
	}

	listResponse, err := s.ledger.List(ctx, &node.ListRequest{
		Filters: node.NodeListFilters{
			IDIn:               req.IdIn,
			NameIn:             req.NameIn,
			VersionGte:         gte,
			VersionLte:         lte,
			VersionEq:          eq,
			StateIn:            stateIn,
			StateNotIn:         stateNotIn,
			RemainingCoresGte:  remainingCoresGte,
			RemainingCoresLte:  remainingCoresLte,
			RemainingMemoryGte: remainingMemoryGte,
			RemainingMemoryLte: remainingMemoryLte,
			ClusterIDIn:        req.ClusterIdIn,
			UpdateDomainIn:     req.UpdateDomainIn,
		},
	})
	if err != nil {
		return nil, err
	}

	records := make([]*mrdspb.Node, len(listResponse.Records))
	for i, record := range listResponse.Records {
		records[i] = s.ledgerRecordToProto(record)
	}

	return &mrdspb.ListNodeResponse{Records: records}, nil
}

// Delete deletes a Node
func (s *NodeService) Delete(ctx context.Context, req *mrdspb.DeleteNodeRequest) (*mrdspb.DeleteNodeResponse, error) {
	err := s.ledger.Delete(ctx, &node.DeleteRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.DeleteNodeResponse{}, nil
}
