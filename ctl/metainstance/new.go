package metainstance

import (
	"github.com/spf13/cobra"
)

func NewInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Manage instances of a deployment plan",
	}

	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newMetaInstanceListCmd())
	cmd.AddCommand(newAddRuntimeInstanceCmd())
	cmd.AddCommand(newStopInstanceCmd())
	cmd.AddCommand(newRestartInstanceCmd())
	cmd.AddCommand(newSwapInstanceCmd())
	cmd.AddCommand(newApproveOperationCmd())
	cmd.AddCommand(newCompleteOperationCmd())

	return cmd
}
