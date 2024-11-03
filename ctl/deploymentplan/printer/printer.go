package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/msanath/gondolf/pkg/printer"
	"github.com/msanath/mrds/ctl/deploymentplan/types"
	metainstanceprinter "github.com/msanath/mrds/ctl/metainstance/printer"
	"github.com/msanath/mrds/gen/api/mrdspb"
)

type Printer struct {
	printer.PlainText
	metaInstancePrinter *metainstanceprinter.Printer
}

func NewPrinter() *Printer {
	return &Printer{
		PlainText:           printer.NewPlainTextPrinter(),
		metaInstancePrinter: metainstanceprinter.NewPrinter(),
	}
}

func (p *Printer) PrintDisplayDeploymentPlan(plan types.DisplayDeploymentPlan) {
	p.PrintHeader("Deployment Plan Info")

	p.PrintDisplayField(plan.GetName())
	p.PrintDisplayField(plan.Metadata.GetID())
	p.PrintDisplayField(plan.Metadata.GetVersion())
	p.PrintDisplayField(plan.GetNamespace())
	p.PrintDisplayField(plan.GetServiceName())
	p.PrintEmptyLine()

	p.PrintHeader("Status")
	p.PrintDisplayFieldWithIndent(plan.Status.GetState())
	p.PrintDisplayFieldWithIndent(plan.Status.GetMessage())
	p.PrintEmptyLine()

	p.PrintHeader("Matching Compute Capabilities")
	if len(plan.MatchingComputeCapabilities) == 0 {
		p.PrintWarning("No matching compute capabilities found")
	} else {
		tableHeaders := []string{"Capability Type", "Comparator", "Capability Names"}
		rows := make([][]string, 0)
		for _, capability := range plan.MatchingComputeCapabilities {
			rows = append(rows,
				[]string{
					capability.CapabilityType,
					capability.Comparator,
					strings.Join(capability.CapabilityNames, ","),
				},
			)
		}
		p.PrintTable(tableHeaders, rows)
	}
	p.PrintEmptyLine()

	p.PrintHeader("Applications")
	if len(plan.Applications) == 0 {
		p.PrintWarning("No applications found")
	} else {
		for _, app := range plan.Applications {

			ports := []string{}
			for _, port := range app.Ports {
				ports = append(ports, port.Protocol+":"+strconv.Itoa(port.Port))
			}

			persistentVolumes := []string{}
			for _, volume := range app.PersistentVolumes {
				persistentVolumes = append(persistentVolumes, fmt.Sprintf("Storage Class: %s, Capacity: %d GB, Mount Path: %s", volume.StorageClass, volume.Capacity, volume.MountPath))
			}

			tableHeaders := []string{"Payload Name", "Cores", "Memory (GiB)", "Ports", "Persistent Volumes"}
			rows := make([][]string, 0)
			rows = append(rows,
				[]string{
					app.GetPayloadName().Value(),
					app.Resources.GetCores().Value(),
					app.Resources.GetMemory().Value(),
					strings.Join(ports, "\n"),
					strings.Join(persistentVolumes, "\n"),
				},
			)
			p.PrintTable(tableHeaders, rows, printer.WithRowSeparator())
		}
	}
	p.PrintEmptyLine()

	p.PrintHeader("Deployments")
	if len(plan.Deployments) == 0 {
		p.PrintWarning("No deployments found")
	} else {
		tableHeaders := []string{"Deployment ID", "Instance Count", "State", "Message", "Payload Information"}
		rows := make([][]string, 0)
		for _, deployment := range plan.Deployments {

			payloadInfo := []string{}
			for _, payloadCoordinate := range deployment.PayloadCoordinates {
				payloadInfo = append(payloadInfo, fmt.Sprintf("Payload Name: %s, Coordinates: %s", payloadCoordinate.PayloadName, payloadCoordinate.Coordinates))
			}

			rows = append(rows,
				[]string{
					deployment.GetID().Value(),
					deployment.GetInstanceCount().Value(),
					deployment.Status.GetState().Value(),
					deployment.Status.GetMessage().Value(),
					strings.Join(payloadInfo, "\n"),
				},
			)
		}
		p.PrintTable(tableHeaders, rows)
	}

	p.PrintEmptyLine()
	p.PrintHeader("Instance Summary")
	if len(plan.InstanceSummary.MetaInstances) > 0 {
		p.PrintDisplayFieldWithIndent(plan.InstanceSummary.GetNumTotalInstances())
		p.PrintDisplayFieldWithIndent(plan.InstanceSummary.GetNumRunningInstances())
		p.PrintDisplayFieldWithIndent(plan.InstanceSummary.GetNumPendingInstances())
		p.PrintDisplayFieldWithIndent(plan.InstanceSummary.GetNumFailedInstances())
		p.PrintEmptyLine()

		p.PrintHeader("Opeartions Summary")
		tableHeaders := []string{"Operation Type", "# Pending Approval", "# Approved", "# Failed", "# Succeeded"}
		rows := [][]string{
			{
				mrdspb.OperationType_OperationType_CREATE.String(),
				plan.InstanceSummary.GetNumCreateOperationsPendingApproval().Value(),
				plan.InstanceSummary.GetNumCreateOperationsApproved().Value(),
				plan.InstanceSummary.GetNumCreateOperationsFailed().Value(),
				plan.InstanceSummary.GetNumCreateOperationsSucceeded().Value(),
			},
			{
				mrdspb.OperationType_OperationType_UPDATE.String(),
				plan.InstanceSummary.GetNumUpdateOperationsPendingApproval().Value(),
				plan.InstanceSummary.GetNumUpdateOperationsApproved().Value(),
				plan.InstanceSummary.GetNumUpdateOperationsFailed().Value(),
				plan.InstanceSummary.GetNumUpdateOperationsSucceeded().Value(),
			},
			{
				mrdspb.OperationType_OperationType_RELOCATE.String(),
				plan.InstanceSummary.GetNumRelocateOperationsPendingApproval().Value(),
				plan.InstanceSummary.GetNumRelocateOperationsApproved().Value(),
				plan.InstanceSummary.GetNumRelocateOperationsFailed().Value(),
				plan.InstanceSummary.GetNumRelocateOperationsSucceeded().Value(),
			},
			{
				mrdspb.OperationType_OperationType_STOP.String(),
				plan.InstanceSummary.GetNumStopOperationsPendingApproval().Value(),
				plan.InstanceSummary.GetNumStopOperationsApproved().Value(),
				plan.InstanceSummary.GetNumStopOperationsFailed().Value(),
				plan.InstanceSummary.GetNumStopOperationsSucceeded().Value(),
			},
			{
				mrdspb.OperationType_OperationType_RESTART.String(),
				plan.InstanceSummary.GetNumRestartOperationsPendingApproval().Value(),
				plan.InstanceSummary.GetNumRestartOperationsApproved().Value(),
				plan.InstanceSummary.GetNumRestartOperationsFailed().Value(),
				plan.InstanceSummary.GetNumRestartOperationsSucceeded().Value(),
			},
			{
				mrdspb.OperationType_OperationType_DELETE.String(),
				plan.InstanceSummary.GetNumDeleteOperationsPendingApproval().Value(),
				plan.InstanceSummary.GetNumDeleteOperationsApproved().Value(),
				plan.InstanceSummary.GetNumDeleteOperationsFailed().Value(),
				plan.InstanceSummary.GetNumDeleteOperationsSucceeded().Value(),
			},
		}
		p.PrintTable(tableHeaders, rows, printer.WithHideTotal())
		p.PrintEmptyLine()

		p.PrintHeader("Instances")
		p.metaInstancePrinter.PrintDisplayMetaInstanceList(plan.InstanceSummary.MetaInstances)
	} else {
		p.PrintWarning("No instances found")
	}
}

func (p *Printer) PrintDisplayDeploymentPlanList(plans []types.DisplayDeploymentPlan) {
	tableHeaders := []string{
		"Deployment Plan Name",
		"Namespace",
		"Service Name",
		"State",
		"Message",
		"# Applications",
		"# Deployments",
	}
	rows := make([][]string, 0)
	for _, plan := range plans {
		rows = append(rows,
			[]string{
				plan.GetName().Value(),
				plan.GetNamespace().Value(),
				plan.GetServiceName().Value(),
				plan.Status.GetState().Value(),
				plan.Status.GetMessage().Value(),
				strconv.Itoa(len(plan.Applications)),
				strconv.Itoa(len(plan.Deployments)),
			},
		)
	}
	p.PrintTable(tableHeaders, rows)
}
