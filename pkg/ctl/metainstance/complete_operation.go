package metainstance

import (
	"context"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/msanath/mrds/pkg/ctl/metainstance/printer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type completeOperationOptions struct {
	metaInstanceName string
	operationId      string

	client  mrdspb.MetaInstancesClient
	printer *printer.Printer
}

func newCompleteOperationCmd() *cobra.Command {
	o := completeOperationOptions{}
	cmd := &cobra.Command{
		Use:   "complete-operation",
		Short: "Mark an operation as complete",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.NewClient("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}
			o.client = mrdspb.NewMetaInstancesClient(conn)
			o.printer = printer.NewPrinter()

			return o.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&o.metaInstanceName, "meta-instance-name", "", "Meta instance name")
	cmd.Flags().StringVar(&o.operationId, "operation-id", "", "Operation ID")

	return cmd
}

func (o completeOperationOptions) Run(ctx context.Context) error {
	// Get meta instance
	metaInstanceResp, err := o.client.GetByName(ctx, &mrdspb.GetMetaInstanceByNameRequest{
		Name: o.metaInstanceName,
	})
	if err != nil {
		return err
	}

	found := false
	for _, operation := range metaInstanceResp.Record.Operations {
		if operation.Id == o.operationId {
			if operation.Status.State != mrdspb.OperationState_OperationState_APPROVED {
				o.printer.PrintWarning("Operation is not approved, so it never took place to complete")
				return nil
			}
			found = true
			break
		}
	}
	if !found {
		o.printer.PrintError("Operation ID not found in meta instance")
	}

	o.printer.PrintDisplayMetaInstance(convertGRPCMetaInstanceToDisplayMetaInstance(metaInstanceResp.Record))
	o.printer.PrintEmptyLine()
	if !o.printer.SeekConfirmation("Are you sure you want to approve this operation?") {
		o.printer.PrintWarning("Operation cancelled")
		return nil
	}

	updateResp, err := o.client.UpdateOperationStatus(ctx, &mrdspb.UpdateOperationStatusRequest{
		Metadata:    metaInstanceResp.Record.Metadata,
		OperationId: o.operationId,
		Status: &mrdspb.OperationStatus{
			State:   mrdspb.OperationState_OperationState_SUCCEEDED,
			Message: "user approved operation",
		},
	})
	if err != nil {
		return err
	}

	o.printer.PrintSuccess("Request sent to swap")
	o.printer.PrintDisplayMetaInstance(convertGRPCMetaInstanceToDisplayMetaInstance(updateResp.Record))
	return nil
}
