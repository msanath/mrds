package deploymentplan

import (
	"context"
	"os"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/deploymentplan/printer"
	"github.com/msanath/mrds/pkg/ctl/deploymentplan/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

type createDeploymentPlanOptions struct {
	manifestFilePath string

	deploymentPlansClient mrdspb.DeploymentPlansClient
	printer               *printer.Printer
}

type planCreateRequest struct {
	Name                        string                     `yaml:"name"`
	Namespace                   string                     `yaml:"namespace"`
	ServiceName                 string                     `yaml:"serviceName"`
	MatchingComputeCapabilities []matchingComputeCapabilty `yaml:"matchingComputeCapabilities"`
	Applications                []application              `yaml:"applications"`
}

type matchingComputeCapabilty struct {
	CapabilityType  string   `yaml:"capabilityType"`
	Comparator      string   `yaml:"comparator"`
	CapabilityNames []string `yaml:"capabilityNames"`
}

type application struct {
	PayloadName       string                        `yaml:"payloadName"`
	Resources         ApplicationResources          `yaml:"resources"`
	Ports             []ApplicationPort             `yaml:"ports"`
	PersistentVolumes []ApplicationPersistentVolume `yaml:"persistentVolumes"`
}

type ApplicationResources struct {
	Cores  uint32 `yaml:"cores"`
	Memory uint32 `yaml:"memory"`
}

type ApplicationPort struct {
	Protocol string `yaml:"protocol"`
	Port     uint32 `yaml:"port"`
}

type ApplicationPersistentVolume struct {
	StorageClass string `yaml:"storageClass"`
	Capacity     uint32 `yaml:"capacity"`
	MountPath    string `yaml:"mountPath"`
}

type deploymentPlanList struct {
	Plans []planCreateRequest `yaml:"plans"`
}

// newDeploymentPlanCreateCmd creates a new deployment plan create command
func newDeploymentPlanCreateCmd() *cobra.Command {
	o := createDeploymentPlanOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new deployment plan",
		RunE: func(cmd *cobra.Command, args []string) error {

			conn, err := grpc.NewClient("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}
			o.deploymentPlansClient = mrdspb.NewDeploymentPlansClient(conn)
			o.printer = printer.NewPrinter()
			return o.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&o.manifestFilePath, "manifest", "m", "", "Path to the deployment plan manifest file")
	cmd.MarkFlagRequired("manifest")
	return cmd
}

func (o *createDeploymentPlanOptions) Run(ctx context.Context) error {
	yamlFile, err := os.Open(o.manifestFilePath)
	if err != nil {
		return err
	}

	req := &deploymentPlanList{}
	err = yaml.NewDecoder(yamlFile).Decode(req)
	if err != nil {
		return err
	}

	createdPlans := make([]*mrdspb.DeploymentPlanRecord, 0, len(req.Plans))
	for _, plan := range req.Plans {
		computeCapabilities := make([]*mrdspb.MatchingComputeCapability, 0, len(plan.MatchingComputeCapabilities))
		for _, cc := range plan.MatchingComputeCapabilities {
			computeCapabilities = append(computeCapabilities, &mrdspb.MatchingComputeCapability{
				CapabilityType:  cc.CapabilityType,
				Comparator:      mrdspb.Comparator(mrdspb.Comparator_value[string(cc.Comparator)]),
				CapabilityNames: cc.CapabilityNames,
			})
		}

		applications := make([]*mrdspb.Application, 0, len(plan.Applications))
		for _, app := range plan.Applications {
			ports := make([]*mrdspb.ApplicationPort, 0, len(app.Ports))
			for _, p := range app.Ports {
				ports = append(ports, &mrdspb.ApplicationPort{
					Protocol: p.Protocol,
					Port:     p.Port,
				})
			}

			persistentVolumes := make([]*mrdspb.ApplicationPersistentVolume, 0, len(app.PersistentVolumes))
			for _, pv := range app.PersistentVolumes {
				persistentVolumes = append(persistentVolumes, &mrdspb.ApplicationPersistentVolume{
					StorageClass: pv.StorageClass,
					Capacity:     pv.Capacity,
					MountPath:    pv.MountPath,
				})
			}

			applications = append(applications, &mrdspb.Application{
				PayloadName: app.PayloadName,
				Resources: &mrdspb.ApplicationResources{
					Cores:  app.Resources.Cores,
					Memory: app.Resources.Memory,
				},
				Ports:             ports,
				PersistentVolumes: persistentVolumes,
			})
		}

		resp, err := o.deploymentPlansClient.Create(ctx, &mrdspb.CreateDeploymentPlanRequest{
			Name:                        plan.Name,
			Namespace:                   plan.Namespace,
			ServiceName:                 plan.ServiceName,
			MatchingComputeCapabilities: computeCapabilities,
			Applications:                applications,
		})
		createdPlans = append(createdPlans, resp.Record)

		if err != nil {
			return err
		}
	}
	o.printer.PrintSuccess("Deployment plan created successfully" + createdPlans[0].Name)

	displayPlans := make([]types.DisplayDeploymentPlan, 0, len(createdPlans))
	for _, p := range createdPlans {
		displayPlans = append(displayPlans, convertGRPCDeploymentPlanToDisplayDeploymentPlan(p))
	}
	o.printer.PrintDisplayDeploymentPlanList(displayPlans)

	return nil
}
