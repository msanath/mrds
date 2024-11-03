package deploymentplan

import (
	"context"

	"github.com/msanath/mrds/ctl/deploymentplan/printer"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type cancelDeploymentOptions struct {
	deploymentPlanName string
	deploymentID       string

	deploymentPlanClient mrdspb.DeploymentPlansClient
	printer              *printer.Printer
}

func newCancelDeploymentCmd() *cobra.Command {
	o := cancelDeploymentOptions{}
	cmd := &cobra.Command{
		Use:   "cancel-deployment",
		Short: "Cancel a deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}

			o.deploymentPlanClient = mrdspb.NewDeploymentPlansClient(conn)
			o.printer = printer.NewPrinter()
			o.deploymentPlanName = args[0]
			return o.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&o.deploymentID, "deployment-id", "D", "", "ID of the deployment")
	cmd.MarkFlagRequired("deployment-id")
	return cmd
}

func (o *cancelDeploymentOptions) Run(ctx context.Context) error {
	// Get deployment by name
	getResp, err := o.deploymentPlanClient.GetByName(ctx, &mrdspb.GetDeploymentPlanByNameRequest{
		Name: o.deploymentPlanName,
	})
	if err != nil {
		return err
	}

	plan := getResp.Record
	foundDeployment := false
	for _, d := range plan.Deployments {
		if d.Id == o.deploymentID {
			foundDeployment = true
			if d.Status.State == mrdspb.DeploymentState_DeploymentState_COMPLETED {
				o.printer.PrintWarning("Deployment is already completed")
				return nil
			}
			if !o.printer.SeekConfirmation("Are you sure you want to cancel the deployment?") {
				o.printer.PrintWarning("Operation canceled")
				return nil
			}
		}
	}
	if !foundDeployment {
		o.printer.PrintWarning("Deployment not found")
		return nil
	}

	updateResp, err := o.deploymentPlanClient.UpdateDeploymentStatus(ctx, &mrdspb.UpdateDeploymentStatusRequest{
		Metadata: plan.Metadata,
		Status: &mrdspb.DeploymentStatus{
			State:   mrdspb.DeploymentState_DeploymentState_CANCELLED,
			Message: "User canceled",
		},
		DeploymentId: o.deploymentID,
	})
	if err != nil {
		return err
	}
	o.printer.PrintSuccess("Deployment canceled")
	o.printer.PrintDisplayDeploymentPlan(convertGRPCDeploymentPlanToDisplayDeploymentPlan(updateResp.Record))
	return nil
}
