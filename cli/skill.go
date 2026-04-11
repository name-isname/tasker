package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	skillOutput string
	skillAuto   bool
	skillGlobal bool

	// Uninstall flags
	uninstallGlobal bool
	uninstallAll    bool
)

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Generate AI skill documentation for taskctl",
	Long: `Generate a comprehensive SKILL.md file that helps AI agents understand
how to use taskctl effectively. The output includes command reference, data model
documentation, usage examples, and best practices for AI collaboration.

This enables Claude Code and other AI agents to effectively help users with
task management, progress tracking, and project organization.

Recommended:
  taskctl skill -a                 # Recommended: Install to local project

Examples:
  taskctl skill                    # Output to stdout
  taskctl skill -a                 # Install to local project (.claude/skills/taskctl/)
  taskctl skill -g                 # Install globally (~/.claude/skills/taskctl/)
  taskctl skill-uninstall          # Uninstall local skill
  taskctl skill-uninstall -g       # Uninstall global skill
  taskctl skill-uninstall --all    # Uninstall both locations`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		content := generateSkillMarkdown()

		// Handle global flag (takes precedence)
		if skillGlobal {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			outputPath := filepath.Join(homeDir, ".claude/skills/taskctl/SKILL.md")
			return writeSkillFile(content, outputPath, true)
		}

		// Handle auto flag
		if skillAuto {
			outputPath := ".claude/skills/taskctl/SKILL.md"
			return writeSkillFile(content, outputPath, false)
		}

		return writeSkillFile(content, skillOutput, false)
	},
}

var skillUninstallCmd = &cobra.Command{
	Use:   "skill-uninstall",
	Short: "Uninstall the taskctl skill",
	Long: `Remove the taskctl skill file from local or global location.

Local uninstall removes .claude/skills/taskctl/
Global uninstall removes ~/.claude/skills/taskctl/

Examples:
  taskctl skill-uninstall          # Uninstall local skill
  taskctl skill-uninstall -g       # Uninstall global skill
  taskctl skill-uninstall --all    # Uninstall both locations`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return uninstallSkill(uninstallGlobal, uninstallAll)
	},
}

func generateSkillMarkdown() string {
	var sb strings.Builder

	// YAML Frontmatter
	sb.WriteString(generateFrontmatter())
	sb.WriteString("\n")

	// Overview
	sb.WriteString("## Overview\n\n")
	sb.WriteString(rootCmd.Long)
	sb.WriteString("\n\n")

	// Data Model
	sb.WriteString(generateDataModel())

	// Commands Reference
	sb.WriteString("## Commands Reference\n\n")
	sb.WriteString(generateCommandsByGroup())

	// Usage Examples
	sb.WriteString(generateUsageExamples())

	// Best Practices
	sb.WriteString(generateBestPractices())

	// AI Features
	sb.WriteString(generateAIFeatures())

	return sb.String()
}

func generateFrontmatter() string {
	return "---\n" +
		"name: taskctl\n" +
		"description: >-\n" +
		"  Process-oriented task management CLI with TUI and Web UI. Use this skill when helping users manage\n" +
		"  tasks, track progress, organize work into hierarchical structures, search task history, or analyze\n" +
		"  productivity data. Trigger for: task management, todo lists, project tracking, progress logging,\n" +
		"  task search, \"show my tasks\", \"create task\", \"task status\", \"what did I work on\", project organization,\n" +
		"  workflow management, process management, or any task/progress/todo related queries. This skill provides\n" +
		"  specialized knowledge for taskctl CLI including all commands (spawn, ps, inspect, block, wake, terminate,\n" +
		"  grep, timeline, stats, tree, export), data model (Process/Log entities with state transitions), and best practices.\n" +
		"---\n\n"
}

func generateDataModel() string {
	return "## Data Model\n\n" +
		"taskctl models tasks as OS \"Processes\" with lifecycle states, not simple todo items. " +
		"This enables rich state tracking and chronological logging.\n\n" +
		"### Process Entity\n\n" +
		"| Field | Type | Description |\n" +
		"|-------|------|-------------|\n" +
		"| id | uint | Primary key (auto-increment) |\n" +
		"| parent_id | uint | Optional parent for hierarchical sub-processes |\n" +
		"| title | string | Short process name |\n" +
		"| description | string | Detailed context (supports Markdown) |\n" +
		"| status | enum | running, blocked, suspended, terminated |\n" +
		"| priority | enum | low, medium, high |\n" +
		"| ranking | float64 | Custom sort weight (like Linux nice value) |\n" +
		"| created_at | timestamp | Creation time |\n" +
		"| updated_at | timestamp | Last modification |\n\n" +
		"**Status Transitions:**\n" +
		"- running ↔ blocked (external dependencies)\n" +
		"- running ↔ suspended (intentional pausing)\n" +
		"- any state → terminated (completed or cancelled)\n\n" +
		"### Log Entity\n\n" +
		"Chronological timeline entries for tracking progress and state changes:\n\n" +
		"| Field | Type | Description |\n" +
		"|-------|------|-------------|\n" +
		"| id | uint | Primary key |\n" +
		"| process_id | uint | Foreign key to Process |\n" +
		"| log_type | enum | state_change (auto) or progress (manual) |\n" +
		"| content | string | Markdown text |\n" +
		"| created_at | timestamp | Entry time |\n\n" +
		"**Important:** State changes via `block`/`wake`/`terminate` commands automatically create state_change logs. " +
		"Users only manually create progress logs.\n\n" +
		"### Key Concepts\n\n" +
		"1. **Process-oriented**: Tasks have lifecycle states and transitions\n" +
		"2. **Hierarchical**: Parent-child relationships enable project breakdown\n" +
		"3. **Chronological**: All activity logged with timestamps\n" +
		"4. **Full-text search**: FTS5 enables instant search across all content\n\n"
}

func generateCommandsByGroup() string {
	groups := []struct {
		id    string
		title string
	}{
		{GroupProcess, "Process Management"},
		{GroupState, "State Management"},
		{GroupLogs, "Log Management"},
		{GroupAnalysis, "Analysis & Search"},
		{GroupUI, "Interface & Export"},
	}

	var sb strings.Builder
	for _, g := range groups {
		sb.WriteString(fmt.Sprintf("### %s\n\n", g.title))
		sb.WriteString(formatCommandTable(g.id))
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatCommandTable(groupID string) string {
	var rows []string
	rows = append(rows, "| Command | Description | Key Flags | Example |")
	rows = append(rows, "|---------|-------------|-----------|---------|")

	for _, cmd := range rootCmd.Commands() {
		// Skip help command, the skill command itself, and completion command
		// Check if command is runnable (has either Run or RunE)
		if cmd.GroupID == groupID && cmd.Use != "help" && cmd.Use != "skill" && cmd.Use != "completion" && cmd.Use != "version" {
			flags := getFlagSummary(cmd)
			example := getExampleSummary(cmd)
			row := fmt.Sprintf("| **%s** | %s | %s | %s |",
				cmd.Use, cmd.Short, flags, example)
			rows = append(rows, row)
		}
	}

	return strings.Join(rows, "\n") + "\n"
}

func getFlagSummary(cmd *cobra.Command) string {
	var flags []string
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		// Skip deprecated or hidden flags
		if !f.Hidden && f.Name != "db" && f.Name != "local" && f.Name != "xml" {
			shorthand := ""
			if f.Shorthand != "" {
				shorthand = "-" + f.Shorthand + ", "
			}
			flags = append(flags, shorthand+"--"+f.Name)
		}
	})
	if len(flags) == 0 {
		return "-"
	}
	if len(flags) > 3 {
		return strings.Join(flags[:2], ", ") + "..."
	}
	return strings.Join(flags, ", ")
}

func getExampleSummary(cmd *cobra.Command) string {
	example := strings.TrimSpace(cmd.Example)
	if example == "" {
		return "-"
	}
	// Get first line only
	lines := strings.Split(example, "\n")
	firstLine := strings.TrimSpace(lines[0])
	if len(firstLine) > 40 {
		return firstLine[:37] + "..."
	}
	return firstLine
}

func generateUsageExamples() string {
	// Note: We use backticks for markdown code blocks in the output
	// In Go raw string literals we can't use backticks, so we construct the string
	return "## Usage Examples\n\n" +
		"### Example 1: Basic Task Management\n\n" +
		"```bash\n" +
		"# Create a new process with description\n" +
		"taskctl spawn \"Build web app\" -D \"Create personal website with React\"\n\n" +
		"# List all running processes\n" +
		"taskctl ps\n\n" +
		"# Add a progress log\n" +
		"taskctl log 1 \"Started with React setup and Vite\"\n\n" +
		"# Block the process with a reason\n" +
		"taskctl block 1 -m \"Waiting for API key from external service\"\n" +
		"```\n\n" +
		"**Output:**\n" +
		"```\n" +
		"Spawned process 1: Build web app\n" +
		"ID  Status    Title           Priority\n" +
		"1   running   Build web app   medium\n" +
		"Added log to process 1\n" +
		"Process 1 blocked: Waiting for API key from external service\n" +
		"```\n\n" +
		"### Example 2: Hierarchical Project Structure\n\n" +
		"```bash\n" +
		"# Create parent process\n" +
		"PARENT=$(taskctl spawn \"Launch SaaS product\" | grep -o '[0-9]*' | head -1)\n\n" +
		"# Add sub-processes for different workstreams\n" +
		"taskctl spawn \"Design database schema\" -p $PARENT -P high\n" +
		"taskctl spawn \"Implement REST API\" -p $PARENT\n" +
		"taskctl spawn \"Build React frontend\" -p $PARENT\n\n" +
		"# View the entire process tree\n" +
		"taskctl tree\n" +
		"```\n\n" +
		"**Output:**\n" +
		"```\n" +
		"Launch SaaS product\n" +
		"├── Design database schema [high]\n" +
		"├── Implement REST API\n" +
		"└── Build React frontend\n" +
		"```\n\n" +
		"### Example 3: AI Integration with XML Output\n\n" +
		"```bash\n" +
		"# Get processes as XML for programmatic parsing\n" +
		"taskctl ps --xml\n\n" +
		"# Get single process with full details\n" +
		"taskctl inspect 1 --xml\n\n" +
		"# Search across all processes and logs\n" +
		"taskctl grep \"database\"\n\n" +
		"# View activity timeline\n" +
		"taskctl timeline\n\n" +
		"# Get statistics for the past 7 days\n" +
		"taskctl stats 7\n" +
		"```\n\n" +
			"```xml\n" +
			"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
			"<process_list>\n" +
			"<process id=\"1\" parent_id=\"null\">\n" +
			"  <name><![CDATA[Build web app]]></name>\n" +
			"  <state>running</state>\n" +
			"  <importance>medium</importance>\n" +
			"</process>\n" +
			"</process_list>\n" +
			"```\n\n" +
			"### Example 4: Using the TUI and Web UI\n\n" +
		"```bash\n" +
		"# Launch the terminal UI (recommended for interactive use)\n" +
		"taskctl tui\n\n" +
		"# Start the web server\n" +
		"taskctl web\n" +
		"# Visit http://localhost:8080\n\n" +
		"# Export a process as Markdown\n" +
		"taskctl export 1\n" +
		"```\n\n"
}

func generateBestPractices() string {
	return "## Best Practices for AI Agents\n\n" +
		"### When to Use Sub-Processes\n\n" +
		"Use sub-processes when:\n" +
		"- Tasks can be worked independently\n" +
		"- Different team members own different parts\n" +
		"- You need to track progress separately\n" +
		"- Each sub-task has a clear deliverable\n\n" +
		"**Keep hierarchy depth to 3-4 levels maximum** for maintainability.\n\n" +
		"### State Management Guidelines\n\n" +
		"| State | When to Use | Example |\n" +
		"|-------|-------------|---------|\n" +
		"| **running** | Default active state | Normal work in progress |\n" +
		"| **blocked** | External dependency | Waiting for code review, API access, approval |\n" +
		"| **suspended** | Intentional pause | Deprioritized, time-boxed, on hold |\n" +
		"| **terminated** | Completed or cancelled | Done, abandoned, merged upstream |\n\n" +
		"**Important:** Use `block`/`wake`/`terminate` commands instead of manually setting status. " +
		"These create automatic state_change logs for audit trail.\n\n" +
		"### Writing Effective Progress Logs\n\n" +
		"**Good logs:**\n" +
		"- \"Decided to use PostgreSQL instead of MySQL for XML support\"\n" +
		"- \"API returns 500 on /users endpoint when email contains '+'\"\n" +
		"- \"Blocked by https://github.com/user/repo/issues/123\"\n\n" +
		"**Poor logs:**\n" +
		"- \"Worked on API\" (too vague)\n" +
		"- \"Fixed bug\" (no details)\n" +
		"- \"Made progress\" (no context)\n\n" +
		"### Priority vs Ranking\n\n" +
		"- **Priority** (low/medium/high): Urgency and importance\n" +
		"- **Ranking** (numeric): Manual sort order for UI display\n\n" +
		"Use ranking when priority alone doesn't capture the desired order. " +
		"Example: Two high-priority items where one must be done first.\n\n" +
		"### Search and Navigation\n\n" +
		"- Use `grep` for full-text search across all processes and logs\n" +
		"- Use `tree` to visualize hierarchy\n" +
		"- Use `timeline` to see recent activity across all processes\n" +
		"- Use `stats` to identify productivity patterns\n\n"
}

func generateAIFeatures() string {
	return "## AI-Specific Features\n\n" +
		"### XML Output Mode\n\n" +
		"All commands support the `--xml` flag for LLM-optimized structured output:\n\n" +
		"```bash\n" +
		"taskctl ps --xml        # All processes with semantic tags\n" +
		"taskctl inspect 1 --xml # Single process with activity log\n" +
		"taskctl logs 1 --xml    # Activity log entries only\n" +
		"```\n\n" +
		"**Why XML for AI?** The XML format uses semantic element names that map to natural language:\n" +
		"- `<process>` instead of generic `item`\n" +
		"- `<name>` instead of `title` for clarity\n" +
		"- `<state>` instead of `status` for process semantics\n" +
		"- `<importance>` instead of `priority` for better context\n" +
		"- `<activity_log>` instead of `logs` for temporal clarity\n" +
		"- CDATA sections protect user content from XML injection\n\n" +
		"**For AI Agents:**\n" +
		"- **Always use `--xml`** when parsing task data programmatically\n" +
		"- XML's semantic tags improve comprehension and reduce parsing errors\n" +
		"- Use `taskctl ps --xml` to get complete task context for analysis\n" +
		"- Use `taskctl inspect 1 --xml` when you need logs with process details\n" +
		"- The format is stable and designed for LLM consumption\n\n" +
		"**XML Structure:**\n" +
		"- Process: `<process id=\"N\" parent_id=\"null\"><name>...</name><state>running</state>...</process>`\n" +
		"- Log Entry: `<entry id=\"N\" type=\"note|state_change\" timestamp=\"...\"><content>...</content></entry>`\n\n" +
		"### Database Management\n\n" +
		"- `--db <path>`: Custom database location\n" +
		"- `--local` / `-L`: Use `./taskctl.db` instead of `~/.taskctl/taskctl.db`\n" +
		"- Default global location (`~/.taskctl/taskctl.db`) ensures consistent access across directories\n\n" +
		"**For AI:** When helping users, prefer `--local` for project-specific databases " +
		"and default global for personal task management.\n\n" +
		"### Search and Analysis Commands\n\n" +
		"| Command | Use Case | AI Scenario |\n" +
		"|---------|----------|-------------|\n" +
		"| `grep <keyword>` | Full-text search | \"What did I do about the database?\" |\n" +
		"| `timeline` | Global activity stream | \"What did I work on this week?\" |\n" +
		"| `stats [days]` | Activity counts per day | \"Show my productivity pattern\" |\n" +
		"| `tree` | Hierarchical view | \"What's the project structure?\" |\n\n" +
		"### Transactional State Changes\n\n" +
		"**Always prefer** `block`/`wake`/`terminate` over direct status updates:\n" +
		"- Automatically creates state_change log\n" +
		"- Ensures data consistency (transactional)\n" +
		"- Provides complete audit trail\n" +
		"- Includes optional reason/message\n\n" +
		"### Common AI Workflows\n\n" +
		"1. **\"Show me my running tasks\"**\n" +
		"   ```bash\n" +
		"   taskctl ps --filter status=running\n" +
		"   ```\n\n" +
		"2. **\"Create a new task for X\"**\n" +
		"   ```bash\n" +
		"   taskctl spawn \"X\" -D \"Detailed description\"\n" +
		"   ```\n\n" +
		"3. **\"What did I work on today?\"**\n" +
		"   ```bash\n" +
		"   taskctl timeline | grep \"$(date +%Y-%m-%d)\"\n" +
		"   ```\n\n" +
		"4. **\"Mark task 1 as completed\"**\n" +
		"   ```bash\n" +
		"   taskctl terminate 1 -m \"Completed all requirements\"\n" +
		"   ```\n\n" +
		"5. **\"Search for database-related tasks\"**\n" +
		"   ```bash\n" +
		"   taskctl grep database\n" +
		"   ```\n\n" +
			"6. **\"Analyze all my tasks and suggest priorities\"**\n" +
			"   ```bash\n" +
			"   taskctl ps --xml  # Get structured data for AI analysis\n" +
			"   ```\n" +
			"   *Use this when AI needs to analyze task data, find patterns, or provide insights*\n" +
			"\n" +
		"### Integration with Claude Code\n\n" +
		"This skill enables Claude Code to:\n" +
		"- Create and manage processes on behalf of users\n" +
		"- Search and analyze task history\n" +
		"- Suggest appropriate state transitions\n" +
		"- Generate reports from timeline/stats data\n" +
		"- Help organize work into hierarchical structures\n\n" +
		"**To install this skill:**\n" +
		"```bash\n" +
		"taskctl skill -a           # Local install\n" +
		"taskctl skill -g           # Global install\n" +
		"```\n\n"
}

func writeSkillFile(content, outputPath string, isGlobal bool) error {
	if outputPath == "" || outputPath == "-" {
		fmt.Println(content)
		return nil
	}

	// Create directory if needed
	dir := filepath.Dir(outputPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write skill file: %w", err)
	}

	location := "locally"
	if isGlobal {
		location = "globally"
	}
	fmt.Printf("Installed skill %s to: %s\n", location, outputPath)
	return nil
}

func uninstallSkill(global, all bool) error {
	pathsToRemove := []string{}

	if all {
		// Add both local and global paths
		localPath := ".claude/skills/taskctl"
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		globalPath := filepath.Join(homeDir, ".claude/skills/taskctl")
		pathsToRemove = append(pathsToRemove, localPath, globalPath)
	} else if global {
		// Only global
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		globalPath := filepath.Join(homeDir, ".claude/skills/taskctl")
		pathsToRemove = append(pathsToRemove, globalPath)
	} else {
		// Only local
		localPath := ".claude/skills/taskctl"
		pathsToRemove = append(pathsToRemove, localPath)
	}

	// Remove each path
	removedCount := 0
	for _, path := range pathsToRemove {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("No skill found at: %s\n", path)
			continue
		}

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}

		// Determine location based on path
		location := "local"
		if strings.HasPrefix(path, "/") || strings.Contains(path, filepath.Join("", ".claude")) {
			// Absolute path or contains .claude (global path has home dir prefix)
			homeDir, _ := os.UserHomeDir()
			if strings.HasPrefix(path, homeDir) || filepath.IsAbs(path) {
				location = "global"
			}
		}
		fmt.Printf("Uninstalled %s skill from: %s\n", location, path)
		removedCount++
	}

	if removedCount == 0 {
		return fmt.Errorf("no skill files found to uninstall")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(skillCmd)
	skillCmd.GroupID = GroupUI
	skillCmd.Flags().StringVarP(&skillOutput, "output", "o", "", "Output file path (default: stdout)")
	skillCmd.Flags().BoolVarP(&skillAuto, "auto", "a", false, "Install to local project (.claude/skills/taskctl/)")
	skillCmd.Flags().BoolVarP(&skillGlobal, "global", "g", false, "Install globally (~/.claude/skills/taskctl/)")

	// Add uninstall subcommand
	rootCmd.AddCommand(skillUninstallCmd)
	skillUninstallCmd.GroupID = GroupUI
	skillUninstallCmd.Flags().BoolVarP(&uninstallGlobal, "global", "g", false, "Uninstall global skill")
	skillUninstallCmd.Flags().BoolVarP(&uninstallAll, "all", "A", false, "Uninstall both local and global")
}
