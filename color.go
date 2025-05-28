package main

import (
	"os"
)

const (
	Reset   = "\033[0m"
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

var knownTags = map[string]string{
	"dalvikvm":        White,
	"Process":         White,
	"ActivityManager": White,
	"ActivityThread":  White,
	"AndroidRuntime":  Cyan,
	"jdwp":            White,
	"StrictMode":      White,
	"DEBUG":           Yellow,
}

var lastUsedColors = []string{Red, Green, Yellow, Blue, Magenta, Cyan}

// Message colors for each log level
var messageColors = map[string]string{
	"V": White,  // Verbose - white
	"D": Blue,   // Debug - blue
	"I": Green,  // Info - green
	"W": Yellow, // Warning - yellow
	"E": Red,    // Error - red
	"F": Red,    // Fatal - red
}

// Background colors for log levels
var levelBgColors = map[string]string{
	"V": BgBlack,  // Verbose - black background
	"D": BgBlue,   // Debug - blue background
	"I": BgGreen,  // Info - green background
	"W": BgYellow, // Warning - yellow background
	"E": BgRed,    // Error - red background
	"F": BgRed,    // Fatal - red background
}

func colorize(text, fg, bg string) string {
	if !isatty() {
		return text
	}
	return fg + bg + text + Reset
}

func isatty() bool {
	if fileInfo, _ := os.Stdout.Stat(); fileInfo != nil {
		return fileInfo.Mode()&os.ModeCharDevice != 0
	}
	return false
}

func allocateColor(tag string) string {
	if color, ok := knownTags[tag]; ok {
		for i, c := range lastUsedColors {
			if c == color {
				lastUsedColors = append(lastUsedColors[:i], lastUsedColors[i+1:]...)
				lastUsedColors = append(lastUsedColors, color)
				break
			}
		}
		return color
	}
	color := lastUsedColors[0]
	lastUsedColors = append(lastUsedColors[1:], color)
	knownTags[tag] = color
	return color
}
