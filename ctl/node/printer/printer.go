package printer

import (
	"strconv"

	"github.com/msanath/gondolf/pkg/printer"
	"github.com/msanath/mrds/ctl/node/types"
)

type Printer struct {
	printer.PlainText
}

func NewPrinter() *Printer {
	return &Printer{
		PlainText: printer.NewPlainTextPrinter(),
	}
}

func (p *Printer) PrintDisplayNode(node types.DisplayNode) {
	p.PrintHeader("Node Info")

	p.PrintDisplayField(node.GetName())
	p.PrintDisplayField(node.Metadata.GetID())
	p.PrintDisplayField(node.Metadata.GetVersion())
	p.PrintDisplayField(node.GetClusterID())
	p.PrintDisplayField(node.GetUpdateDomain())
	p.PrintEmptyLine()

	p.PrintHeader("Status")
	p.PrintDisplayFieldWithIndent(node.Status.GetState())
	p.PrintDisplayFieldWithIndent(node.Status.GetMessage())
	p.PrintEmptyLine()

	p.PrintHeader("Capabilities")
	p.PrintDisplayFieldWithIndent(node.GetCapabilityIDs())
	p.PrintEmptyLine()

	p.PrintHeader("Total Resources")
	p.PrintDisplayFieldWithIndent(node.TotalResources.GetCores())
	p.PrintDisplayFieldWithIndent(node.TotalResources.GetMemory())
	p.PrintEmptyLine()

	p.PrintHeader("System Reserved Resources")
	p.PrintDisplayFieldWithIndent(node.SystemReserved.GetCores())
	p.PrintDisplayFieldWithIndent(node.SystemReserved.GetMemory())
	p.PrintEmptyLine()

	p.PrintHeader("Remaining Resources")
	p.PrintDisplayFieldWithIndent(node.RemainingResources.GetCores())
	p.PrintDisplayFieldWithIndent(node.RemainingResources.GetMemory())
	p.PrintEmptyLine()

	p.PrintHeader("Local Volumes")
	if len(node.LocalVolumes) == 0 {
		p.PrintWarning("No local volumes found")
	} else {
		tableHeaders := []string{"Mount Path", "Storage Class", "Storage Capacity"}
		rows := make([][]string, 0)
		for _, volume := range node.LocalVolumes {
			rows = append(rows,
				[]string{volume.MountPath, volume.StorageClass, strconv.Itoa(volume.StorageCapacity)},
			)
		}
		p.PrintTable(tableHeaders, rows)
	}
	p.PrintEmptyLine()

	p.PrintHeader("Disruptions")
	if len(node.Disruptions) == 0 {
		p.PrintWarning("No disruptions found")
	} else {
		tableHeaders := []string{"Disruption ID", "Should Evict", "Start Time", "State", "Message"}
		rows := make([][]string, 0)
		for _, disruption := range node.Disruptions {
			rows = append(rows,
				[]string{disruption.ID, strconv.FormatBool(disruption.ShouldEvict), disruption.StartTime.String(), disruption.Status.State, disruption.Status.Message},
			)
		}
		p.PrintTable(tableHeaders, rows)
	}

}

func (p *Printer) PrintDisplayNodeList(nodes []types.DisplayNode) {
	tableHeaders := []string{
		"Node Name",
		"Update Domain",
		"Total Cores",
		"Total Memory",
		"System Reserved Cores",
		"System Reserved Memory",
		"Remaining Cores",
		"Remaining Memory",
		"# Disruptions",
		"Cluster ID",
		"State",
		"Message",
	}
	rows := make([][]string, 0)
	for _, node := range nodes {
		rows = append(rows,
			[]string{
				node.GetName().Value(),
				node.GetUpdateDomain().Value(),
				node.TotalResources.GetCores().Value(),
				node.TotalResources.GetMemory().Value(),
				node.SystemReserved.GetCores().Value(),
				node.SystemReserved.GetMemory().Value(),
				node.RemainingResources.GetCores().Value(),
				node.RemainingResources.GetMemory().Value(),
				strconv.Itoa(len(node.Disruptions)),
				node.GetClusterID().Value(),
				node.Status.GetState().Value(),
				node.Status.GetMessage().Value(),
			},
		)
	}
	p.PrintTable(tableHeaders, rows)
}
