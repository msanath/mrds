package main

import (
	"github.com/msanath/mrds/ctl/deploymentplan"
	"github.com/msanath/mrds/ctl/metainstance"
	"github.com/msanath/mrds/ctl/node"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use: "mrds-ctl",
	}

	cmd.AddCommand(node.NewNodeCmd())
	cmd.AddCommand(deploymentplan.NewDeploymentPlanCmd())
	cmd.AddCommand(metainstance.NewInstanceCmd())

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
