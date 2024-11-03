package metainstance

import (
	"context"

	"github.com/google/uuid"
	"github.com/msanath/mrds/ctl/metainstance/getter"
	"github.com/msanath/mrds/ctl/metainstance/printer"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type stopInstanceOptions struct {
	metaInstanceName string

	client  mrdspb.MetaInstancesClient
	printer *printer.Printer
	getter  *getter.Getter
}

func newStopInstanceCmd() *cobra.Command {
	o := stopInstanceOptions{}
	cmd := &cobra.Command{
		Use:   "stop-instance",
		Short: "Stop a meta instance",
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

func (o stopInstanceOptions) Run(ctx context.Context) error {
	// Get meta instance
	metaInstanceResp, err := o.client.GetByName(ctx, &mrdspb.GetMetaInstanceByNameRequest{
		Name: o.metaInstanceName,
	})
	if err != nil {
		return err
	}

	displayRecord, err := o.getter.GetDisplayMetaInstance(ctx, metaInstanceResp.Record)
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
			Type:     mrdspb.OperationType_OperationType_STOP,
			IntentId: "User-Stop",
			Status: &mrdspb.OperationStatus{
				State:   mrdspb.OperationState_OperationState_PENDING,
				Message: "user requested stop",
			},
		},
	})
	if err != nil {
		return err
	}

	o.printer.PrintSuccess("Request sent to stop")
	o.printer.PrintEmptyLine()
	displayRecord, err = o.getter.GetDisplayMetaInstance(ctx, updateResp.Record)
	if err != nil {
		return err
	}
	o.printer.PrintDisplayMetaInstance(displayRecord)
	return nil
}
