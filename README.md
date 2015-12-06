# noti

Trigger a notification when a terminal process finishes.

## Types

These are the currently available notifications. The default type is a banner
notification.

### Banner

On OS X, `noti` uses NSUserNotification.

![OS X notification]

On Linux and FreeBSD, `noti` uses NotifyNotifications.

![Linux Mint notification]

(Linux Mint 17.2)

![Dunst notification]

([Dunst 1.1.0])

### Smartphone

On smartphones, `noti` uses Pushbullet notifications. You can get an access
token from Pushbullet, under Settings. Set the env var `PUSHBULLET_ACCESS_TOKEN`
on your computer. For example, `export PUSHBULLET_ACCESS_TOKEN='12345'`.

![Pushbullet Android notification]

### Speech

On OS X, `noti` uses NSSpeechSynthesizer. On Linux and FreeBSD, `noti` uses
eSpeak.

## Install

Download the standalone binary for your OS and architecture from the [releases
page]. Then, add it to your `PATH`.

### Alternative

#### Compile from source OS X

```
go get -u github.com/variadico/noti/cmd/noti
```

#### Compile from source Ubuntu

```
sudo apt-get install libespeak-dev libnotify-dev
go get -u github.com/variadico/noti/cmd/noti
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
    Set notification message. Default is "Done!"
-p, -pushbullet
    Send a Pushbullet notification. Access token must be set in
    PUSHBULLET_ACCESS_TOKEN environment variable.
-s, -save
    Save current flag set.
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
    Set voice.
```

#### Linux and FreeBSD

```
-i, -icon
    Set icon name. You can pass a name from /usr/share/icons/gnome/32x32/ or
    /usr/share/notify-osd/icons/. Alternatively, you can specify a full
    filepath.
-V, -voice
    Set voice.
```

## Quick Walkthrough

The shortest usage of `noti` is this.

```
noti
```

That will display a banner notification that says "noti Done!".

### Banner notification

A more practical example could be compressing your huge music collection. You
don't want to sit and stare at your terminal while `tar` compresses all your
files. You also don't have to have to `alt+tab` back and forth to check if `tar`
is done. Instead, you can use `noti` to notify you when your music library is
compressed. You can relax and focus on that cat video.

```
noti tar -cjf music.tar.bz2 Music/
```

### Smartphone notification

Here's another situation where `noti` can help you save time. Imagine you're
working on a huge software project. Builds take 10s of minutes, not to mention
tests take more than 1 hour (true story). You don't want to sit at your computer
twiddling your thumbs. Maybe you should go play Ping Pong while you wait. You
can use `noti` to send you a Pushbullet notification when your project finishes
building.

```
export PUSHBULLET_ACCESS_TOKEN='abc123'
noti -p ./build.sh
```

## More Examples

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

Have your Mac tell you what happened. This will speak "docker pull done".

```
noti -V zarvox docker pull ubuntu
```

[releases page]: https://github.com/variadico/noti/releases/latest
[OS X notification]: https://raw.githubusercontent.com/variadico/noti/master/screenshots/osx.png
[Linux Mint notification]: https://raw.githubusercontent.com/variadico/noti/master/screenshots/linux_mint.png
[Dunst notification]: https://raw.githubusercontent.com/variadico/noti/master/screenshots/bsd_dunst.png
[Dunst 1.1.0]: http://knopwob.org/dunst/index.html
[Pushbullet Android notification]: https://raw.githubusercontent.com/variadico/noti/master/screenshots/pushbullet_android.png
