package command

import (
	"fmt"
)

var (
	// vbsEnabled indicates whether or not something should be printed.
	vbsEnabled bool
)

// vbsPrintln prints to output if Enabled is true.
func vbsPrintln(a ...interface{}) {
	if vbsEnabled {
		fmt.Println(a...)
	}
}

// vbsPrintf prints to output if Enabled is true.
func vbsPrintf(format string, a ...interface{}) {
	if vbsEnabled {
		fmt.Printf(format, a...)
	}
}
