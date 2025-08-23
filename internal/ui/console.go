package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"goup/internal/config"
	"goup/internal/dependency"
)

// Modern ANSI color palette
const (
	// Reset and modifiers
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	// Modern color palette
	Black   = "\033[30m"
	Red     = "\033[91m" // Bright red
	Green   = "\033[92m" // Bright green
	Yellow  = "\033[93m" // Bright yellow
	Blue    = "\033[94m" // Bright blue
	Magenta = "\033[95m" // Bright magenta
	Cyan    = "\033[96m" // Bright cyan
	White   = "\033[97m" // Bright white
	Gray    = "\033[90m" // Dark gray

	// Semantic colors
	Primary   = Cyan + Bold
	Success   = Green + Bold
	Error     = Red + Bold
	Warning   = Yellow + Bold
	Info      = Blue + Bold
	Debug     = Magenta + Dim
	Secondary = Gray
	Accent    = White + Bold

	// Beautiful symbols
	SymbolSuccess  = "✨"
	SymbolError    = "💥"
	SymbolWarning  = "⚡"
	SymbolInfo     = "💡"
	SymbolProgress = "🚀"
	SymbolDebug    = "🔍"
	SymbolPackage  = "📦"
	SymbolIndirect = "🔗"
	SymbolArrow    = "▶"
	SymbolCheck    = "✓"
	SymbolCross    = "✗"
	SymbolBullet   = "•"
	SymbolDivider  = "─"

	// Table characters
	TableTopLeft     = "┌"
	TableTopRight    = "┐"
	TableBottomLeft  = "└"
	TableBottomRight = "┘"
	TableHorizontal  = "─"
	TableVertical    = "│"
	TableCross       = "┼"
	TableTeeDown     = "┬"
	TableTeeUp       = "┴"
	TableTeeRight    = "├"
	TableTeeLeft     = "┤"
	TableDotted      = "┈"
)

type console struct {
	noColor bool
	verbose bool
	reader  *bufio.Reader
}

func NewConsole(cfg *config.Config) Console {
	return &console{
		noColor: cfg.NoColor,
		verbose: cfg.Verbose,
		reader:  bufio.NewReader(os.Stdin),
	}
}

func (c *console) Header() {
	fmt.Println()
	c.printBox("🚀 goup - Go Dependency Updater", Primary)
	fmt.Println()
}

func (c *console) Info(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	c.printMessage(SymbolInfo, "INFO", Info, message)
}

func (c *console) Success(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	c.printMessage(SymbolSuccess, "SUCCESS", Success, message)
}

func (c *console) Warning(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	c.printMessage(SymbolWarning, "WARNING", Warning, message)
}

func (c *console) Error(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	c.printMessage(SymbolError, "ERROR", Error, message)
}

func (c *console) Debug(format string, args ...any) {
	if c.verbose {
		message := fmt.Sprintf(format, args...)
		c.printMessage(SymbolDebug, "DEBUG", Debug, message)
	}
}

func (c *console) Progress(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	c.printMessage(SymbolProgress, "PROGRESS", Info, message)
}

func (c *console) ReadInput(prompt string) (string, error) {
	if c.noColor {
		fmt.Printf("\n%s: ", prompt)
	} else {
		fmt.Printf("\n%s%s❯%s %s%s%s: ",
			Primary, Bold, Reset, Accent, prompt, Reset)
	}

	response, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response), nil
}

func (c *console) Confirm(message string) bool {
	if c.noColor {
		fmt.Printf("\n%s (y/N): ", message)
	} else {
		fmt.Printf("\n%s%s?%s %s%s%s %s(y/N):%s ",
			Warning, Bold, Reset,
			Accent, message, Reset,
			Secondary, Reset)
	}

	response, err := c.reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func (c *console) PrintDependencies(deps []dependency.Dependency, title string) {
	if title != "" {
		c.Info("%s", title)
	}
	fmt.Println()

	if len(deps) == 0 {
		return
	}

	// Calculate optimal column widths
	maxPathWidth := c.calculateMaxPathWidth(deps)
	indexWidth := c.calculateIndexWidth(len(deps))
	versionWidth := 15 // Fixed width for versions
	typeWidth := 8     // Fixed width for type

	if c.noColor {
		c.printSimpleTable(deps, indexWidth, maxPathWidth, versionWidth, typeWidth)
		return
	}

	c.printStyledTable(deps, indexWidth, maxPathWidth, versionWidth, typeWidth)
}

func (c *console) calculateMaxPathWidth(deps []dependency.Dependency) int {
	maxWidth := 20 // minimum width
	for _, dep := range deps {
		if len(dep.Path) > maxWidth {
			maxWidth = len(dep.Path)
		}
	}
	// Cap the maximum width to keep table readable
	if maxWidth > 60 {
		maxWidth = 60
	}
	return maxWidth
}

func (c *console) calculateIndexWidth(total int) int {
	totalStr := fmt.Sprintf("%d/%d", total, total)
	return len(totalStr)
}

func (c *console) printSimpleTable(deps []dependency.Dependency, indexWidth, pathWidth, versionWidth, typeWidth int) {
	// Header
	fmt.Printf(" %-*s │ %-*s │ %-*s │ %-*s\n",
		indexWidth, "#",
		pathWidth, "Package",
		versionWidth, "Version",
		typeWidth, "Type")

	// Separator
	fmt.Printf("%s┼%s┼%s┼%s\n",
		strings.Repeat("─", indexWidth+2),
		strings.Repeat("─", pathWidth+2),
		strings.Repeat("─", versionWidth+2),
		strings.Repeat("─", typeWidth+2))

	// Rows
	for i, dep := range deps {
		c.printSimpleDependencyRow(i+1, len(deps), dep, indexWidth, pathWidth, versionWidth, typeWidth)
	}

	// Bottom border
	fmt.Printf("%s┴%s┴%s┴%s\n",
		strings.Repeat("─", indexWidth+2),
		strings.Repeat("─", pathWidth+2),
		strings.Repeat("─", versionWidth+2),
		strings.Repeat("─", typeWidth+2))
	fmt.Println()
}

func (c *console) printStyledTable(deps []dependency.Dependency, indexWidth, pathWidth, versionWidth, typeWidth int) {
	// Top border
	fmt.Printf("   %s%s%s\n", Secondary,
		TableTopLeft+strings.Repeat(TableHorizontal, indexWidth+2)+
			TableTeeDown+strings.Repeat(TableHorizontal, pathWidth+2)+
			TableTeeDown+strings.Repeat(TableHorizontal, versionWidth+2)+
			TableTeeDown+strings.Repeat(TableHorizontal, typeWidth+2)+
			TableTopRight, Reset)

	// Header
	fmt.Printf("   %s%s%s %s%s%-*s%s %s%s%s %s%s%-*s%s %s%s%s %s%s%-*s%s %s%s%s %s%s%-*s%s %s%s%s\n",
		Secondary, TableVertical, Reset,
		Primary, Bold, indexWidth, "#", Reset,
		Secondary, TableVertical, Reset,
		Primary, Bold, pathWidth, "Package", Reset,
		Secondary, TableVertical, Reset,
		Primary, Bold, versionWidth, "Version", Reset,
		Secondary, TableVertical, Reset,
		Primary, Bold, typeWidth, "Type", Reset,
		Secondary, TableVertical, Reset)

	// Header separator
	fmt.Printf("   %s%s%s\n", Secondary,
		TableTeeRight+strings.Repeat(TableHorizontal, indexWidth+2)+
			TableCross+strings.Repeat(TableHorizontal, pathWidth+2)+
			TableCross+strings.Repeat(TableHorizontal, versionWidth+2)+
			TableCross+strings.Repeat(TableHorizontal, typeWidth+2)+
			TableTeeLeft, Reset)

	// Rows
	for i, dep := range deps {
		c.printStyledDependencyRow(i+1, len(deps), dep, indexWidth, pathWidth, versionWidth, typeWidth)

		// Row separator (except for last row)
		if i < len(deps)-1 {
			fmt.Printf("   %s%s%s\n", Secondary,
				TableTeeRight+strings.Repeat(TableDotted, indexWidth+2)+
					TableCross+strings.Repeat(TableDotted, pathWidth+2)+
					TableCross+strings.Repeat(TableDotted, versionWidth+2)+
					TableCross+strings.Repeat(TableDotted, typeWidth+2)+
					TableTeeLeft, Reset)
		}
	}

	// Bottom border
	fmt.Printf("   %s%s%s\n", Secondary,
		TableBottomLeft+strings.Repeat(TableHorizontal, indexWidth+2)+
			TableTeeUp+strings.Repeat(TableHorizontal, pathWidth+2)+
			TableTeeUp+strings.Repeat(TableHorizontal, versionWidth+2)+
			TableTeeUp+strings.Repeat(TableHorizontal, typeWidth+2)+
			TableBottomRight, Reset)
	fmt.Println()
}

func (c *console) printSimpleDependencyRow(index, total int, dep dependency.Dependency, indexWidth, pathWidth, versionWidth, typeWidth int) {
	indexStr := fmt.Sprintf("%d/%d", index, total)
	pathStr := c.truncateString(dep.Path, pathWidth)
	versionStr := c.truncateString(dep.Version, versionWidth)

	typeStr := "direct"
	if dep.Indirect {
		typeStr = "indirect"
	}

	fmt.Printf(" %-*s │ %-*s │ %-*s │ %-*s\n",
		indexWidth, indexStr,
		pathWidth, pathStr,
		versionWidth, versionStr,
		typeWidth, typeStr)
}

func (c *console) printStyledDependencyRow(index, total int, dep dependency.Dependency, indexWidth, pathWidth, versionWidth, typeWidth int) {
	indexStr := fmt.Sprintf("%d/%d", index, total)
	pathStr := c.truncateString(dep.Path, pathWidth)
	versionStr := c.truncateString(dep.Version, versionWidth)

	// Colors and symbols based on dependency type
	var pathColor, typeColor, typeStr string
	if dep.Indirect {
		pathColor = Yellow
		typeColor = Secondary
		typeStr = "indirect"
	} else {
		pathColor = Green
		typeColor = Primary
		typeStr = "direct"
	}

	// Format each column separately to avoid printf complexity
	indexCol := fmt.Sprintf("%-*s", indexWidth, indexStr)
	pathCol := fmt.Sprintf("%-*s", pathWidth, pathStr)
	versionCol := fmt.Sprintf("%-*s", versionWidth, versionStr)
	typeCol := fmt.Sprintf("%-*s", typeWidth, typeStr)

	// Print the formatted row
	fmt.Printf("   %s%s%s %s%s%s %s%s%s %s%s%s%s %s%s%s %s%s%s %s%s%s %s%s%s %s%s%s\n",
		Secondary, TableVertical, Reset,
		Secondary, indexCol, Reset,
		Secondary, TableVertical, Reset,
		pathColor, Bold, pathCol, Reset,
		Secondary, TableVertical, Reset,
		Cyan, versionCol, Reset,
		Secondary, TableVertical, Reset,
		typeColor, typeCol, Reset,
		Secondary, TableVertical, Reset)
}

func (c *console) PrintUpdateResult(updated, total int, hasErrors bool) {
	if c.noColor {
		if hasErrors {
			fmt.Printf("\n[WARNING] Completed with %d/%d dependencies updated\n", updated, total)
		} else {
			fmt.Printf("\n[SUCCESS] All %d dependencies updated successfully!\n", total)
		}
	} else {
		fmt.Println()
		if hasErrors {
			c.printBox(fmt.Sprintf("⚡ Partial Success: %d/%d updated", updated, total), Warning)
		} else {
			c.printBox(fmt.Sprintf("🎉 Complete Success: All %d dependencies updated!", total), Success)
		}
		fmt.Println()
	}
}

// Helper methods
func (c *console) printMessage(symbol, label, color, message string) {
	if c.noColor {
		fmt.Printf("[%s] %s\n", label, message)
	} else {
		fmt.Printf(" %s %s%s[%s]%s %s%s%s\n",
			symbol, color, Bold, label, Reset,
			Accent, message, Reset)
	}
}

func (c *console) printBox(message, color string) {
	if c.noColor {
		fmt.Printf("=== %s ===\n", message)
		return
	}

	width := utf8.RuneCountInString(message) + 4
	top := "╭" + strings.Repeat("─", width) + "╮"
	bottom := "╰" + strings.Repeat("─", width) + "╯"

	fmt.Printf("%s%s%s\n", Primary, top, Reset)
	fmt.Printf("%s│%s  %s%s%s%s %s│%s\n",
		Primary, Reset,
		color, Bold, message, Reset,
		Primary, Reset)
	fmt.Printf("%s%s%s\n", Primary, bottom, Reset)
}

func (c *console) truncateString(s string, maxWidth int) string {
	if utf8.RuneCountInString(s) <= maxWidth {
		return s
	}
	if maxWidth <= 1 {
		return "…"
	}
	// Convert to runes to handle Unicode properly
	runes := []rune(s)
	if len(runes) > maxWidth-1 {
		return string(runes[:maxWidth-1]) + "…"
	}
	return s
}
