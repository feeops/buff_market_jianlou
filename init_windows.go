package main

import (
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	winConsole := windows.Handle(os.Stdin.Fd())

	var mode uint32
	err := windows.GetConsoleMode(winConsole, &mode)
	if err != nil {
		panic(err)
	}

	// Disable this mode
	mode &^= windows.ENABLE_QUICK_EDIT_MODE

	// Enable this mode
	mode |= windows.ENABLE_EXTENDED_FLAGS

	err = windows.SetConsoleMode(winConsole, mode)
	if err != nil {
		panic(err)
	}
}
