package deploymentplan

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

type operateOptions struct {
	deploymentPlanName string
	metaInstanceName   string

	restart  bool
	stop     bool
	relocate bool

	deploymentsClient   mrdspb.DeploymentPlansClient
	metaInstanceClient  mrdspb.MetaInstancesClient
	deploymentPrinter   *deploymentprinter.Printer
	metaInstancePrinter *metainstanceprinter.Printer
	deploymentGetter    *deploymentgetter.Getter
	metaInstanceGetter  *metainstancegetter.Getter
}

func newOperateOption() *cobra.Command {
	o := operateOptions{}
	cmd := &cobra.Command{
		Use:   "operate",
		Short: "Operate on the plan and instances",
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
	cmd.Flags().BoolVar(&o.restart, "start", false, "Start the instance")
	cmd.Flags().BoolVar(&o.stop, "stop", false, "Stop the instance")
	cmd.Flags().BoolVar(&o.relocate, "relocate", false, "Relocate the instance")
	cmd.MarkFlagsMutuallyExclusive("start", "stop", "relocate")
	cmd.MarkFlagsOneRequired("start", "stop", "relocate")

	return cmd
}

func (o operateOptions) Run(ctx context.Context) error {

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
	if o.metaInstanceName == "" {
		selectors := make([]instanceSelector, 0)
		selectorStr := make([]string, 0)

		for _, metaInstance := range displayPlan.InstanceSummary.MetaInstances {
			nodeName := ""
			instanceState := ""
			for _, instance := range metaInstance.RuntimeInstances {
				if instance.IsActive {
					nodeName = instance.NodeName
					instanceState = instance.Status.GetState().Value()
					break
				}
			}
			opSelector := instanceSelector{
				metaInstanceName:    metaInstance.Name,
				nodeName:            nodeName,
				ActiveInstanceState: instanceState,
			}
			selectorStr = append(selectorStr, opSelector.String())
			selectors = append(selectors, opSelector)
		}

		if len(selectorStr) == 0 {
			o.deploymentPrinter.PrintWarning("No instances found.")
			return nil
		}

		prompt := promptui.Select{
			Label: "Select the instance to operate",
			Items: selectorStr,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return err
		}

		selectedOp := selectors[idx]
		o.metaInstanceName = selectedOp.metaInstanceName
	}

	// Get meta instance
	metaInstanceResp, err := o.metaInstanceClient.GetByName(ctx, &mrdspb.GetMetaInstanceByNameRequest{
		Name: o.metaInstanceName,
	})
	if err != nil {
		return err
	}

	displayRecord, err := o.metaInstanceGetter.GetDisplayMetaInstance(ctx, metaInstanceResp.Record)
	if err != nil {
		return err
	}
	o.metaInstancePrinter.PrintDisplayMetaInstance(displayRecord)
	o.deploymentPrinter.PrintEmptyLine()
	if !o.deploymentPrinter.SeekConfirmation("Are you sure you want to continue?") {
		o.deploymentPrinter.PrintWarning("Operation cancelled")
		return nil
	}

	var operationType mrdspb.OperationType
	if o.restart {
		operationType = mrdspb.OperationType_OperationType_RESTART
	} else if o.stop {
		operationType = mrdspb.OperationType_OperationType_STOP
	} else if o.relocate {
		operationType = mrdspb.OperationType_OperationType_RELOCATE
	}

	// Create the operation
	updateResp, err := o.metaInstanceClient.AddOperation(ctx, &mrdspb.AddOperationRequest{
		Metadata: metaInstanceResp.Record.Metadata,
		Operation: &mrdspb.Operation{
			Id:       fmt.Sprintf("%s-%s", operationType.String(), uuid.New().String()),
			Type:     operationType,
			IntentId: "User",
			Status: &mrdspb.OperationStatus{
				State:   mrdspb.OperationState_OperationState_PENDING,
				Message: "user initiated",
			},
		},
	})
	if err != nil {
		return err
	}
	o.deploymentPrinter.PrintSuccess("Request sent to stop")
	o.deploymentPrinter.PrintEmptyLine()

	displayRecord, err = o.metaInstanceGetter.GetDisplayMetaInstance(ctx, updateResp.Record)
	if err != nil {
		return err
	}
	o.metaInstancePrinter.PrintDisplayMetaInstance(displayRecord)
	return nil
}

type instanceSelector struct {
	metaInstanceName    string
	nodeName            string
	ActiveInstanceState string
}

func (is instanceSelector) String() string {
	return fmt.Sprintf(
		"MetaInstanceID: %s - NodeName: %s - ActiveInstanceState: %s",
		printer.GreenText(is.metaInstanceName),
		printer.GreenText(is.nodeName),
		printer.GreenText(is.ActiveInstanceState),
	)
}
