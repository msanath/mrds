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
	ledgerRecordToProto func(record node.NodeRecord) *mrdspb.Node

	mrdspb.UnimplementedNodesServer
}

func nodeLedgerRecordToProto(record node.NodeRecord) *mrdspb.Node {
	node := &mrdspb.Node{
		Metadata: &mrdspb.Metadata{
			Id:      record.Metadata.ID,
			Version: record.Metadata.Version,
		},
		Name: record.Name,
		Status: &mrdspb.NodeStatus{
			State:   mrdspb.NodeState(mrdspb.NodeState_value[string(record.Status.State)]),
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
		CapabilityIds: record.CapabilityIDs,
	}

	for _, disruption := range record.Disruptions {
		node.Disruptions = append(node.Disruptions, &mrdspb.NodeDisruption{
			Id:          disruption.ID,
			ShouldEvict: disruption.ShouldEvict,
			Status: &mrdspb.DisruptionStatus{
				State:   mrdspb.DisruptionState(mrdspb.DisruptionState_value[string(disruption.Status.State)]),
				Message: disruption.Status.Message,
			},
		})
	}

	for _, localVolume := range record.LocalVolumes {
		node.LocalVolumes = append(node.LocalVolumes, &mrdspb.NodeLocalVolume{
			MountPath:       localVolume.MountPath,
			StorageClass:    localVolume.StorageClass,
			StorageCapacity: localVolume.StorageCapacity,
		})
	}

	return node
}

func NewNodeService(ledger node.Ledger) *NodeService {
	return &NodeService{
		ledger:              ledger,
		ledgerRecordToProto: nodeLedgerRecordToProto,
	}
}

// Create creates a new Node
func (s *NodeService) Create(ctx context.Context, req *mrdspb.CreateNodeRequest) (*mrdspb.CreateNodeResponse, error) {
	if req.TotalResources == nil || req.SystemReservedResources == nil {
		return nil, fmt.Errorf("TotalResources and SystemReservedResources are required")
	}
	cr := &node.CreateRequest{
		Name: req.Name,
		TotalResources: node.Resources{
			Cores:  req.TotalResources.Cores,
			Memory: req.TotalResources.Memory,
		},
		SystemReservedResources: node.Resources{
			Cores:  req.SystemReservedResources.Cores,
			Memory: req.SystemReservedResources.Memory,
		},
		UpdateDomain:  req.UpdateDomain,
		CapabilityIDs: req.CapabilityIds,
	}

	for _, localVolume := range req.LocalVolumes {
		cr.LocalVolumes = append(cr.LocalVolumes, node.LocalVolume{
			MountPath:       localVolume.MountPath,
			StorageClass:    localVolume.StorageClass,
			StorageCapacity: localVolume.StorageCapacity,
		})
	}

	createResponse, err := s.ledger.Create(ctx, cr)
	if err != nil {
		return nil, err
	}

	return &mrdspb.CreateNodeResponse{Record: s.ledgerRecordToProto(createResponse.Record)}, nil
}

// GetByID retrieves a Node by its ID
func (s *NodeService) GetByID(ctx context.Context, req *mrdspb.GetNodeByIDRequest) (*mrdspb.GetNodeResponse, error) {
	getResponse, err := s.ledger.GetByID(ctx, req.Id)
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
		stateIn[i] = node.NodeState(state.String())
	}

	stateNotIn := make([]node.NodeState, len(req.StateNotIn))
	for i, state := range req.StateNotIn {
		stateNotIn[i] = node.NodeState(state.String())
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

// AddDisruption adds a disruption to a Node
func (s *NodeService) AddDisruption(ctx context.Context, req *mrdspb.AddDisruptionRequest) (*mrdspb.UpdateNodeResponse, error) {
	addDisruptionResponse, err := s.ledger.AddDisruption(ctx, &node.AddDisruptionRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		Disruption: node.Disruption{
			ID:          req.Disruption.Id,
			ShouldEvict: req.Disruption.ShouldEvict,
			Status: node.DisruptionStatus{
				State:   node.DisruptionState(req.Disruption.Status.State.String()),
				Message: req.Disruption.Status.Message,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateNodeResponse{Record: s.ledgerRecordToProto(addDisruptionResponse.Record)}, nil
}

// UpdateDisruptionStatus updates the status of a disruption on a Node
func (s *NodeService) UpdateDisruptionStatus(ctx context.Context, req *mrdspb.UpdateDisruptionStatusRequest) (*mrdspb.UpdateNodeResponse, error) {
	updateDisruptionStatusResponse, err := s.ledger.UpdateDisruptionStatus(ctx, &node.UpdateDisruptionStatusRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		DisruptionID: req.DisruptionId,
		Status: node.DisruptionStatus{
			State:   node.DisruptionState(req.Status.State.String()),
			Message: req.Status.Message,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateNodeResponse{Record: s.ledgerRecordToProto(updateDisruptionStatusResponse.Record)}, nil
}

func (s *NodeService) RemoveDisruption(ctx context.Context, req *mrdspb.RemoveDisruptionRequest) (*mrdspb.UpdateNodeResponse, error) {
	removeDisruptionResponse, err := s.ledger.RemoveDisruption(ctx, &node.RemoveDisruptionRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		DisruptionID: req.DisruptionId,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateNodeResponse{Record: s.ledgerRecordToProto(removeDisruptionResponse.Record)}, nil
}

func (s *NodeService) AddCapability(ctx context.Context, req *mrdspb.AddCapabilityRequest) (*mrdspb.UpdateNodeResponse, error) {
	addCapabilityResponse, err := s.ledger.AddCapability(ctx, &node.AddCapabilityRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		CapabilityID: req.CapabilityId,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateNodeResponse{Record: s.ledgerRecordToProto(addCapabilityResponse.Record)}, nil
}

func (s *NodeService) RemoveCapability(ctx context.Context, req *mrdspb.RemoveCapabilityRequest) (*mrdspb.UpdateNodeResponse, error) {
	removeCapabilityResponse, err := s.ledger.RemoveCapability(ctx, &node.RemoveCapabilityRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		CapabilityID: req.CapabilityId,
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateNodeResponse{Record: s.ledgerRecordToProto(removeCapabilityResponse.Record)}, nil
}
