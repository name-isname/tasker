package tui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Markdown renderer for terminal UI
// Supports basic markdown syntax: bold, italic, code, headers, lists, quotes, code blocks

// RenderMarkdown renders markdown text to terminal-friendly output
func RenderMarkdown(text string, width int) string {
	if text == "" {
		return ""
	}

	// Split into lines for processing
	lines := strings.Split(text, "\n")
	var result strings.Builder

	inCodeBlock := false
	codeBlockLines := []string{}
	codeBlockLang := ""

	for _, line := range lines {
		// Check for code block
		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				// Start code block
				inCodeBlock = true
				codeBlockLang = strings.TrimPrefix(line, "```")
				codeBlockLines = []string{}
			} else {
				// End code block
				inCodeBlock = false
				if len(codeBlockLines) > 0 {
					result.WriteString(renderCodeBlock(strings.Join(codeBlockLines, "\n"), codeBlockLang, width))
				}
				codeBlockLines = []string{}
				codeBlockLang = ""
			}
			continue
		}

		if inCodeBlock {
			codeBlockLines = append(codeBlockLines, line)
			continue
		}

		// Process regular line
		processedLine := renderInlineMarkdown(line, width)
		result.WriteString(processedLine + "\n")
	}

	return result.String()
}

// renderInlineMarkdown renders inline markdown (bold, italic, code, links)
func renderInlineMarkdown(text string, width int) string {
	// Trim leading whitespace for processing
	trimmed := strings.TrimLeft(text, " ")
	indent := strings.Repeat(" ", len(text)-len(trimmed))
	text = trimmed

	// Check for headers
	if strings.HasPrefix(text, "#") {
		return renderHeader(text, width)
	}

	// Check for quote
	if strings.HasPrefix(text, ">") {
		return renderQuote(text, width)
	}

	// Check for list
	if strings.HasPrefix(text, "-") || strings.HasPrefix(text, "*") {
		return renderListItem(text, width)
	}

	// Process inline elements (order matters for nested patterns)
	result := text

	// Links: [text](url) -> text (url in gray if exists)
	result = renderLinks(result)

	// Code: `text`
	result = renderInlineCode(result)

	// Bold: **text**
	result = renderBold(result)

	// Italic: *text* (but not **text**)
	result = renderItalic(result)

	return indent + result
}

// renderCodeBlock renders a code block with styling
func renderCodeBlock(code, lang string, width int) string {
	// Code block style: gray background
	style := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252")).
		Padding(0, 1).
		Width(width)

	lines := strings.Split(code, "\n")
	for i, line := range lines {
		lines[i] = " " + line
	}

	return style.Render(strings.Join(lines, "\n")) + "\n"
}

// renderInlineCode renders inline code
func renderInlineCode(text string) string {
	re := regexp.MustCompile("`([^`]+)`")
	return re.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(match, "`")
		return mdCodeStyle.Render(content)
	})
}

// renderBold renders bold text
func renderBold(text string) string {
	re := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(strings.Trim(match, "*"), "*")
		return mdBoldStyle.Render(content)
	})
}

// renderItalic renders italic text (but not bold)
func renderItalic(text string) string {
	// First, skip already processed bold sections
	re := regexp.MustCompile(`\*([^*]+)\*`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		// Skip if this was part of bold (already processed)
		if strings.HasPrefix(match, "**") {
			return match
		}
		content := strings.Trim(match, "*")
		return mdItalicStyle.Render(content)
	})
}

// renderLinks renders markdown links
func renderLinks(text string) string {
	re := regexp.MustCompile(`\[([^\]]+)\]\(([^\)]*)\)`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		linkText := parts[1]
		url := parts[2]

		if url != "" {
			// Show link text with URL in parentheses (gray)
			return mdLinkStyle.Render(linkText) + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(" ("+url+")")
		}
		return mdLinkStyle.Render(linkText)
	})
}

// renderHeader renders markdown headers
func renderHeader(text string, width int) string {
	// Count leading #
	level := 0
	for _, ch := range text {
		if ch == '#' {
			level++
		} else {
			break
		}
	}

	if level > 6 || level == 0 {
		return text
	}

	// Extract header text
	content := strings.TrimLeft(text, "#")
	content = strings.TrimLeft(content, " ")

	var style lipgloss.Style
	var prefix string

	switch level {
	case 1:
		style = mdHeader1Style
		prefix = "━━ "
	case 2:
		style = mdHeader2Style
		prefix = "── "
	case 3:
		style = mdHeader3Style
		prefix = "··· "
	default:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		prefix = strings.Repeat("·", level) + " "
	}

	return style.Render(prefix + content)
}

// renderListItem renders list items
func renderListItem(text string, width int) string {
	// Extract content after - or *
	content := strings.TrimLeft(text, "-*")
	content = strings.TrimLeft(content, " ")

	bullet := "•"
	return "  " + bullet + " " + renderInlineMarkdown(content, width-4)
}

// renderQuote renders block quotes
func renderQuote(text string, width int) string {
	// Extract content after >
	content := strings.TrimLeft(text, ">")
	content = strings.TrimLeft(content, " ")

	return "┃ " + mdQuoteStyle.Render(content)
}
