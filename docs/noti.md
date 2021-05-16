# Noti Manual

##  Name

noti - monitor a process and trigger a notification

## Synopsis

```
noti [flags] [utility [args...]]
```

## Description

Never sit and wait for some long-running process to finish. Noti can alert you
when it's done. You can receive messages on your computer or phone.

## Services

Noti can send notifications on a number of services.

```
           | macOS | Linux | Windows
--------------------------------------
Banner     |   ✔   |   ✔   |    ✔
Speech     |   ✔   |   ✔   |    ✔
BearyChat  |   ✔   |   ✔   |    ✔
Keybase    |   ✔   |   ✔   |    ✔
Pushbullet |   ✔   |   ✔   |    ✔
Pushover   |   ✔   |   ✔   |    ✔
Pushsafer  |   ✔   |   ✔   |    ✔
Simplepush |   ✔   |   ✔   |    ✔
Slack      |   ✔   |   ✔   |    ✔
Twilio     |   ✔   |   ✔   |    ✔
```

## Installation

The `master` branch always contains the latest tagged release.

```shell
# Install the latest version on any platform.
go get -u github.com/variadico/noti/cmd/noti

# Install the latest version on macOS.
brew install noti
```

If you don't want to build from source or install anything extra, just download
the latest binary.

```shell
# macOS
curl -L $(curl -s https://api.github.com/repos/variadico/noti/releases/latest | awk '/browser_download_url/ { print $2 }' | grep 'darwin-amd64' | sed 's/"//g') | tar -xz

# Linux
curl -L $(curl -s https://api.github.com/repos/variadico/noti/releases/latest | awk '/browser_download_url/ { print $2 }' | grep 'linux-amd64' | sed 's/"//g') | tar -xz
```

## Options

```
-t <string>, --title <string>
    Set notification title.  Default is utility name.

-m <string>, --message <string>
    Set notification message.  Default is "Done!". Read from stdin with "-".

-e, --time
    Show command execution time in message.

-b, --banner
    Trigger a banner notification.  This is enabled by default.  To disable
    this service, set this flag to false.  This will be either nsuser,
    freedesktop, or notifyicon notification, depending on the OS.

-s, --speech
    Trigger a speech notification.  This will be either say, espeak, or
    speechsynthesizer notification, depending on the OS.

-c, --bearychat
    Trigger a BearyChat notification.  This requires
    bearychat.incomingHookURI to be set.

--keybase
    Trigger a Keybase notification.  This requires keybase.conversation
    to be set.

-p, --pushbullet
    Trigger a Pushbullet notification.  This requires pushbullet.accessToken
    to be set.

-o, --pushover
    Trigger a Pushover notification.  This requires pushover.apiToken and
    pushover.userKey to be set.

-u, --pushsafer
    Trigger a Pushsafer notification.  This requires pushsafer.key to be set.

-l, --simplepush
    Trigger a Simplepush notification.  This requires simplepush.key to be
    set.

-k, --slack
    Trigger a Slack notification.  This requires slack.token and
    slack.channel to be set.

-g, --telegram
    Trigger a Telegram notification.  This requires telegeram.token and
    telegram.chatId to be set.

--twilio
    Trigger a Twilio notification.  This requires twilio.authToken, twilio.accountSid, 
    twilio.numberFrom and twilio.numberTo to be set.

-w <pid>, --pwatch <pid>
    Monitor a process by PID and trigger a notification when the pid
    disappears.

-f, --file
    Path to noti.yaml configuration file.

--verbose
    Enable verbose mode.

-v, --version
    Print noti version and exit.

-h, --help
    Print noti help and exit.
```

## Environment

* `NOTI_DEFAULT`
* `NOTI_NSUSER_SOUNDNAME`
* `NOTI_NSUSER_SOUNDNAMEFAIL`
* `NOTI_SAY_VOICE`
* `NOTI_ESPEAK_VOICENAME`
* `NOTI_SPEECHSYNTHESIZER_VOICE`
* `NOTI_BEARYCHAT_INCOMINGHOOKURI`
* `NOTI_KEYBASE_CONVERSATION`
* `NOTI_KEYBASE_CHANNEL`
* `NOTI_KEYBASE_PUBLIC`
* `NOTI_KEYBASE_EXPLODINGLIFETIME`
* `NOTI_PUSHBULLET_ACCESSTOKEN`
* `NOTI_PUSHBULLET_DEVICEIDEN`
* `NOTI_PUSHOVER_APITOKEN`
* `NOTI_PUSHOVER_USERKEY`
* `NOTI_PUSHSAFER_KEY`
* `NOTI_SIMPLEPUSH_KEY`
* `NOTI_SIMPLEPUSH_EVENT`
* `NOTI_SLACK_TOKEN`
* `NOTI_SLACK_CHANNEL`
* `NOTI_SLACK_USERNAME`
* `NOTI_TWILIO_TO`
* `NOTI_TWILIO_FROM`
* `NOTI_TWILIO_ACCOUNTSID`
* `NOTI_TWILIO_AUTHTOKEN`
* `NOTI_TELEGRAM_CHATID`
* `NOTI_TELEGRAM_TOKEN`


## Files

If not explicitly set with `--file`, then `noti` will check the following paths,
in the following order.

* `./.noti.yaml`
* `$XDG_CONFIG_HOME/noti/noti.yaml`

If `$XDG_CONFIG_HOME` is empty, then `$HOME/.config` will be used as its default
value and `noti` will check `$HOME/.config/noti/noti.yaml`.

## Configuration

```
NSUSER

soundName
    Banner success sound. Default is Ping. Possible options are Basso, Blow,
    Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
    Submarine, Tink. See /System/Library/Sounds for available sounds.

soundNameFail
    Banner failure sound. Default is Basso. Possible options are Basso,
    Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
    Submarine, Tink. See /System/Library/Sounds for available sounds.

SAY

voice
    Name of voice used for speech notifications.

ESPEAK

voiceName
    Name of voice used for speech notifications.

SPEECHSYNTHESIZER

voice
    Name of voice used for speech notifications.

BEARYCHAT

incomingHookURI
    BearyChat incoming URI.

KEYBASE

conversation
    Keybase message destination. Can be either users (comma-separated) or team.

channel
    Keybase team's chat channel to send to. Conversation must be a team.
    If empty, the team's default channel will be used (typically "general").

explodingLifetime
    Keybase self-destructing message, after the specified time. Times are
    written like `30s` (30 seconds), `15m` (15 minutes), `24h` (24 hours).

public
    Enables broadcasting a message to everyone (when `conversation` is
    your username), or to teams (when `conversation` is your team name).

PUSHBULLET

accessToken
    Pushbullet access token. Log into your Pushbullet account and retrieve a
    token from the Account Settings page.

deviceIden
    Pushbullet device iden of the target device, if sending to a single device.

PUSHOVER

apiToken
    Pushover access token. Log into your Pushover account and create a
    token from the Create New Application/Plugin page.

userKey
    Pushover message destination. Should be your User Key.

PUSHSAFER

key
    Pushsafer private or alias key. Log into your Pushsafer account and note
    your private or alias key.

SIMPLEPUSH

key
    Simplepush key. Install the Simplepush app and retrieve your key there.

event
    Customize ringtone and vibration.

SLACK

token
    Slack access token. Log into your Slack account and retrieve a token
    from the Slack Web API page.

channel
    Slack message destination. Can be either a #channel or a @username.

username
    Noti bot username.

TELEGRAM

token
    Telegram access token. The token can be retrieved using the [BotFather](https://telegram.me/botfather).

chatId
    Telegram message destination: Can be either a chat or a channel.

TWILIO

AuthToken
    Twilio access token. Log into your Twilio account and copy the AuthToken from your 
    project dashboard.

accountSid
    Twilio account id. Log into your Twilio account and copy the accountSid from your 
    project dashboard.

numberTo
    This parameter determines the destination phone number for your SMS message. 
    Format this number with a '+' and a country code, e.g., +16175551212

numberFrom
    From specifies the Twilio phone number, short code, or Messaging Service that sends 
    this message.
    This must be a Twilio phone number that you own,
    formatted with a '+' and country code, e.g. +16175551212 (E.164 format)
```

## Examples

Display a notification when `tar` finishes compressing files.

```
noti tar -cjf music.tar.bz2 Music/
```

Add `noti` after a command, in case you forgot at the beginning.

```
clang foo.c -Wall -lm -L/usr/X11R6/lib -lX11 -o bizz; noti
```

If you already started a command, but forgot to use `noti`, then you can do
this to get notified when that process' PID disappears.

```
noti --pwatch $(pgrep docker-machine)
```

Receive your message from stdin with `-`.

```
rsync -az --stats ~/  server:/backups/homedir | noti -t "backup stats" -m -
```

Sample configuration file.

```yaml
nsuser:
  soundName: Ping
  soundNameFail: Basso
say:
  voice: Alex
espeak:
  voiceName: english-us
speechsynthesizer:
  voice: Microsoft David Desktop
bearychat:
  incomingHookURI: 1234567890abcdefg
keybase:
  conversation: yourteam
  channel: general
pushbullet:
  accessToken: 1234567890abcdefg
  deviceIden: 1234567890abcdefg
pushover:
  userKey: 1234567890abcdefg
  apiToken: 1234567890abcdefg
pushsafer:
  key: 1234567890abcdefg
simplepush:
  key: 1234567890abcdefg
  event: 1234567890abcdefg
slack:
  token: 1234567890abcdefg
  channel: '@jaime'
  username: noti
telegram:
  token: '110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw'
  chatId: 1234567890
twilio:
  numberto: +972542877978
  numberfrom: +18111119711
  accountsid: 1234567890abcdefg
  authtoken: 1234567890abcdefg

```

## Setting up cloud accounts

### BearyChat

Log into your BearyChat account. Then create an [incoming robot][bc-incoming].
Next, look for the "Hook Address" (or "Hook 地址" in Chinese), this is what
you'll set `bearychat.incomingHookURI` to.

### Keybase

Log into your [Keybase] account and go to the Chat section.
Pick a team or user to send to, and set your `keybase.conversation` to that
name. For teams, you can also specify a channel with `keybase.channel` or have
it send to `general`. You can also send notifications that auto-delete with
`keybase.explodingLifetime` set between `30s` (seconds) to `168h` (hours).

As long as the keybase service is running and the `keybase` binary is in your
system path, you should be all set.

### Pushbullet

Log into your Pushbullet account. Next, click on [Settings] on the left sidebar.
Scroll down to "Access Tokens" and click "Create Access Token". The text that
appears will be what you'll set `pushbullet.accessToken` to.

### Pushover

Log into your [Pushover] account. Next, look for the "User Key". That's what
you'll set `pushover.userKey` to.

Next [create a new application]. Fill in the fields. Under "Type", select
"Script". Finally, go to the application page. Look for "API Token/Key". This is
what you'll set `pushover.apiToken` to.

### Pushsafer

Log into your [Pushsafer] account. Next, look for the "Private or Alias Key".
That's what you'll set `pushsafer.key` to.

### Simplepush

Install the Simplepush Android app to get your Simplepush key.
That's the key you'll set to `simplepush.key`.
Simplepush requires no registration and sending notifications is completely free.

In the app you can create events to customize ringtone and vibration patterns for
different kinds of notifications.
The event id you can set in the app, translates to `simplepush.event` in noti.

### Slack

Log into your Slack account. Then go to the [OAuth Tokens for Testing and
Development] page. Create a token. This is what you'll set `slack.token` to.

The variable `slack.channel` can be set to a channel like `#general` or
`#random`. You can also set it to someone's username, like `@jaime` or
`@variadico`.

### Telegram

Open your telegram app, and start a conversation with the [BotFather](https://telegram.me/botfather).
Use the `/newbot` command to get the BotFather to create your bot, and lastly copy and keep the bot token.
Send the new bot a message from the chat you want to notify, then go to `https://api.telegram.org/bot<TOKEN>/getUpdates` and find the `id` for the chat.

### Twilio

Log into your [Twilio] account. Next go to the "Project info" tab. look for "ACCOUNT SID" that's what you will set for `twilio.accountsid` and look for "AUTH TOKEN" that's what you will set for `twilio.authtoken`

## Reporting bugs

Report bugs on GitHub at https://github.com/variadico/noti/issues.


[Settings]: https://www.pushbullet.com/#settings
[Pushover]: https://pushover.net
[create a new application]: https://pushover.net/apps/build
[Pushsafer]: https://www.pushsafer.com
[OAuth Tokens for Testing and Development]: https://api.slack.com/docs/oauth-test-tokens
[bc-incoming]: https://bearychat.com/integrations/incoming
[Keybase]: https://keybase.io/download
[Twilio]: https://www.twilio.com/console
