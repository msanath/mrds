package metainstance

import (
	"context"

	"github.com/google/uuid"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/metainstance/getter"
	"github.com/msanath/mrds/pkg/ctl/metainstance/printer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type swapInstanceOptions struct {
	metaInstanceName string

	client  mrdspb.MetaInstancesClient
	printer *printer.Printer
	getter  *getter.Getter
}

func newSwapInstanceCmd() *cobra.Command {
	o := swapInstanceOptions{}
	cmd := &cobra.Command{
		Use:   "swap-instance",
		Short: "Swap a runtim instance to another node",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.NewClient("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}
			o.client = mrdspb.NewMetaInstancesClient(conn)
			o.printer = printer.NewPrinter()
			o.getter = getter.NewGetter(conn)

			return o.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&o.metaInstanceName, "meta-instance-name", "", "Meta instance name")
	return cmd
}

func (o swapInstanceOptions) Run(ctx context.Context) error {
	// Get meta instance
	metaInstanceResp, err := o.client.GetByName(ctx, &mrdspb.GetMetaInstanceByNameRequest{
		Name: o.metaInstanceName,
	})
	if err != nil {
		return err
	}

	displayRecord, err := o.getter.GetDisplayMetaInstances(ctx, metaInstanceResp.Record)
	if err != nil {
		return err
	}
	o.printer.PrintDisplayMetaInstance(displayRecord)
	o.printer.PrintEmptyLine()
	if !o.printer.SeekConfirmation("Are you sure you want to stop this instance?") {
		o.printer.PrintWarning("Operation cancelled")
		return nil
	}

	updateResp, err := o.client.AddOperation(ctx, &mrdspb.AddOperationRequest{
		Metadata: metaInstanceResp.Record.Metadata,
		Operation: &mrdspb.Operation{
			Id:       uuid.New().String(),
			Type:     "SWAP",
			IntentId: "User-swap",
			Status: &mrdspb.OperationStatus{
				State:   mrdspb.OperationState_OperationState_PENDING_APPROVAL,
				Message: "user requested swap",
			},
		},
	})
	if err != nil {
		return err
	}

	o.printer.PrintSuccess("Request sent to swap")
	o.printer.PrintEmptyLine()
	displayRecord, err = o.getter.GetDisplayMetaInstances(ctx, updateResp.Record)
	if err != nil {
		return err
	}
	o.printer.PrintDisplayMetaInstance(displayRecord)
	return nil
}
