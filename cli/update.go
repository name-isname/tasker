package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	updateTitle     string
	updateDesc      string
	updatePriority  string
	updateRanking   float64
	updateParent    string
)

var updateCmd = &cobra.Command{
	Use:   "update <pid>",
	Short: "Update process attributes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		// Check if at least one flag is provided
		cmd.Flags()
		if updateTitle == "" && updateDesc == "" && updatePriority == "" && updateRanking == 0 && updateParent == "" {
			return fmt.Errorf("at least one attribute must be specified")
		}

		var titlePtr *string
		var descPtr *string
		var priorityPtr *core.ProcessPriority
		var parentIDPtr *uint

		if cmd.Flags().Changed("title") {
			titlePtr = &updateTitle
		}
		if cmd.Flags().Changed("desc") {
			descPtr = &updateDesc
		}
		if updatePriority != "" {
			p := core.ProcessPriority(updatePriority)
			priorityPtr = &p
		}
		if cmd.Flags().Changed("parent") {
			parentID, err := strconv.ParseUint(updateParent, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid parent ID: %w", err)
			}
			pid := uint(parentID)
			parentIDPtr = &pid
		}

		// Handle ranking separately
		if cmd.Flags().Changed("ranking") {
			if err := core.SetProcessRanking(uint(id), updateRanking); err != nil {
				return err
			}
		}

		if err := core.UpdateProcess(uint(id), titlePtr, descPtr, priorityPtr, parentIDPtr); err != nil {
			return err
		}

		fmt.Printf("Process %d updated\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.GroupID = GroupProcess
	updateCmd.Flags().StringVarP(&updateTitle, "title", "t", "", "New title")
	updateCmd.Flags().StringVarP(&updateDesc, "desc", "D", "", "New description")
	updateCmd.Flags().StringVar(&updatePriority, "priority", "", "New priority (low, medium, high)")
	updateCmd.Flags().Float64Var(&updateRanking, "ranking", 0, "New ranking weight")
	updateCmd.Flags().StringVar(&updateParent, "parent", "", "Parent process ID")
}
