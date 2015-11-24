package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/variadico/noti/pushbullet"
)

const usageTmpl = `NOTI
    trigger a notification after a terminal process finishes.
USAGE
    noti [options] [utility [args...]]
OPTIONS
    -t, -title
        Set the notification title. If no arguments passed, default is "noti",
        otherwise default is utility name.
    -m, -message
        Set notification message. Default is "Done!"%s
    -p, -pushbullet
        Send a Pushbullet notification. Access token must be set in
        PUSHBULLET_ACCESS_TOKEN environment variable.
    -v, -version
        Print noti version and exit.
    -h, -help
        Display usage information and exit.
EXAMPLES
    Display a notification when tar finishes compressing files.
        noti tar -cjf music.tar.bz2 Music/
    You can also add noti after a command, in case you forgot at the beginning.
        clang foo.c -Wall -lm -L/usr/X11R6/lib -lX11 -o bizz; noti
    Create a reminder to get back to a friend.
        noti -t "Reply to Pedro" gsleep 5m &
`

var usageText string

func main() {
	title := flag.String("t", "noti", "")
	flag.StringVar(title, "title", "noti", "")
	message := flag.String("m", "Done!", "")
	flag.StringVar(message, "message", "Done!", "")
	pb := flag.Bool("p", false, "")
	flag.BoolVar(pb, "pushbullet", false, "")
	version := flag.Bool("v", false, "")
	flag.BoolVar(version, "version", false, "")
	help := flag.Bool("h", false, "")
	flag.BoolVar(help, "help", false, "")
	flag.Usage = func() { log.Println(usageText) }
	flag.Parse()

	if *help {
		fmt.Println(usageText)
		return
	}

	if *version {
		fmt.Println("noti version dev")
		return
	}

	if *pb {
		if err := pushbulletNotify(*title, *message); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := notify(*title, *message); err != nil {
			log.Fatal(err)
		}
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

// pushbulletNotify sends a pushbullet notification.
func pushbulletNotify(title, message string) error {
	nt := &pushbullet.Notification{
		AccessToken: os.Getenv(pushbullet.AccessTokenEnv),
		Title:       title,
		Body:        message,
	}

	return nt.Notify()
}
