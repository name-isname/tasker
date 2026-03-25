package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	blockReason string
)

var blockCmd = &cobra.Command{
	Use:   "block <pid>",
	Short: "Block a process (set status to blocked)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		if err := core.ChangeProcessState(uint(id), core.StatusBlocked, blockReason); err != nil {
			return err
		}

		if blockReason != "" {
			fmt.Printf("Process %d blocked: %s\n", id, blockReason)
		} else {
			fmt.Printf("Process %d blocked\n", id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(blockCmd)
	blockCmd.GroupID = GroupState
	blockCmd.Flags().StringVarP(&blockReason, "message", "m", "", "Reason for blocking")
}
