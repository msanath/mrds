package node

import (
	"context"

	"github.com/msanath/mrds/ctl/node/printer"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type nodeShowOptions struct {
	name string

	nodesClient mrdspb.NodesClient
	printer     *printer.Printer
}

func newNodeShowCmd() *cobra.Command {
	o := nodeShowOptions{}
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show node by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.NewClient("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}
			o.name = args[0]
			o.nodesClient = mrdspb.NewNodesClient(conn)
			o.printer = printer.NewPrinter()
			return o.Run(cmd.Context())
		},
	}

	return cmd
}

func (o *nodeShowOptions) Run(ctx context.Context) error {
	resp, err := o.nodesClient.GetByName(ctx, &mrdspb.GetNodeByNameRequest{Name: o.name})
	if err != nil {
		return err
	}
	displayNode := convertGRPCNodeToDisplayNode(resp.Record)
	o.printer.PrintDisplayNode(displayNode)
	return nil
}
