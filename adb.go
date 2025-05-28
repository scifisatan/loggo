package main

import (
	"bufio"
	"os/exec"
	"regexp"
	"strings"
)

func getCurrentPackage(baseCommand []string) string {
	cmd := append(baseCommand, "shell", "dumpsys", "activity", "activities")
	output, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`.*TaskRecord.*A[= ]([^ ^}]*)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func getInitialPids(baseCommand []string, catchallPackages []string) map[string]bool {
	pids := make(map[string]bool)

	cmd := append(baseCommand, "shell", "ps")
	output, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return pids
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		matches := pidLine.FindStringSubmatch(line)
		if matches != nil {
			pid := matches[1]
			proc := matches[2]
			for _, pkg := range catchallPackages {
				if proc == pkg {
					pids[pid] = true
					break
				}
			}
		}
	}

	return pids
}
