package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const version = "2.1.0"

// Log levels
// const logLevels = "VDIWEF"

var logLevelsMap = map[string]int{
	"V": 0, "D": 1, "I": 2, "W": 3, "E": 4, "F": 5,
	"v": 0, "d": 1, "i": 2, "w": 3, "e": 4, "f": 5,
}

func main() {
	config := parseArgs()

	if config.version {
		fmt.Printf("loggo %s\n", version)
		return
	}

	minLevel := logLevelsMap[strings.ToUpper(config.minLevel)]

	baseAdbCommand := []string{"adb"}
	if config.deviceSerial != "" {
		baseAdbCommand = append(baseAdbCommand, "-s", config.deviceSerial)
	}
	if config.useDevice {
		baseAdbCommand = append(baseAdbCommand, "-d")
	}
	if config.useEmulator {
		baseAdbCommand = append(baseAdbCommand, "-e")
	}

	packages := config.packages

	if config.currentApp {
		currentPkg := getCurrentPackage(baseAdbCommand)
		if currentPkg != "" {
			packages = append(packages, currentPkg)
		}
	}

	catchallPackages := []string{}
	namedProcesses := []string{}

	for _, pkg := range packages {
		if strings.Contains(pkg, ":") {
			if strings.HasSuffix(pkg, ":") {
				namedProcesses = append(namedProcesses, pkg[:len(pkg)-1])
			} else {
				namedProcesses = append(namedProcesses, pkg)
			}
		} else {
			catchallPackages = append(catchallPackages, pkg)
		}
	}

	headerSize := config.tagWidth + 1 + 3 + 1
	width := getTerminalWidth()

	// Clear logcat if requested
	clearCmd := append(baseAdbCommand, "logcat", "-c")
	exec.Command(clearCmd[0], clearCmd[1:]...).Run()

	// Get initial PIDs
	pids := getInitialPids(baseAdbCommand, catchallPackages)

	// Start logcat
	adbCommand := append(baseAdbCommand, "logcat", "-v", "brief")
	cmd := exec.Command(adbCommand[0], adbCommand[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting logcat: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(stdout)
	appPid := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Skip bug lines
		if bugLine.MatchString(line) {
			continue
		}

		// Parse log line
		matches := logLine.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		level := matches[1]
		tag := strings.TrimSpace(matches[2])
		owner := matches[3]
		message := matches[4]

		// Handle process start
		start := parseStartProc(line)
		if start != nil {
			if matchPackages(start.packageName, catchallPackages, namedProcesses) {
				pids[start.pid] = true
				appPid = start.pid

				linebuf := "\n"
				linebuf += colorize(strings.Repeat(" ", headerSize-1), "", BgWhite)
				linebuf += indentWrap(fmt.Sprintf(" Process %s created for %s\n", start.packageName, start.target), width, headerSize)
				linebuf += colorize(strings.Repeat(" ", headerSize-1), "", BgWhite)
				linebuf += fmt.Sprintf(" PID: %s   UID: %s   GIDs: %s", start.pid, start.uid, start.gids)
				linebuf += "\n"
				fmt.Print(linebuf)
			}
		}

		// Handle process death
		deadPid, deadPname := parseDeath(tag, message, catchallPackages, namedProcesses, pids)
		if deadPid != "" {
			delete(pids, deadPid)
			linebuf := "\n"
			linebuf += colorize(strings.Repeat(" ", headerSize-1), "", BgRed)
			linebuf += fmt.Sprintf(" Process %s (PID: %s) ended", deadPname, deadPid)
			linebuf += "\n"
			fmt.Print(linebuf)
		}

		// Handle backtrace
		if tag == "DEBUG" {
			if backtraceLine.MatchString(strings.TrimLeft(message, " ")) {
				message = strings.TrimLeft(message, " ")
				owner = appPid
			}
		}

		// Filter by PID
		if len(packages) > 0 && !pids[owner] {
			continue
		}

		// Filter by level
		if logLevelsMap[level] < minLevel {
			continue
		}

		// Build output line
		linebuf := ""

		// Add tag with color (bold)
		if config.tagWidth > 0 {
			color := allocateColor(tag)
			tagFormatted := rightAlign(tag, config.tagWidth)
			linebuf += "\033[1m" + colorize(tagFormatted, color, "") + Reset
			linebuf += " "
		}

		// Add level indicator with background color and black text
		if bgColor, ok := levelBgColors[level]; ok {
			linebuf += colorize(" "+level+" ", Black, bgColor)
		} else {
			linebuf += " " + level + " "
		}
		linebuf += " "

		// Add message with color based on log level
		messageColor := messageColors[level]
		coloredMessage := colorize(message, messageColor, "")
		linebuf += indentWrap(coloredMessage, width, headerSize)

		fmt.Println(linebuf)
	}

	cmd.Wait()
}
