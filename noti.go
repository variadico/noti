package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const usageText = `NOTI
    display a notification in os x after a terminal process finishes.

USAGE
    noti [options] [utility [args...]]

OPTIONS
    -f, -foreground
        Bring the terminal to the foreground.

    -t, -title
        Set the notification title. If no arguments passed, default is "noti",
        else default is utility name.

    -m, -message
        Set notification message. Default is "Done!"

    -s, -sound
        Set notification sound. Default is Ping. Possible options are Basso,
        Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
        Submarine, Tink. Check /System/Library/Sounds for available sounds.

    -v, -version
        Print noti version and exit.

    -h, -help
        Display usage information and exit.

EXAMPLES
    Display a notification when tar finishes compressing files.
        noti tar -cjf music.tar.bz2 Music/

    Display a notification when apt-get finishes updating on a remote server.
        noti ssh you@server.com apt-get update

    Set the notification title to "homebrew" and message to "Up to date" and
    display it after Homebrew finishes updating.
        noti -t "homebrew" -m "up to date" brew update

    You can also add noti after a command, in case you forgot at the beginning.
        clang foo.c -Wall -lm -L/usr/X11R6/lib -lX11 -o bizz; noti
`

const (
	activateReopen = `tell application "Terminal"
	activate
	reopen
end tell`

	displayNotification = "display notification %q with title %q sound name %q"
)

func main() {
	foreground := flag.Bool("f", false, "")
	title := flag.String("t", "", "")
	mesg := flag.String("m", "Done!", "")
	sound := flag.String("s", "Ping", "")
	version := flag.Bool("v", false, "")
	help := flag.Bool("h", false, "")
	flag.BoolVar(foreground, "foreground", false, "")
	flag.StringVar(title, "title", "", "")
	flag.StringVar(mesg, "message", "Done!", "")
	flag.StringVar(sound, "sound", "Ping", "")
	flag.BoolVar(version, "version", false, "")
	flag.BoolVar(help, "help", false, "")
	flag.Usage = func() { log.Println(usageText) }
	flag.Parse()

	if *help {
		fmt.Println(usageText)
		return
	}

	if *version {
		fmt.Println("noti version 1.1.0")
		return
	}

	// noti called by itself
	if len(flag.Args()) == 0 {
		if err := notify("noti", *mesg, *sound, *foreground); err != nil {
			log.Fatal(err)
		}

		return
	}

	if *title == "" {
		// title = utility's name
		*title = flag.Args()[0]
	}

	// run a binary and its arguments
	if err := run(flag.Args()[0], flag.Args()[1:]); err != nil {
		notify(*title, "Failed. See terminal.", "Basso", *foreground)
		os.Exit(1)
	}

	if err := notify(*title, *mesg, *sound, *foreground); err != nil {
		log.Fatal(err)
	}
}

// run executes a program and waits for it to finish. The stdin, stdout, and
// stderr of noti are passed to the program.
func run(bin string, args []string) error {
	cmd := exec.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// notify displays a notification in OS X's notification center with a given
// title, message, and sound.
func notify(title, mesg, sound string, foreground bool) error {
	if foreground {
		cmd := exec.Command("osascript", "-e", activateReopen)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	script := fmt.Sprintf(displayNotification, mesg, title, sound)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}
