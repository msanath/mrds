package metainstance

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/metainstance/printer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type createInstanceOptions struct {
	deploymentPlanName string
	deploymentID       string

	metaInstancesClient   mrdspb.MetaInstancesClient
	deploymentPlansClient mrdspb.DeploymentPlansClient
	printer               *printer.Printer
}

func newCreateCmd() *cobra.Command {
	o := createInstanceOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.NewClient("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}
			o.metaInstancesClient = mrdspb.NewMetaInstancesClient(conn)
			o.deploymentPlansClient = mrdspb.NewDeploymentPlansClient(conn)
			o.printer = printer.NewPrinter()

			return o.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&o.deploymentPlanName, "deployment-plan", "", "Deployment plan name")
	cmd.Flags().StringVar(&o.deploymentID, "deployment-id", "", "Deployment ID")

	return cmd
}

func (o createInstanceOptions) Run(ctx context.Context) error {
	// Get deployment plan
	deploymentResp, err := o.deploymentPlansClient.GetByName(ctx, &mrdspb.GetDeploymentPlanByNameRequest{
		Name: o.deploymentPlanName,
	})
	if err != nil {
		return err
	}

	foundDeployment := false
	for _, deployments := range deploymentResp.Record.Deployments {
		if deployments.Id == o.deploymentID {
			foundDeployment = true
			break
		}
	}
	if !foundDeployment {
		o.printer.PrintError("Deployment ID not found in deployment plan")
	}

	// Create instance
	createResp, err := o.metaInstancesClient.Create(ctx, &mrdspb.CreateMetaInstanceRequest{
		Name:             fmt.Sprintf("instance-%s", uuid.New().String()),
		DeploymentPlanId: deploymentResp.Record.Metadata.Id,
		DeploymentId:     o.deploymentID,
	})
	if err != nil {
		return err
	}

	o.printer.PrintDisplayMetaInstance(convertGRPCMetaInstanceToDisplayMetaInstance(createResp.Record))
	return nil
}
