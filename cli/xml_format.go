package cli

import (
	"fmt"
	"strings"
	"taskctl/core"
	"time"
)

// LLM-optimized XML format with semantic tags
func formatProcessAsXML(process *core.Process) string {
	var sb strings.Builder

	parentID := "null"
	if process.ParentID != nil {
		parentID = fmt.Sprintf("%d", *process.ParentID)
	}

	sb.WriteString(fmt.Sprintf(`<process id="%d" parent_id="%s">
  <name><![CDATA[%s]]></name>
  <description><![CDATA[%s]]></description>
  <state>%s</state>
  <importance>%s</importance>
  <created>%s</created>
  <modified>%s</modified>
</process>`,
		process.ID,
		parentID,
		escapeXML(process.Title),
		escapeXML(process.Description),
		process.Status,
		process.Priority,
		process.CreatedAt.Format(time.RFC3339),
		process.UpdatedAt.Format(time.RFC3339),
	))

	return sb.String()
}

func formatProcessListAsXML(processes []core.Process) string {
	var sb strings.Builder
	sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	sb.WriteString("<process_list>\n")

	for _, p := range processes {
		sb.WriteString(formatProcessAsXML(&p))
		sb.WriteString("\n")
	}

	sb.WriteString("</process_list>")
	return sb.String()
}

func formatProcessDetailAsXML(process *core.Process) string {
	var sb strings.Builder
	sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	sb.WriteString("<process_detail>\n")
	sb.WriteString(formatProcessAsXML(process))
	sb.WriteString("\n")

	// Include logs
	if len(process.Logs) > 0 {
		sb.WriteString("  <activity_log>\n")
		for _, log := range process.Logs {
			sb.WriteString(formatLogAsXML(&log, "    "))
		}
		sb.WriteString("  </activity_log>\n")
	}

	sb.WriteString("</process_detail>")
	return sb.String()
}

func formatLogsAsXML(logs []core.Log) string {
	var sb strings.Builder
	sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	sb.WriteString("<activity_log>\n")

	for _, log := range logs {
		sb.WriteString(formatLogAsXML(&log, "  "))
	}

	sb.WriteString("</activity_log>")
	return sb.String()
}

func formatLogAsXML(log *core.Log, indent string) string {
	logType := "note"
	if log.LogType == core.LogTypeStateChange {
		logType = "state_change"
	}

	return fmt.Sprintf(`%s<entry id="%d" type="%s" timestamp="%s">
%s  <content><![CDATA[%s]]></content>
%s</entry>
`,
		indent,
		log.ID,
		logType,
		log.CreatedAt.Format(time.RFC3339),
		indent,
		escapeXML(log.Content),
		indent,
	)
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
