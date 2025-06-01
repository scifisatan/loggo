package main

import (
	"os"
)

const (
	Reset   = "\033[0m"
	Black   = "\033[30m"
	Red     = "\033[91m" // Bright Red
	Green   = "\033[92m" // Bright Green
	Yellow  = "\033[93m" // Bright Yellow
	Blue    = "\033[94m" // Bright Blue
	Magenta = "\033[95m" // Bright Magenta
	Cyan    = "\033[96m" // Bright Cyan
	White   = "\033[97m" // Bright White

	BgBlack   = "\033[100m" // Bright Black background
	BgRed     = "\033[101m" // Bright Red background
	BgGreen   = "\033[102m" // Bright Green background
	BgYellow  = "\033[103m" // Bright Yellow background
	BgBlue    = "\033[104m" // Bright Blue background
	BgMagenta = "\033[105m" // Bright Magenta background
	BgCyan    = "\033[106m" // Bright Cyan background
	BgWhite   = "\033[107m" // Bright White background
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
	return fg + bg + text + Reset // Resets color
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
