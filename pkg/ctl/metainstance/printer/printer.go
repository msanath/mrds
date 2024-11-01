package printer

import (
	"fmt"
	"strings"

	"github.com/msanath/gondolf/pkg/printer"
	"github.com/msanath/mrds/pkg/ctl/metainstance/types"
)

type Printer struct {
	printer.PlainText
}

func NewPrinter() *Printer {
	return &Printer{
		PlainText: printer.NewPlainTextPrinter(),
	}
}

func (p *Printer) PrintDisplayMetaInstance(metaInstance types.DisplayMetaInstance) {
	p.PrintHeader("Meta Instance Info")

	p.PrintDisplayField(metaInstance.GetName())
	p.PrintDisplayField(metaInstance.Metadata.GetID())
	p.PrintDisplayField(metaInstance.Metadata.GetVersion())
	p.PrintDisplayField(metaInstance.GetDeploymentPlanName())
	p.PrintDisplayField(metaInstance.GetDeploymentID())
	p.PrintEmptyLine()

	p.PrintHeader("Status")
	p.PrintDisplayFieldWithIndent(metaInstance.Status.GetState())
	p.PrintDisplayFieldWithIndent(metaInstance.Status.GetMessage())
	p.PrintEmptyLine()

	p.PrintHeader("Runtime Instances")
	if len(metaInstance.RuntimeInstances) == 0 {
		p.PrintWarning("No runtime instances found")
	} else {
		tableHeaders := []string{"Instance ID", "Node ID", "Is Active", "State", "Status Message"}
		rows := make([][]string, 0)
		for _, instance := range metaInstance.RuntimeInstances {
			rows = append(rows, []string{
				instance.GetID().Value(),
				instance.GetNodeName().Value(),
				instance.GetIsActive().Value(),
				instance.Status.GetState().Value(),
				instance.Status.GetMessage().Value(),
			})
		}
		p.PrintTable(tableHeaders, rows)
	}
	p.PrintEmptyLine()

	p.PrintHeader("Operations")
	if len(metaInstance.Operations) == 0 {
		p.PrintWarning("No operations found")
	} else {
		tableHeaders := []string{"Operation ID", "Type", "Intent ID", "State", "Status Message"}
		rows := make([][]string, 0)
		for _, operation := range metaInstance.Operations {
			rows = append(rows, []string{
				operation.GetID().Value(),
				operation.GetType().Value(),
				operation.GetIntentID().Value(),
				operation.Status.GetState().Value(),
				operation.Status.GetMessage().Value(),
			})
		}
		p.PrintTable(tableHeaders, rows)
	}
}

func (p *Printer) PrintDisplayMetaInstanceList(metaInstances []types.DisplayMetaInstance) {
	tableHeaders := []string{
		"Name", "Deployment Plan Name", "Deployment ID", "State", "Status Message",
		"Instance Info", "Operations",
	}

	rows := make([][]string, 0)
	for _, metaInstance := range metaInstances {

		instanceInfo := []string{}

		for _, runtimeInstance := range metaInstance.RuntimeInstances {
			if len(metaInstance.RuntimeInstances) > 0 {
				instanceInfo = append(instanceInfo, fmt.Sprintf("Node: %s, IsActive: %s", runtimeInstance.NodeName, runtimeInstance.GetIsActive().Value()))
			}
		}
		instanceStr := strings.Join(instanceInfo, "\n")

		operationsInfo := []string{}
		for _, operation := range metaInstance.Operations {
			operationsInfo = append(operationsInfo, fmt.Sprintf("ID: %s, Type: %s, State: %s", operation.ID, operation.Type, operation.Status.GetState().Value()))
		}
		operationsStr := strings.Join(operationsInfo, "\n")

		rows = append(rows, []string{
			metaInstance.GetName().Value(),
			metaInstance.GetDeploymentPlanName().Value(),
			metaInstance.GetDeploymentID().Value(),
			metaInstance.Status.GetState().Value(),
			metaInstance.Status.GetMessage().Value(),
			instanceStr,
			operationsStr,
		})
	}
	p.PrintTable(tableHeaders, rows, printer.WithRowSeparator())
}
