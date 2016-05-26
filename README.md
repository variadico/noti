# noti

[![Gitter]](https://gitter.im/variadico/noti)
[![Open Hub]](https://www.openhub.net/p/variadico-noti)

Trigger notifications when a process completes.

Never sit and wait for some long-running process to finish! `noti` will alert
you when it's done—on your computer or smartphone—so you can stop worrying about
constantly checking the terminal.

![OS X Banner Notification]

## Types

### Desktop

Supported on OS X, Linux, and FreeBSD.

* Banner
* Speech

### Mobile

Supported on any platform.

* HipChat
* Pushbullet
* Pushover
* Slack

Check out the [Wiki] for more information on how to configure mobile
notifications. If you're curious, you can also browse the [screenshots]
directory.

## Installation

```
# Anywhere, latest
go get -u github.com/variadico/noti/cmd/noti

# Alternatively on OS X
brew install noti
```

If you only want a binary, get it from the [releases page].

## Usage
Put `noti` at the beginning or end of your regular commands.

```
noti [options] [utility [args...]]
```

### Options

```
-t <string>, -title <string>
    Notification title. Default is utility name.
-m <string>, -message <string>
    Notification message. Default is "Done!"
-w <pid>, -pwatch <pid>
    Trigger notification after PID disappears.

-b, -banner
    Trigger a banner notification. Default is true. To disable this
    notification set this flag to false.
-s, -speech
    Trigger a speech notification. Optionally, customize the voice with
    NOTI_VOICE.

-i, -hipchat
    Trigger a HipChat notification. Requires NOTI_HIPCHAT_TOK and
    NOTI_HIPCHAT_DEST to be set.
-p, -pushbullet
    Trigger a Pushbullet notification. Requires NOTI_PUSHBULLET_TOK to
    be set.
-o, -pushover
    Trigger a Pushover notification. Requires NOTI_PUSHOVER_TOK and
    NOTI_PUSHOVER_DEST to be set.
-k, -slack
    Trigger a Slack notification. Requires NOTI_SLACK_TOK and
    NOTI_SLACK_DEST to be set.

-v, -version
    Print noti version and exit.
-h, -help
    Display help information and exit.
```

### Environment

You can further configure `noti` by setting the following environment variables.
Some are required for specific notifications. For example, `NOTI_PUSHBULLET_TOK`
**must** be set to use Pushbullet notifications. Others can be used to
optionally change the default behavior, like setting `NOTI_VOICE` to change the
voice used in speech notifications.

```
NOTI_DEFAULT
    Notification types noti should trigger in a space-delimited list. For
    example, set NOTI_DEFAULT="banner speech pushbullet slack" to enable
    all available notifications to fire sequentially.
NOTI_HIPCHAT_TOK
    HipChat access token. Log into your HipChat account and retrieve a token
    from the Room Notification Tokens page.
NOTI_HIPCHAT_DEST
    HipChat message destination. Can be either a Room name or ID.
NOTI_PUSHBULLET_TOK
    Pushbullet access token. Log into your Pushbullet account and retrieve a
    token from the Account Settings page.
NOTI_PUSHOVER_TOK
    Pushover access token. Log into your Pushover account and create a
    token from the Create New Application/Plugin page.
NOTI_PUSHOVER_DEST
    Pushover message destination. Should be your User Key.
NOTI_SLACK_TOK
    Slack access token. Log into your Slack account and retrieve a token
    from the Slack Web API page.
NOTI_SLACK_DEST
    Slack message destination. Can be either a #channel or a @username.
NOTI_VOICE
    Name of voice used for speech notifications.
```

#### OS X only

```
NOTI_SOUND
    Banner success sound. Default is Ping. Possible options are Basso, Blow,
    Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
    Submarine, Tink. See /System/Library/Sounds for available sounds.
NOTI_SOUND_FAIL
    Banner failure sound. Default is Basso. Possible options are Basso,
    Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
    Submarine, Tink. See /System/Library/Sounds for available sounds.
```

## Examples

Display a notification when `tar` finishes compressing files.

```
noti tar -cjf music.tar.bz2 Music/
```

You can also add `noti` after a command, in case you forgot at the beginning.

```
clang foo.c -Wall -lm -L/usr/X11R6/lib -lX11 -o bizz; noti
```

If you already started a command, but forgot to use `noti`, then you can do this
to get notified when that process' PID disappears.

```
noti -pwatch $(pgrep docker-machine)
```

[OS X Banner Notification]: https://raw.githubusercontent.com/variadico/noti/master/.github/screenshots/osx_banner.png
[Wiki]: https://github.com/variadico/noti/wiki
[screenshots]: https://github.com/variadico/noti/tree/master/.github/screenshots
[releases page]: https://github.com/variadico/noti/releases
[Gitter]: https://badges.gitter.im/variadico/noti.svg
[Open Hub]: https://img.shields.io/badge/open%20hub-metrics-blue.svg
