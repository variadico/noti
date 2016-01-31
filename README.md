# noti

Trigger notifications when a process completes.

Never sit and wait for some long-running process to finish! `noti` will alert
you when it's done—on your computer or smartphone—so you can stop worrying about
constantly checking the terminal.

## Types

### Banner notifications

On OS X, these are the normal app notifications you already know. On Linux and
FreeBSD, these are libnotify notifications and requre `notify-send` to be
installed.

### Speech notifications

On OS X, these use the built-in speech command. On Linux and FreeBSD, these use
`espeak` and requre `espeak` to be installed.

### Pushbullet notifications

These are sent to all devices registered with Pushbullet. You need to create a
Pushbullet account and get an [access token][1].

### Slack notifications

These are sent to all devices that have the Slack app installed. You need to
create a Slack account and get an [access token][2].

## Installation

```
# Anywhere
go get -u github.com/variadico/noti

# Alternatively on OS X
brew install noti
```

## Usage
Put `noti` at the beginning or end of your regular commands.

```
noti [options] [utility [args...]]
```

### Options

```
-t, -title
    Notification title. Default is utility name.
-m, -message
    Notification message. Default is "Done!"
-b, -banner
    Trigger a banner notification. Default is true. To disable this
    notification set this flag to false.
-s, -speech
    Trigger a speech notification. Optionally, customize the voice with
    NOTI_VOICE.
-p, -pushbullet
    Trigger a Pushbullet notification. This requires NOTI_PUSHBULLET_TOK to
    be set.
-k, -slack
    Trigger a Slack notification. This requires NOTI_SLACK_TOK and
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
NOTI_PUSHBULLET_TOK
    Pushbullet access token. Log into your Pushbullet account and retrieve a
    token from the Account Settings page.
NOTI_SLACK_TOK
    Slack access token. Log into your Slack account and retrieve a token
    from the Slack Web API page.
NOTI_SLACK_DEST
    Slack channel to send message to. Can be either a #channel or a
    @username.
NOTI_VOICE
    Name of voice used for speech notifications.
```

#### OS X

```
NOTI_SOUND
    Banner notification sound. Default is Ping. Possible options are Basso,
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

[1]: https://www.pushbullet.com/#settings/account
[2]: https://api.slack.com/web
