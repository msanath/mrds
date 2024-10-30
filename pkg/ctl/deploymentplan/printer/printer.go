package printer

import (
	"strconv"
	"strings"

	"github.com/msanath/gondolf/pkg/printer"
	"github.com/msanath/mrds/pkg/ctl/deploymentplan/types"
)

type Printer struct {
	printer.PlainText
}

func NewPrinter() *Printer {
	return &Printer{
		PlainText: printer.NewPlainTextPrinter(),
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
			p.PrintKeyValueWithIndent("Payload Name", app.PayloadName)
			p.PrintKeyValueWithIndent("Cores", strconv.Itoa(app.Resources.Cores))
			p.PrintKeyValueWithIndent("Memory (MB)", strconv.Itoa(app.Resources.Memory))
			p.PrintHeader("Ports")
			if len(app.Ports) == 0 {
				p.PrintWarning("No ports found")
			} else {
				tableHeaders := []string{"Protocol", "Port"}
				rows := make([][]string, 0)
				for _, port := range app.Ports {
					rows = append(rows, []string{port.Protocol, strconv.Itoa(port.Port)})
				}
				p.PrintTable(tableHeaders, rows)
			}

			p.PrintHeader("Persistent Volumes")
			if len(app.PersistentVolumes) == 0 {
				p.PrintWarning("No persistent volumes found")
			} else {
				tableHeaders := []string{"Storage Class", "Capacity (GB)", "Mount Path"}
				rows := make([][]string, 0)
				for _, volume := range app.PersistentVolumes {
					rows = append(rows, []string{volume.StorageClass, strconv.Itoa(volume.Capacity), volume.MountPath})
				}
				p.PrintTable(tableHeaders, rows)
			}
			p.PrintLineSeparator()
		}
	}
	p.PrintEmptyLine()

	p.PrintHeader("Deployments")
	if len(plan.Deployments) == 0 {
		p.PrintWarning("No deployments found")
	} else {
		tableHeaders := []string{"Deployment ID", "Instance Count", "State", "Message"}
		rows := make([][]string, 0)
		for _, deployment := range plan.Deployments {
			rows = append(rows,
				[]string{
					deployment.ID,
					strconv.Itoa(deployment.InstanceCount),
					deployment.Status.State,
					deployment.Status.Message,
				},
			)
		}
		p.PrintTable(tableHeaders, rows)
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
