package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var treeCmd = &cobra.Command{
	Use:   "tree [id]",
	Short: "Show process tree (like pstree)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Show full tree
			output, err := core.FormatFullTree()
			if err != nil {
				return err
			}
			fmt.Print(output)
			return nil
		}

		// Show tree for specific process
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid process ID: %w", err)
		}

		node, err := core.GetProcessTree(uint(id))
		if err != nil {
			return err
		}

		fmt.Print(core.FormatProcessTree(node, "", true))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)
}
