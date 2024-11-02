package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/msanath/gondolf/pkg/printer"
	"github.com/msanath/mrds/pkg/ctl/deploymentplan/types"
	metainstanceprinter "github.com/msanath/mrds/pkg/ctl/metainstance/printer"
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

	if len(plan.InstanceSummary.MetaInstances) > 0 {
		p.PrintEmptyLine()
		p.PrintHeader("Instance Summary")
		p.metaInstancePrinter.PrintDisplayMetaInstanceList(plan.InstanceSummary.MetaInstances)
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
