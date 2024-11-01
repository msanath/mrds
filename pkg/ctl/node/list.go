package node

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/node/printer"
	"github.com/msanath/mrds/pkg/ctl/node/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type nodeListOptions struct {
	payloadNamesNotIn  []string
	remainingCoresGTE  int
	remainingMemoryGTE int

	nodesClient mrdspb.NodesClient
	printer     *printer.Printer
}

func newNodeListCmd() *cobra.Command {
	o := nodeListOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.NewClient("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}
			o.nodesClient = mrdspb.NewNodesClient(conn)
			o.printer = printer.NewPrinter()
			return o.Run(cmd.Context())
		},
	}

	cmd.Flags().StringSliceVar(&o.payloadNamesNotIn, "payload-names-not-in", nil, "Filter nodes by payload names not in")
	cmd.Flags().IntVar(&o.remainingCoresGTE, "remaining-cores-gte", 0, "Filter nodes by remaining cores greater than or equal to")
	cmd.Flags().IntVar(&o.remainingMemoryGTE, "remaining-memory-gte", 0, "Filter nodes by remaining memory greater than or equal to")

	return cmd
}

func (o *nodeListOptions) Run(ctx context.Context) error {
	resp, err := o.nodesClient.List(context.Background(), &mrdspb.ListNodeRequest{
		PayloadNameNotIn:   o.payloadNamesNotIn,
		RemainingCoresGte:  uint32(o.remainingCoresGTE),
		RemainingMemoryGte: uint32(o.remainingMemoryGTE),
	})
	if err != nil {
		return err
	}
	displayNodes := make([]types.DisplayNode, 0, len(resp.Records))
	for _, n := range resp.Records {
		displayNodes = append(displayNodes, convertGRPCNodeToDisplayNode(n))
	}

	if len(displayNodes) == 0 {
		o.printer.PrintWarning("No nodes found")
		return nil
	}

	o.printer.PrintDisplayNodeList(displayNodes)
	return nil
}
