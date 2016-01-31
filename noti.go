package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const usageText = `NOTI
    display a notification after a terminal process finishes.

USAGE
    noti [options] [utility [args...]]

OPTIONS
    -t, -title
        Set the notification title. If no arguments passed, default is "noti",
        otherwise default is utility name.

    -m, -message
        Set notification message. Default is "Done!"

    -s, -sound
        Set notification sound (OS X only). Default is Ping.
        Possible options are Basso, Blow, Bottle, Frog, Funk, Glass, Hero,
        Morse, Ping, Pop, Purr, Sosumi, Submarine, Tink.
        Check /System/Library/Sounds for available sounds.

    -f, -foreground
        Bring the terminal to the foreground. (OS X only)

    -p, -pushbullet
        Send a Pushbullet notification. Access token must be set in NOTI_PB
        environment variable.

		-c  -cluster
				Don't send a notification using notify-send or osascript. For usage
				in environments w/o graphical user interface.

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

const (
	pushbulletEnv = "NOTI_PB"
	pushbulletAPI = "https://api.pushbullet.com/v2/pushes"
)

func main() {
	foreground := flag.Bool("f", false, "")
	title := flag.String("t", "noti", "")
	mesg := flag.String("m", "Done!", "")
	sound := flag.String("s", "Ping", "")
	pbullet := flag.Bool("p", false, "")
	cluster := flag.Bool("c", false, "")
	version := flag.Bool("v", false, "")
	help := flag.Bool("h", false, "")
	flag.BoolVar(foreground, "foreground", false, "")
	flag.StringVar(title, "title", "noti", "")
	flag.StringVar(mesg, "message", "Done!", "")
	flag.StringVar(sound, "sound", "Ping", "")
	flag.BoolVar(pbullet, "pushbullet", false, "")
	flag.BoolVar(cluster, "cluster", false, "")
	flag.BoolVar(version, "version", false, "")
	flag.BoolVar(help, "help", false, "")
	flag.Usage = func() { log.Println(usageText) }
	flag.Parse()

	if *help {
		fmt.Println(usageText)
		return
	}

	if *version {
		fmt.Println("noti version 1.3.1")
		return
	}

	// noti called with a utility
	if utilArgs := flag.Args(); len(utilArgs) > 0 {
		// "noti" is default, so flag probably wasn't passed.
		if *title == "noti" {
			// title = utility's name
			*title = utilArgs[0]

			// Also, show flag or subcommand, if long enough.
			if len(utilArgs) > 1 {
				*title += " " + utilArgs[1]
			}
		}

		// run a binary and its arguments
		if err := run(utilArgs[0], utilArgs[1:]); err != nil {
			notify(*title, "Failed. See terminal.", "Basso", *foreground, *pbullet, *cluster)
			os.Exit(1)
		}
	}

	if err := notify(*title, *mesg, *sound, *foreground, *pbullet, *cluster); err != nil {
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

// pbulletNotify sends a Pushbullet notification to all devices associated with
// a given access token.
func pbulletNotify(title, mesg string) error {
	apiKey := os.Getenv(pushbulletEnv)
	if apiKey == "" {
		return errors.New("Pushbullet access token is not set in environment")
	}

	payload := bytes.NewBuffer([]byte(
		fmt.Sprintf(`{"body":"%s","title":"%s","type":"note"}`, mesg, title),
	))

	req, err := http.NewRequest("POST", pushbulletAPI, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Access-Token", apiKey)
	req.Header.Set("Content-Type", "application/json")

	if _, err = http.DefaultClient.Do(req); err != nil {
		return err
	}

	return nil
}
