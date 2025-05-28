package main

import (
	"strings"
	"syscall"
	"unsafe"
)

func rightAlign(text string, width int) string {
	if len(text) > width {
		return text[len(text)-width:]
	}
	return strings.Repeat(" ", width-len(text)) + text
}

func getTerminalWidth() int {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	ws := &winsize{}
	retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		return -1
	}
	return int(ws.Col)
}

func indentWrap(message string, width, headerSize int) string {
	if width == -1 {
		return message
	}

	message = strings.ReplaceAll(message, "\t", "    ")
	wrapArea := width - headerSize
	messagebuf := ""
	current := 0

	for current < len(message) {
		next := current + wrapArea
		if next > len(message) {
			next = len(message)
		}
		messagebuf += message[current:next]
		if next < len(message) {
			messagebuf += "\n"
			messagebuf += strings.Repeat(" ", headerSize)
		}
		current = next
	}
	return messagebuf
}
