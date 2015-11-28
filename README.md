# noti

Display a notification after a terminal process finishes.

## Types

These are the different types of notifications currently available.

### Desktop notifications

This is the default. Supported on OS X and Linux/FreeBSD through libnotify. This
is great so you don't have to keep checking the terminal to see if your process
is done.

![OS X notification](https://raw.githubusercontent.com/variadico/noti/master/screenshots/osx.png)
(OS X)

![Linux Mint notification](https://raw.githubusercontent.com/variadico/noti/master/screenshots/linux_mint.png)
(Linux Mint 17.2)

![Dunst](https://raw.githubusercontent.com/variadico/noti/master/screenshots/bsd_dunst.png)
([Dunst](http://knopwob.org/dunst/index.html))

### Pushbullet notifications

This is great if you want to leave sight of your computer and grab some coffee.
These notifications will get sent to all your Pushbullet devices, including your
phone.

![Pushbullet notification](https://raw.githubusercontent.com/variadico/noti/master/screenshots/pushbullet.png)

![Pushbullet Android notification](https://raw.githubusercontent.com/variadico/noti/master/screenshots/pushbullet_android.png)

## Installation

If you have Go installed, then you can do this.

```
go get -u github.com/variadico/noti
```

Otherwise, you can download the standalone binary on the
[releases page](https://github.com/variadico/noti/releases/latest). Then give it
execute permissions.

```
chmod u+x noti
```

## Usage

Just put `noti` at the beginning or end of your regular commands.

```
noti [options] [utility [args...]]
```

### Options

```
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
```

#### OS X

```
-s, -sound
    Set notification sound. Default is Ping. Possible options are Basso,
    Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
    Submarine, Tink. Check /System/Library/Sounds for available sounds.
-V, -voice
    Set voice. Check System Preferences > Dictation & Speech for available
    voices.
```

#### Linux/FreeBSD

```
-i, -icon
    Set icon name.
```

## Examples

Display a desktop notification when `tar` finishes compressing files.

```
noti tar -cjf music.tar.bz2 Music/
```

Display a notification when `apt-get` finishes updating on a remote server.

```
noti ssh you@server.com apt-get update
```

You can also add `noti` after a command, in case you forgot at the beginning.

```
clang foo.c -Wall -lm -L/usr/X11R6/lib -lX11 -o bizz; noti
```

Create a reminder to get back to a friend.

```
noti -t "Reply to Pedro" gsleep 5m &
```

Send a Pushbullet notification to your phone and other registered devices after
tests finish.

```
noti -p go test ./...
```

Have your Mac tell you what happened.

```
noti -V alex -m 'Ubuntu download finished.' docker pull ubuntu
```
