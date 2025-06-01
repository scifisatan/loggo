package main

import (
	"flag"
	"fmt"
)

type Config struct {
	packages     []string
	tagWidth     int
	minLevel     string
	currentApp   bool
	deviceSerial string
	useDevice    bool
	useEmulator  bool
	version      bool
}

type Option struct {
	long        string
	short       string
	arg         string
	desc        string
	defaultVal  any
	configField any
}

func assertPointer[T any](field any, flagName string) (*T, bool) {
	ptr, ok := field.(*T)
	if !ok {
		fmt.Printf("Error: configField for '%s' must be *%T, got %T\n", flagName, *new(T), field)
	}
	return ptr, ok
}

func registerFlag(opt Option) {
	if opt.long == "help" || opt.short == "h" {
		return
	}

	switch def := opt.defaultVal.(type) {
	case string:
		if ptr, ok := assertPointer[string](opt.configField, opt.long); ok {
			flag.StringVar(ptr, opt.long, def, opt.desc)
			if opt.short != "" {
				flag.StringVar(ptr, opt.short, def, opt.desc)
			}
		}
	case int:
		if ptr, ok := assertPointer[int](opt.configField, opt.long); ok {
			flag.IntVar(ptr, opt.long, def, opt.desc)
			if opt.short != "" {
				flag.IntVar(ptr, opt.short, def, opt.desc)
			}
		}
	case bool:
		if ptr, ok := assertPointer[bool](opt.configField, opt.long); ok {
			flag.BoolVar(ptr, opt.long, def, opt.desc)
			if opt.short != "" {
				flag.BoolVar(ptr, opt.short, def, opt.desc)
			}
		}
	default:
		fmt.Printf("Unsupported flag type for '%s'\n", opt.long)
	}
}
func printHelpLine(opt Option) {
	flags := fmt.Sprintf("--%s", opt.long)
	if opt.short != "" {
		flags += fmt.Sprintf(", -%s", opt.short)
	}

	arg := ""
	if opt.arg != "" {
		arg = fmt.Sprintf("<%s>", opt.arg)
	}

	fmt.Printf("  %s %s %s\n",
		colorize(fmt.Sprintf("%-18s", flags), Cyan, ""),
		colorize(fmt.Sprintf("%-10s", arg), Magenta, ""),
		colorize(opt.desc, White, ""),
	)
}

func parseArgs() Config {
	config := Config{}

	var options = []Option{
		{
			long:        "min-level",
			short:       "l",
			arg:         "level",
			desc:        "Minimum log level to display (V, D, I, W, E, F)",
			defaultVal:  "V",
			configField: &config.minLevel,
		},
		{
			long:        "tag-width",
			short:       "w",
			arg:         "width",
			desc:        "Width of log tag (default 23)",
			defaultVal:  23,
			configField: &config.tagWidth,
		},
		{
			long:        "current",
			short:       "",
			arg:         "",
			desc:        "Filter logcat by current running app",
			defaultVal:  false,
			configField: &config.currentApp,
		},
		{
			long:        "serial",
			short:       "s",
			arg:         "serial",
			desc:        "Device serial number",
			defaultVal:  "",
			configField: &config.deviceSerial,
		},
		{
			long:        "device",
			short:       "d",
			arg:         "",
			desc:        "Use first device for log input",
			defaultVal:  false,
			configField: &config.useDevice,
		},
		{
			long:        "emulator",
			short:       "e",
			arg:         "",
			desc:        "Use first emulator for log input",
			defaultVal:  false,
			configField: &config.useEmulator,
		},
		{
			long:        "version",
			short:       "v",
			arg:         "",
			desc:        "Print version and exit",
			defaultVal:  false,
			configField: &config.version,
		},
	}

	for _, opt := range options {
		registerFlag(opt)
	}

	flag.Usage = func() {
		fmt.Println(colorize("\nUsage:", Blue, ""), colorize("loggo [options] [package ...]", Yellow, ""))
		fmt.Println(colorize("\nOptions:", Yellow, ""))

		for _, opt := range options {

			printHelpLine(opt)
		}
	}

	flag.Parse()
	config.packages = flag.Args()
	return config

}
