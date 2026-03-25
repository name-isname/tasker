package cli

import (
	"testing"
	"github.com/spf13/cobra"
)

// TestBasicCommands verifies that all command structures are valid
func TestBasicCommands(t *testing.T) {
	commands := []*cobra.Command{
		spawnCmd,
		psCmd,
		killCmd,
		inspectCmd,
		updateCmd,
		blockCmd,
		wakeCmd,
		terminateCmd,
		logCmd,
		logsCmd,
		grepCmd,
		treeCmd,
		timelineCmd,
		statsCmd,
		exportCmd,
	}

	for _, cmd := range commands {
		if cmd == nil {
			t.Errorf("Command is nil")
			continue
		}

		// Validate command structure
		if cmd.Use == "" {
			t.Errorf("Command %s has empty Use field", cmd.Use)
		}
		if cmd.Short == "" {
			t.Errorf("Command %s has empty Short field", cmd.Use)
		}
	}
}
