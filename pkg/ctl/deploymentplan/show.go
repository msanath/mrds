package deploymentplan

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/deploymentplan/printer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type deploymentPlanShowOptions struct {
	name string

	deploymentPlanClient mrdspb.DeploymentPlansClient
	printer              *printer.Printer
}

func newDeploymentPlanShowCmd() *cobra.Command {
	o := deploymentPlanShowOptions{}
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show deployment plan by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}
			o.name = args[0]
			o.deploymentPlanClient = mrdspb.NewDeploymentPlansClient(conn)
			o.printer = printer.NewPrinter()
			return o.Run(cmd.Context())
		},
	}

	return cmd
}

func (o *deploymentPlanShowOptions) Run(ctx context.Context) error {
	resp, err := o.deploymentPlanClient.GetByName(ctx, &mrdspb.GetDeploymentPlanByNameRequest{Name: o.name})
	if err != nil {
		return err
	}
	displayDeploymentPlan := convertGRPCDeploymentPlanToDisplayDeploymentPlan(resp.Record)
	o.printer.PrintDisplayDeploymentPlan(displayDeploymentPlan)
	return nil
}
