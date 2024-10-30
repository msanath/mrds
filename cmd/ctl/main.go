package main

import (
	"github.com/msanath/mrds/pkg/ctl/deploymentplan"
	"github.com/msanath/mrds/pkg/ctl/metainstance"
	"github.com/msanath/mrds/pkg/ctl/node"
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
