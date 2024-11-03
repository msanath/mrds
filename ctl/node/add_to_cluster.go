package node

import (
	"context"

	"github.com/msanath/mrds/ctl/node/printer"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type addToClusterOptions struct {
	name      string
	clusterID string

	nodesClient mrdspb.NodesClient
	printer     *printer.Printer
}

func newAddToClusterCmd() *cobra.Command {
	o := addToClusterOptions{}
	cmd := &cobra.Command{
		Use:   "add-to-cluster",
		Short: "Add node to cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.name = args[0]

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

	cmd.Flags().StringVarP(&o.clusterID, "cluster-id", "c", "", "ID of the cluster")

	return cmd
}

func (o *addToClusterOptions) Run(ctx context.Context) error {
	// Get node by name
	getResp, err := o.nodesClient.GetByName(ctx, &mrdspb.GetNodeByNameRequest{Name: o.name})
	if err != nil {
		return err
	}

	updateResp, err := o.nodesClient.UpdateStatus(ctx, &mrdspb.UpdateNodeStatusRequest{
		Metadata: getResp.Record.Metadata,
		Status: &mrdspb.NodeStatus{
			State: mrdspb.NodeState_NodeState_ALLOCATING,
		},
		ClusterId: o.clusterID,
	})
	if err != nil {
		return err
	}
	o.printer.PrintDisplayNode(convertGRPCNodeToDisplayNode(updateResp.Record))
	updateResp, err = o.nodesClient.UpdateStatus(ctx, &mrdspb.UpdateNodeStatusRequest{
		Metadata: updateResp.Record.Metadata,
		Status: &mrdspb.NodeStatus{
			State: mrdspb.NodeState_NodeState_ALLOCATED,
		},
		ClusterId: o.clusterID,
	})
	if err != nil {
		return err
	}

	o.printer.PrintSuccess("Node added to cluster")
	o.printer.PrintDisplayNode(convertGRPCNodeToDisplayNode(updateResp.Record))
	return nil
}
