package metainstance

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/metainstance/getter"
	"github.com/msanath/mrds/pkg/ctl/metainstance/printer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type addRuntimeInstanceOptions struct {
	metaInstanceName string
	nodeName         string

	metaInstancesClient mrdspb.MetaInstancesClient
	nodesClient         mrdspb.NodesClient
	printer             *printer.Printer
	getter              *getter.Getter
}

func newAddRuntimeInstanceCmd() *cobra.Command {
	o := addRuntimeInstanceOptions{}
	cmd := &cobra.Command{
		Use:   "add-runtime-instance",
		Short: "Add a new runtime instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return fmt.Errorf("failed to connect to gRPC server: %w", err)
			}
			defer conn.Close()

			o.metaInstancesClient = mrdspb.NewMetaInstancesClient(conn)
			o.nodesClient = mrdspb.NewNodesClient(conn)
			o.printer = printer.NewPrinter()
			o.getter = getter.NewGetter(conn)

			return o.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&o.metaInstanceName, "meta-instance", "", "Meta instance name")
	cmd.Flags().StringVar(&o.nodeName, "node", "", "Node name")

	return cmd
}

func (o addRuntimeInstanceOptions) Run(ctx context.Context) error {
	// Get meta instance
	metaInstanceResp, err := o.metaInstancesClient.GetByName(ctx, &mrdspb.GetMetaInstanceByNameRequest{
		Name: o.metaInstanceName,
	})
	if err != nil {
		return err
	}

	// Get node
	nodeResp, err := o.nodesClient.GetByName(ctx, &mrdspb.GetNodeByNameRequest{
		Name: o.nodeName,
	})
	if err != nil {
		return err
	}

	// Add runtime instance
	updateResp, err := o.metaInstancesClient.AddRuntimeInstance(ctx, &mrdspb.AddRuntimeInstanceRequest{
		Metadata: metaInstanceResp.Record.Metadata,
		RuntimeInstance: &mrdspb.RuntimeInstance{
			Id:     fmt.Sprintf("RuntimeInstance-%s", uuid.New().String()),
			NodeId: nodeResp.Record.Metadata.Id,
			Status: &mrdspb.RuntimeInstanceStatus{
				State:   mrdspb.RuntimeInstanceState_RuntimeState_PENDING,
				Message: "",
			},
		},
	})
	if err != nil {
		return err
	}

	o.printer.PrintSuccess("Runtime instance added successfully")
	displayRecord, err := o.getter.GetDisplayMetaInstance(ctx, updateResp.Record)
	if err != nil {
		return err
	}
	o.printer.PrintDisplayMetaInstance(displayRecord)
	return nil
}
