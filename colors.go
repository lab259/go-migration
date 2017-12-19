package migration

import "github.com/fatih/color"

type cprintf func(format string, args ...interface{}) string

var (
	styleBold = color.New(color.Bold).SprintfFunc()

	styleNormal         = color.New(color.FgHiBlack).SprintfFunc()
	styleMigrationTitle = color.New(color.FgWhite).SprintfFunc()
	styleMigrationId    = color.New(color.FgCyan).SprintfFunc()

	styleDuration         = color.New(color.FgWhite).SprintfFunc()
	styleDurationSlow     = color.New(color.FgYellow).SprintfFunc()
	styleDurationVerySlow = color.New(color.FgRed).SprintfFunc()

	styleSuccess = color.New(color.Bold, color.FgGreen).SprintfFunc()
	styleWarning = color.New(color.FgYellow).SprintfFunc()
	styleError   = color.New(color.FgRed).SprintfFunc()
)
