# vbs

```
PACKAGE DOCUMENTATION

package vbs
    import "github.com/variadico/vbs"

    Package vbs prints text to a writer if Enabled is true. This package
    isn't threadsafe. Use a mutex if you're going to be changing Enabled in
    goroutines.

VARIABLES

var (
    // Enabled indicates whether or not something should be printed.
    Enabled bool
)

FUNCTIONS

func Printf(format string, a ...interface{})
    Printf prints to output if Enabled is true.

func Println(a ...interface{})
    Println prints to output if Enabled is true.

func SetOutput(w io.Writer)
    SetOutput sets the output for the global printer. By default, the global
    output is os.Stdout.

TYPES

type Printer struct {
    Enabled bool
    // contains filtered or unexported fields
}
    Printer is a conditional printer that prints if Enabled is true.

func New(out io.Writer) Printer
    New returns a new verbose Printer.

func (p Printer) Printf(format string, a ...interface{})
    Printf prints to output if Enabled is true.

func (p Printer) Println(a ...interface{})
    Println prints to output if Enabled is true.
```
