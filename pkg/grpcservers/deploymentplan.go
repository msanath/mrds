package grpcservers

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/internal/ledger/core"
	"github.com/msanath/mrds/internal/ledger/deploymentplan"
)

type DeploymentPlanService struct {
	ledger              deploymentplan.Ledger
	ledgerRecordToProto func(record deploymentplan.DeploymentPlanRecord) *mrdspb.DeploymentPlanRecord

	mrdspb.UnimplementedDeploymentPlansServer
}

func deploymentPlanLedgerRecordToProto(record deploymentplan.DeploymentPlanRecord) *mrdspb.DeploymentPlanRecord {
	return &mrdspb.DeploymentPlanRecord{
		Metadata: &mrdspb.Metadata{
			Id:      record.Metadata.ID,
			Version: record.Metadata.Version,
		},
		Name: record.Name,
		Status: &mrdspb.DeploymentPlanStatus{
			State:   mrdspb.DeploymentPlanState(mrdspb.DeploymentPlanState_value[string(record.Status.State)]),
			Message: record.Status.Message,
		},
		Namespace:    record.Namespace,
		ServiceName:  record.ServiceName,
		Deployments:  deploymentPlanDeploymentsToProto(record.Deployments),
		Applications: deploymentPlanApplicationsToProto(record.Applications),
	}
}

func deploymentPlanDeploymentsToProto(deployments []deploymentplan.Deployment) []*mrdspb.Deployment {
	var protoDeployments []*mrdspb.Deployment
	for _, d := range deployments {
		protoDeployments = append(protoDeployments, &mrdspb.Deployment{
			Id: d.ID,
			Status: &mrdspb.DeploymentStatus{
				State:   mrdspb.DeploymentState(mrdspb.DeploymentState_value[string(d.Status.State)]),
				Message: d.Status.Message,
			},
			PayloadCoordinates: deploymentPlanPayloadCoordinatesToProto(d.PayloadCoordinates),
			InstanceCount:      d.InstanceCount,
		})
	}
	return protoDeployments
}

func deploymentPlanApplicationsToProto(applications []deploymentplan.Application) []*mrdspb.Application {
	var protoApps []*mrdspb.Application
	for _, a := range applications {
		protoApps = append(protoApps, &mrdspb.Application{
			PayloadName: a.PayloadName,
			Resources: &mrdspb.ApplicationResources{
				Cores:  a.Resources.Cores,
				Memory: a.Resources.Memory,
			},
			Ports:             applicationPortsToProto(a.Ports),
			PersistentVolumes: applicationPersistentVolumesToProto(a.PersistentVolumes),
		})
	}
	return protoApps
}

func applicationPortsToProto(ports []deploymentplan.ApplicationPort) []*mrdspb.ApplicationPort {
	var protoPorts []*mrdspb.ApplicationPort
	for _, p := range ports {
		protoPorts = append(protoPorts, &mrdspb.ApplicationPort{
			Protocol: p.Protocol,
			Port:     p.Port,
		})
	}
	return protoPorts
}

func applicationPersistentVolumesToProto(persistentVolumes []deploymentplan.ApplicationPersistentVolume) []*mrdspb.ApplicationPersistentVolume {
	var protoPVs []*mrdspb.ApplicationPersistentVolume
	for _, pv := range persistentVolumes {
		protoPVs = append(protoPVs, &mrdspb.ApplicationPersistentVolume{
			StorageClass: pv.StorageClass,
			Capacity:     pv.Capacity,
			MountPath:    pv.MountPath,
		})
	}
	return protoPVs
}

func deploymentPlanPayloadCoordinatesToProto(coords []deploymentplan.PayloadCoordinates) []*mrdspb.PayloadCoordinates {
	var protoCoords []*mrdspb.PayloadCoordinates
	for _, coord := range coords {
		protoCoords = append(protoCoords, &mrdspb.PayloadCoordinates{
			PayloadName: coord.PayloadName,
			Coordinates: coord.Coordinates,
		})
	}
	return protoCoords
}

func NewDeploymentPlanService(ledger deploymentplan.Ledger) *DeploymentPlanService {
	return &DeploymentPlanService{
		ledger:              ledger,
		ledgerRecordToProto: deploymentPlanLedgerRecordToProto,
	}
}

// Create creates a new DeploymentPlan
func (s *DeploymentPlanService) Create(ctx context.Context, req *mrdspb.CreateDeploymentPlanRequest) (*mrdspb.CreateDeploymentPlanResponse, error) {
	// Convert MatchingComputeCapabilities from protobuf type to interface type
	var matchingComputeCapabilities []deploymentplan.MatchingComputeCapability
	for _, cap := range req.MatchingComputeCapabilities {
		matchingComputeCapabilities = append(matchingComputeCapabilities, deploymentplan.MatchingComputeCapability{
			CapabilityType:  cap.CapabilityType,
			Comparator:      deploymentplan.ComparatorType(cap.Comparator.String()),
			CapabilityNames: cap.CapabilityNames,
		})
	}

	// Convert Applications from protobuf type to interface type
	var applications []deploymentplan.Application
	for _, app := range req.Applications {
		// Convert Application resources
		resources := deploymentplan.ApplicationResources{
			Cores:  app.Resources.Cores,
			Memory: app.Resources.Memory,
		}

		// Convert Application ports
		var ports []deploymentplan.ApplicationPort
		for _, port := range app.Ports {
			ports = append(ports, deploymentplan.ApplicationPort{
				Protocol: port.Protocol,
				Port:     port.Port,
			})
		}

		// Convert Application persistent volumes
		var persistentVolumes []deploymentplan.ApplicationPersistentVolume
		for _, pv := range app.PersistentVolumes {
			persistentVolumes = append(persistentVolumes, deploymentplan.ApplicationPersistentVolume{
				StorageClass: pv.StorageClass,
				Capacity:     pv.Capacity,
				MountPath:    pv.MountPath,
			})
		}

		// Construct the Application
		applications = append(applications, deploymentplan.Application{
			PayloadName:       app.PayloadName,
			Resources:         resources,
			Ports:             ports,
			PersistentVolumes: persistentVolumes,
		})
	}

	// Construct the CreateRequest with converted fields
	cr := &deploymentplan.CreateRequest{
		Name:                        req.Name,
		Namespace:                   req.Namespace,
		ServiceName:                 req.ServiceName,
		MatchingComputeCapabilities: matchingComputeCapabilities,
		Applications:                applications,
	}

	// Call the ledger's Create function
	createResponse, err := s.ledger.Create(ctx, cr)
	if err != nil {
		return nil, err
	}

	// Return the response
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

// UpdateStatus updates the state of an existing DeploymentPlan
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
	// Convert StateIn from protobuf enum to interface enum
	var stateIn []deploymentplan.DeploymentPlanState
	for _, state := range req.Filters.StateIn {
		stateIn = append(stateIn, deploymentplan.DeploymentPlanState(state.String()))
	}

	// Convert StateNotIn from protobuf enum to interface enum
	var stateNotIn []deploymentplan.DeploymentPlanState
	for _, state := range req.Filters.StateNotIn {
		stateNotIn = append(stateNotIn, deploymentplan.DeploymentPlanState(state.String()))
	}

	// Build the ListRequest with converted fields
	listResponse, err := s.ledger.List(ctx, &deploymentplan.ListRequest{
		Filters: deploymentplan.DeploymentPlanListFilters{
			IDIn:               req.Filters.IdIn,
			NameIn:             req.Filters.NameIn,
			VersionGte:         req.Filters.VersionGte,
			VersionLte:         req.Filters.VersionLte,
			VersionEq:          req.Filters.VersionEq,
			StateIn:            stateIn,
			StateNotIn:         stateNotIn,
			IncludeDeleted:     req.Filters.IncludeDeleted,
			Limit:              req.Filters.Limit,
			ServiceNameIn:      req.Filters.ServiceNameIn,
			PayloadNameIn:      req.Filters.PayloadNameIn,
			DeploymentPlanIDIn: req.Filters.DeploymentPlanIdIn,
		},
	})
	if err != nil {
		return nil, err
	}

	// Convert list of DeploymentPlanRecords to protobuf
	records := make([]*mrdspb.DeploymentPlanRecord, len(listResponse.Records))
	for i, record := range listResponse.Records {
		records[i] = s.ledgerRecordToProto(record)
	}

	// Return the list response
	return &mrdspb.ListDeploymentPlanResponse{Records: records}, nil
}

// Delete deletes a DeploymentPlan by metadata
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
	return &mrdspb.DeleteDeploymentPlanResponse{Success: true}, nil
}

// AddDeployment adds a Deployment to an existing DeploymentPlan
func (s *DeploymentPlanService) AddDeployment(ctx context.Context, req *mrdspb.AddDeploymentRequest) (*mrdspb.UpdateDeploymentPlanResponse, error) {
	var payloadCoordinates []deploymentplan.PayloadCoordinates
	for _, coord := range req.PayloadCoordinates {
		// Convert the map[string]string Coordinates
		coordinates := make(map[string]string)
		for k, v := range coord.Coordinates {
			coordinates[k] = v
		}

		// Construct the PayloadCoordinates object
		payloadCoordinates = append(payloadCoordinates, deploymentplan.PayloadCoordinates{
			PayloadName: coord.PayloadName,
			Coordinates: coordinates,
		})
	}

	// Build the AddDeploymentRequest with converted fields
	addResponse, err := s.ledger.AddDeployment(ctx, &deploymentplan.AddDeploymentRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		DeploymentID:       req.DeploymentId,
		PayloadCoordinates: payloadCoordinates,
		InstanceCount:      req.InstanceCount,
	})
	if err != nil {
		return nil, err
	}

	// Return the response
	return &mrdspb.UpdateDeploymentPlanResponse{Record: s.ledgerRecordToProto(addResponse.Record)}, nil
}

// UpdateDeploymentStatus updates the status of an existing Deployment
func (s *DeploymentPlanService) UpdateDeploymentStatus(ctx context.Context, req *mrdspb.UpdateDeploymentStatusRequest) (*mrdspb.UpdateDeploymentPlanResponse, error) {
	updateResponse, err := s.ledger.UpdateDeploymentStatus(ctx, &deploymentplan.UpdateDeploymentStatusRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		DeploymentID: req.DeploymentId,
		Status: deploymentplan.DeploymentStatus{
			State:   deploymentplan.DeploymentState(req.Status.State.String()),
			Message: req.Status.Message,
		},
	})
	if err != nil {
		return nil, err
	}
	return &mrdspb.UpdateDeploymentPlanResponse{Record: s.ledgerRecordToProto(updateResponse.Record)}, nil
}
