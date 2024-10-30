package metainstance

import (
	"context"
	"fmt"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/metainstance/printer"
	"github.com/msanath/mrds/pkg/ctl/metainstance/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type metaInstanceListOptions struct {
	client  mrdspb.MetaInstancesClient
	printer *printer.Printer
}

func newMetaInstanceListCmd() *cobra.Command {
	o := metaInstanceListOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all meta instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.Dial("localhost:12345", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return fmt.Errorf("failed to connect to gRPC server: %w", err)
			}
			defer conn.Close()

			o.client = mrdspb.NewMetaInstancesClient(conn)
			o.printer = printer.NewPrinter()
			return o.Run(cmd.Context())
		},
	}

	return cmd
}

func (o *metaInstanceListOptions) Run(ctx context.Context) error {
	// Call gRPC service to list meta instances
	resp, err := o.client.List(ctx, &mrdspb.ListMetaInstanceRequest{})
	if err != nil {
		return fmt.Errorf("failed to list meta instances: %w", err)
	}

	// Convert gRPC MetaInstance records to display format
	displayMetaInstances := make([]types.DisplayMetaInstance, 0, len(resp.Records))
	for _, record := range resp.Records {
		displayMetaInstances = append(displayMetaInstances, convertGRPCMetaInstanceToDisplayMetaInstance(record))
	}

	// Print list of meta instances
	if len(displayMetaInstances) == 0 {
		o.printer.PrintWarning("No meta instances found")
		return nil
	}

	o.printer.PrintDisplayMetaInstanceList(displayMetaInstances)
	return nil
}
