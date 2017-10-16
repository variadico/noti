// Package vbs prints text to a writer if Enabled is true. This
// package isn't threadsafe. Use a mutex if you're going to be changing
// Enabled in goroutines.
package vbs

import (
	"fmt"
	"io"
	"os"
)

var (
	// Enabled indicates whether or not something should be printed.
	Enabled bool

	// output is the destination for the global printer.
	output io.Writer = os.Stdout
)

// SetOutput sets the output for the global printer. By default, the global
// output is os.Stdout.
func SetOutput(w io.Writer) {
	output = w
}

// Println prints to output if Enabled is true.
func Println(a ...interface{}) {
	if Enabled {
		fmt.Fprintln(output, a...)
	}
}

// Printf prints to output if Enabled is true.
func Printf(format string, a ...interface{}) {
	if Enabled {
		fmt.Fprintf(output, format, a...)
	}
}

// Printer is a conditional printer that prints if Enabled is true.
type Printer struct {
	Enabled bool
	output  io.Writer
}

// New returns a new verbose Printer.
func New(out io.Writer) Printer {
	return Printer{
		Enabled: false,
		output:  out,
	}
}

// Println prints to output if Enabled is true.
func (p Printer) Println(a ...interface{}) {
	if p.Enabled {
		fmt.Fprintln(p.output, a...)
	}
}

// Printf prints to output if Enabled is true.
func (p Printer) Printf(format string, a ...interface{}) {
	if p.Enabled {
		fmt.Fprintf(p.output, format, a...)
	}
}
