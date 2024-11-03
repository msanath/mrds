package deploymentplan

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/msanath/gondolf/pkg/printer"
	deploymentgetter "github.com/msanath/mrds/ctl/deploymentplan/getter"
	deploymentprinter "github.com/msanath/mrds/ctl/deploymentplan/printer"
	metainstancegetter "github.com/msanath/mrds/ctl/metainstance/getter"
	metainstanceprinter "github.com/msanath/mrds/ctl/metainstance/printer"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type approveOperationOptions struct {
	deploymentPlanName  string
	metaInstanceName    string
	operationId         string
	deploymentsClient   mrdspb.DeploymentPlansClient
	metaInstanceClient  mrdspb.MetaInstancesClient
	deploymentPrinter   *deploymentprinter.Printer
	metaInstancePrinter *metainstanceprinter.Printer
	deploymentGetter    *deploymentgetter.Getter
	metaInstanceGetter  *metainstancegetter.Getter
}

func newApproveOperationCmd() *cobra.Command {
	o := approveOperationOptions{}
	cmd := &cobra.Command{
		Use:   "approve-operation",
		Short: "Approve a pending operation.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.NewClient("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}
			o.deploymentsClient = mrdspb.NewDeploymentPlansClient(conn)
			o.metaInstanceClient = mrdspb.NewMetaInstancesClient(conn)
			o.deploymentPrinter = deploymentprinter.NewPrinter()
			o.metaInstancePrinter = metainstanceprinter.NewPrinter()
			o.deploymentGetter = deploymentgetter.NewGetter(conn)
			o.metaInstanceGetter = metainstancegetter.NewGetter(conn)

			o.deploymentPlanName = args[0]

			return o.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&o.metaInstanceName, "meta-instance-name", "", "Meta instance name")
	cmd.Flags().StringVar(&o.operationId, "operation-id", "", "Operation ID")

	return cmd
}

func (o approveOperationOptions) Run(ctx context.Context) error {
	// Get deployment plan
	deploymentPlanResp, err := o.deploymentsClient.GetByName(ctx, &mrdspb.GetDeploymentPlanByNameRequest{
		Name: o.deploymentPlanName,
	})
	if err != nil {
		return err
	}
	displayPlan, err := o.deploymentGetter.ConvertGRPCDeploymentPlanToDisplayDeploymentPlan(ctx, deploymentPlanResp.Record)
	if err != nil {
		return err
	}
	if o.metaInstanceName == "" || o.operationId == "" {
		opSelectors := make([]operationSelector, 0)
		opSelectorStr := make([]string, 0)

		for _, metaInstance := range displayPlan.InstanceSummary.MetaInstances {
			for _, operation := range metaInstance.Operations {
				if operation.Status.State == mrdspb.OperationState_OperationState_PENDING_APPROVAL.String() {
					nodeName := ""
					for _, instance := range metaInstance.RuntimeInstances {
						if instance.IsActive {
							nodeName = instance.NodeName
							break
						}
					}
					opSelector := operationSelector{
						metaInstanceName: metaInstance.Name,
						nodeName:         nodeName,
						operationType:    operation.Type,
						operationID:      operation.ID,
					}
					opSelectorStr = append(opSelectorStr, opSelector.String())
					opSelectors = append(opSelectors, opSelector)
				}
			}
		}

		if len(opSelectorStr) == 0 {
			o.deploymentPrinter.PrintWarning("No operations pending approval")
			return nil
		}

		prompt := promptui.Select{
			Label: "Select Operation",
			Items: opSelectorStr,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return err
		}

		selectedOp := opSelectors[idx]
		o.metaInstanceName = selectedOp.metaInstanceName
		o.operationId = selectedOp.operationID
	}

	// Get meta instance
	metaInstanceResp, err := o.metaInstanceClient.GetByName(ctx, &mrdspb.GetMetaInstanceByNameRequest{
		Name: o.metaInstanceName,
	})
	if err != nil {
		return err
	}

	found := false
	for _, operation := range metaInstanceResp.Record.Operations {
		if operation.Id == o.operationId {
			if operation.Status.State != mrdspb.OperationState_OperationState_PENDING_APPROVAL {
				o.deploymentPrinter.PrintWarning("Operation is not pending approval")
				return nil
			}
			found = true
			break
		}
	}
	if !found {
		o.deploymentPrinter.PrintError("Operation ID not found in meta instance")
	}

	displayRecord, err := o.metaInstanceGetter.GetDisplayMetaInstance(ctx, metaInstanceResp.Record)
	if err != nil {
		return err
	}

	o.metaInstancePrinter.PrintDisplayMetaInstance(displayRecord)
	o.deploymentPrinter.PrintEmptyLine()
	if !o.deploymentPrinter.SeekConfirmation("Are you sure you want to approve this operation?") {
		o.deploymentPrinter.PrintWarning("Operation cancelled")
		return nil
	}

	updateResp, err := o.metaInstanceClient.UpdateOperationStatus(ctx, &mrdspb.UpdateOperationStatusRequest{
		Metadata:    metaInstanceResp.Record.Metadata,
		OperationId: o.operationId,
		Status: &mrdspb.OperationStatus{
			State:   mrdspb.OperationState_OperationState_APPROVED,
			Message: "user approved operation",
		},
	})
	if err != nil {
		return err
	}

	o.deploymentPrinter.PrintSuccess("Request sent to swap")
	o.deploymentPrinter.PrintEmptyLine()
	displayRecord, err = o.metaInstanceGetter.GetDisplayMetaInstance(ctx, updateResp.Record)
	if err != nil {
		return err
	}
	o.metaInstancePrinter.PrintDisplayMetaInstance(displayRecord)
	return nil
}

type operationSelector struct {
	metaInstanceName string
	nodeName         string
	operationID      string
	operationType    string
}

func (os operationSelector) String() string {
	return fmt.Sprintf(
		"MetaInstanceID: %s - OpetationType: %s - NodeName: %s - OperationID: %s",
		printer.GreenText(os.metaInstanceName),
		printer.GreenText(os.operationType),
		printer.GreenText(os.nodeName),
		printer.GreenText(os.operationID),
	)
}
