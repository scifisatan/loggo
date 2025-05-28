package main

import (
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	packages     []string
	tagWidth     int
	minLevel     string
	colorGC      bool
	alwaysTags   bool
	currentApp   bool
	deviceSerial string
	useDevice    bool
	useEmulator  bool
	clearLogcat  bool
	tags         []string
	ignoredTags  []string
	all          bool
	version      bool
}

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func parseArgs() Config {
	config := Config{}

	flag.StringVar(&config.minLevel, "min-level", "V", "Minimum level to be displayed")
	flag.StringVar(&config.minLevel, "l", "V", "Minimum level to be displayed")
	flag.IntVar(&config.tagWidth, "tag-width", 23, "Width of log tag")
	flag.IntVar(&config.tagWidth, "w", 23, "Width of log tag")
	flag.BoolVar(&config.colorGC, "color-gc", false, "Color garbage collection")
	flag.BoolVar(&config.alwaysTags, "always-display-tags", false, "Always display the tag name")
	flag.BoolVar(&config.currentApp, "current", false, "Filter logcat by current running app")
	flag.StringVar(&config.deviceSerial, "serial", "", "Device serial number")
	flag.StringVar(&config.deviceSerial, "s", "", "Device serial number")
	flag.BoolVar(&config.useDevice, "device", false, "Use first device for log input")
	flag.BoolVar(&config.useDevice, "d", false, "Use first device for log input")
	flag.BoolVar(&config.useEmulator, "emulator", false, "Use first emulator for log input")
	flag.BoolVar(&config.useEmulator, "e", false, "Use first emulator for log input")
	flag.BoolVar(&config.clearLogcat, "clear", false, "Clear the entire log before running")
	flag.BoolVar(&config.clearLogcat, "c", false, "Clear the entire log before running")
	flag.BoolVar(&config.all, "all", false, "Print all log messages")
	flag.BoolVar(&config.all, "a", false, "Print all log messages")
	flag.BoolVar(&config.version, "version", false, "Print version and exit")
	flag.BoolVar(&config.version, "v", false, "Print version and exit")

	// Custom parsing for repeated flags and packages
	var tags, ignoredTags stringSlice
	flag.Var(&tags, "tag", "Filter output by specified tag(s)")
	flag.Var(&tags, "t", "Filter output by specified tag(s)")
	flag.Var(&ignoredTags, "ignore-tag", "Filter output by ignoring specified tag(s)")
	flag.Var(&ignoredTags, "i", "Filter output by ignoring specified tag(s)")

	flag.Usage = func() {
		fmt.Println(colorize("\nUsage:", Blue, ""), colorize("loggo [options] [package ...]", Yellow, ""))
		fmt.Println(colorize("\nOptions:", Yellow, ""))

		// Format: flag(s)   arg   description
		options := []struct {
			flags string
			arg   string
			desc  string
		}{
			{"-min-level, -l", "<level>", "Minimum log level to display (V, D, I, W, E, F)"},
			{"-tag-width, -w", "<n>", "Width of log tag (default 23)"},
			{"-color-gc", "", "Color garbage collection"},
			{"-always-display-tags", "", "Always display the tag name"},
			{"-current", "", "Filter logcat by current running app"},
			{"-serial, -s", "<serial>", "Device serial number"},
			{"-device, -d", "", "Use first device for log input"},
			{"-emulator, -e", "", "Use first emulator for log input"},
			{"-clear, -c", "", "Clear the entire log before running"},
			{"-all, -a", "", "Print all log messages"},
			{"-tag, -t", "<tag>", "Filter output by specified tag(s)"},
			{"-ignore-tag, -i", "<tag>", "Ignore specified tag(s)"},
			{"-version, -v", "", "Print version and exit"},
			{"-h, --help", "", "Show this help message"},
		}

		for _, opt := range options {
			fmt.Printf("  %s %s %s\n",
				colorize(fmt.Sprintf("%-18s", opt.flags), Cyan, ""),
				colorize(fmt.Sprintf("%-10s", opt.arg), Magenta, ""),
				colorize(opt.desc, White, ""),
			)
		}
	}

	flag.Parse()

	config.packages = flag.Args()
	config.tags = []string(tags)
	config.ignoredTags = []string(ignoredTags)

	return config
}
