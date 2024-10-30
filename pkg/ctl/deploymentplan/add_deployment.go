package deploymentplan

import (
	"context"
	"os"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/deploymentplan/printer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

type addDeploymenOptions struct {
	manifestFilePath string

	deploymentPlanClient mrdspb.DeploymentPlansClient
	printer              *printer.Printer
}

type deploymentCreateRequest struct {
	DeploymentID       string               `yaml:"deployment_id"`
	DeploymentPlanName string               `yaml:"deployment_plan_name"`
	PayloadCoordinates []payloadCoordinates `yaml:"payload_coordinates"`
	InstanceCount      uint32               `yaml:"instance_count"`
}

type payloadCoordinates struct {
	PayloadName string            `yaml:"payload_name"`
	Coordinates map[string]string `yaml:"coordinates"`
}

func newAddDeploymentCmd() *cobra.Command {
	o := addDeploymenOptions{}
	cmd := &cobra.Command{
		Use:   "add-deployment",
		Short: "Add a deployment",
		RunE: func(cmd *cobra.Command, args []string) error {

			conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}
			o.deploymentPlanClient = mrdspb.NewDeploymentPlansClient(conn)
			o.printer = printer.NewPrinter()
			return o.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&o.manifestFilePath, "manifest", "m", "", "Path to the deployment plan manifest file")
	cmd.MarkFlagRequired("manifest")
	return cmd
}

func (o *addDeploymenOptions) Run(ctx context.Context) error {
	yamlFile, err := os.Open(o.manifestFilePath)
	if err != nil {
		return err
	}

	req := &deploymentCreateRequest{}
	err = yaml.NewDecoder(yamlFile).Decode(req)
	if err != nil {
		return err
	}

	// Get node by name
	getResp, err := o.deploymentPlanClient.GetByName(ctx, &mrdspb.GetDeploymentPlanByNameRequest{Name: req.DeploymentPlanName})
	if err != nil {
		return err
	}

	payloadCoordinatesProto := make([]*mrdspb.PayloadCoordinates, 0)
	for _, pc := range req.PayloadCoordinates {
		coordinates := make(map[string]string)
		for k, v := range pc.Coordinates {
			coordinates[k] = v
		}
		payloadCoordinatesProto = append(payloadCoordinatesProto, &mrdspb.PayloadCoordinates{
			PayloadName: pc.PayloadName,
			Coordinates: coordinates,
		})
	}

	updateResp, err := o.deploymentPlanClient.AddDeployment(ctx, &mrdspb.AddDeploymentRequest{
		Metadata:           getResp.Record.GetMetadata(),
		DeploymentId:       req.DeploymentID,
		PayloadCoordinates: payloadCoordinatesProto,
		InstanceCount:      req.InstanceCount,
	})
	if err != nil {
		return err
	}

	o.printer.PrintDisplayDeploymentPlan(convertGRPCDeploymentPlanToDisplayDeploymentPlan(updateResp.Record))
	return nil
}
