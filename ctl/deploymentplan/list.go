package deploymentplan

import (
	"context"

	"github.com/msanath/mrds/ctl/deploymentplan/printer"
	"github.com/msanath/mrds/ctl/deploymentplan/types"
	"github.com/msanath/mrds/gen/api/mrdspb"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type deploymentPlanListOptions struct {
	deploymentPlanClient mrdspb.DeploymentPlansClient
	printer              *printer.Printer
}

func newDeploymentPlanListCmd() *cobra.Command {
	o := deploymentPlanListOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all deployment plans",
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

	return cmd
}

func (o *deploymentPlanListOptions) Run(ctx context.Context) error {
	resp, err := o.deploymentPlanClient.List(ctx, &mrdspb.ListDeploymentPlanRequest{})
	if err != nil {
		return err
	}

	displayDeploymentPlans := make([]types.DisplayDeploymentPlan, 0, len(resp.Records))
	for _, d := range resp.Records {
		displayDeploymentPlans = append(displayDeploymentPlans, convertGRPCDeploymentPlanToDisplayDeploymentPlan(d))
	}

	if len(displayDeploymentPlans) == 0 {
		o.printer.PrintWarning("No deployment plans found")
		return nil
	}

	o.printer.PrintDisplayDeploymentPlanList(displayDeploymentPlans)
	return nil
}
