package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	title = flag.String("t", "Utility Name", "Title of notification.")
	mesg  = flag.String("m", "Done!", "Message notification will display.")
	sound = flag.String("s", "Ping", "Sound to play when notified.")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage: noti [-tms] <utility> [args...]

    -t    Title of notification. Default is the utility name.

    -m    Message notification will display. Default is "Done!"

    -s    Sound to play when notified. Default is Ping. Possible options
          are Basso, Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping,
          Pop, Purr, Sosumi, Submarine, Tink. Check /System/Library/Sounds
          for more info.

    -h    Display usage information and exit.`)
	}
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}

	if *title == "Utility Name" {
		*title = flag.Args()[0]
	}

	if err := Exec(flag.Args()[0], flag.Args()[1:]); err != nil {
		Notify(*title, "Failed. See terminal.", "Basso")
		return
	}

	if err := Notify(*title, *mesg, *sound); err != nil {
		log.Println(err)
	}
}

// Exec executes a program. The stdout and stderr of the executing program are
// sent to noti.
func Exec(bin string, args []string) error {
	cmd := exec.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Notify displays a notification in OS X's notification center with a given
// title, message, and sound.
func Notify(title, mesg, sound string) error {
	script := fmt.Sprintf("display notification %q with title %q sound name %q",
		mesg, title, sound)

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}
