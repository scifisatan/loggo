package main

import (
	"regexp"
	"strings"
)

var (
	pidLine        = regexp.MustCompile(`^\w+\s+(\w+)\s+\w+\s+\w+\s+\w+\s+\w+\s+\w+\s+\w\s([\w|\.|\/]+)$`)
	pidStart       = regexp.MustCompile(`^.*: Start proc ([a-zA-Z0-9._:]+) for ([a-z]+ [^:]+): pid=(\d+) uid=(\d+) gids=(.*)$`)
	pidStart51     = regexp.MustCompile(`^.*: Start proc (\d+):([a-zA-Z0-9._:]+)/[a-z0-9]+ for (.*)$`)
	pidStartDalvik = regexp.MustCompile(`^E/dalvikvm\(\s*(\d+)\): >>>>> ([a-zA-Z0-9._:]+) \[ userId:0 \| appId:(\d+) \]$`)
	pidKill        = regexp.MustCompile(`^Killing (\d+):([a-zA-Z0-9._:]+)/[^:]+: (.*)$`)
	pidLeave       = regexp.MustCompile(`^No longer want ([a-zA-Z0-9._:]+) \(pid (\d+)\): .*$`)
	pidDeath       = regexp.MustCompile(`^Process ([a-zA-Z0-9._:]+) \(pid (\d+)\) has died.?$`)
	logLine        = regexp.MustCompile(`^([A-Z])/(.+?)\( *(\d+)\): (.*?)$`)
	bugLine        = regexp.MustCompile(`.*nativeGetEnabledTags.*`)
	backtraceLine  = regexp.MustCompile(`^#(.*?)pc\s(.*?)$`)
)

type StartProc struct {
	packageName string
	target      string
	pid         string
	uid         string
	gids        string
}

func parseStartProc(line string) *StartProc {
	matches := pidStart51.FindStringSubmatch(line)
	if matches != nil {
		return &StartProc{
			packageName: matches[2],
			target:      matches[3],
			pid:         matches[1],
			uid:         "",
			gids:        "",
		}
	}
	matches = pidStart.FindStringSubmatch(line)
	if matches != nil {
		return &StartProc{
			packageName: matches[1],
			target:      matches[2],
			pid:         matches[3],
			uid:         matches[4],
			gids:        matches[5],
		}
	}
	matches = pidStartDalvik.FindStringSubmatch(line)
	if matches != nil {
		return &StartProc{
			packageName: matches[2],
			target:      "",
			pid:         matches[1],
			uid:         matches[3],
			gids:        "",
		}
	}
	return nil
}

func parseDeath(tag, message string, catchallPackages, namedProcesses []string, pids map[string]bool) (string, string) {
	if tag != "ActivityManager" {
		return "", ""
	}
	matches := pidKill.FindStringSubmatch(message)
	if matches != nil {
		pid := matches[1]
		packageLine := matches[2]
		if matchPackages(packageLine, catchallPackages, namedProcesses) && pids[pid] {
			return pid, packageLine
		}
	}
	matches = pidLeave.FindStringSubmatch(message)
	if matches != nil {
		pid := matches[2]
		packageLine := matches[1]
		if matchPackages(packageLine, catchallPackages, namedProcesses) && pids[pid] {
			return pid, packageLine
		}
	}
	matches = pidDeath.FindStringSubmatch(message)
	if matches != nil {
		pid := matches[2]
		packageLine := matches[1]
		if matchPackages(packageLine, catchallPackages, namedProcesses) && pids[pid] {
			return pid, packageLine
		}
	}
	return "", ""
}

func matchPackages(token string, catchallPackages, namedProcesses []string) bool {
	for _, proc := range namedProcesses {
		if token == proc {
			return true
		}
	}
	index := strings.Index(token, ":")
	if index == -1 {
		for _, pkg := range catchallPackages {
			if token == pkg {
				return true
			}
		}
	} else {
		packagePart := token[:index]
		for _, pkg := range catchallPackages {
			if packagePart == pkg {
				return true
			}
		}
	}
	return false
}
