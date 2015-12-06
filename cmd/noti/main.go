package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

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
    -s, -save
        Save current flag set.
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
	flagState = ".noticonfig.json"
)

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
	save := flag.Bool("S", false, "")
	flag.BoolVar(save, "save", false, "")
	flag.Usage = func() { log.Println(usageText) }
	flag.Parse()

	if *save {
		if err := saveFlags(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := loadFlags(); err != nil {
			log.Fatal(err)
		}
	}

	if *help {
		fmt.Println(usageText)
		return
	}

	if *version {
		fmt.Println("noti version v2dev")
		return
	}

	processArgs(title, message, flag.Args())

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

func saveFlags() error {
	config := make(map[string]string)
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "save" || f.Name == "S" {
			return
		}
		config[f.Name] = f.Value.String()
	})

	bs, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(usr.HomeDir, flagState), bs, 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadFlags() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	configPath := filepath.Join(usr.HomeDir, flagState)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	bs, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	config := make(map[string]string)
	err = json.Unmarshal(bs, &config)
	if err != nil {
		return err
	}

	flag.VisitAll(func(f *flag.Flag) {
		if v, exists := config[f.Name]; exists {
			f.Value.Set(v)
		}
	})

	return nil
}

// processArgs executes any utility args. In addition, it also sets the
// notification title and message, as necessary.
func processArgs(title, message *string, args []string) {
	if len(args) < 1 {
		return
	}

	if err := run(args); err != nil {
		*message = fmt.Sprint(err)
	}

	if *title == "noti" {
		*title = genTitle(args, 2)
	}
}

// run executes a program and waits for it to finish. The stdin, stdout, and
// stderr of noti are passed to the program.
func run(args []string) error {
	var cmd *exec.Cmd

	if ln := len(args); ln < 1 {
		return nil
	} else if ln == 1 {
		cmd = exec.Command(args[0])
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}

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

// genTitle takes a list of arguments and constructs a title by using the first
// 2 arguments.
func genTitle(args []string, max int) string {
	var t string

	for i, a := range args {
		t += a + " "

		if i+1 >= max {
			break
		}
	}

	return t
}
