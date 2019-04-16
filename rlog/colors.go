package rlog

import (
	"github.com/fatih/color"
)

var (
	styleBold = color.New(color.Bold).SprintFunc()

	styleNormal         = color.New(color.FgHiBlack).SprintFunc()
	styleMigrationTitle = color.New(color.FgWhite).SprintFunc()
	styleMigrationID    = color.New(color.FgCyan).SprintFunc()

	styleDuration         = color.New(color.FgWhite).SprintFunc()
	styleDurationSlow     = color.New(color.FgYellow).SprintFunc()
	styleDurationVerySlow = color.New(color.FgRed).SprintFunc()

	styleSuccess = color.New(color.Bold, color.FgGreen).SprintFunc()
	styleWarning = color.New(color.FgYellow).SprintFunc()
	styleError   = color.New(color.FgRed).SprintFunc()
)
