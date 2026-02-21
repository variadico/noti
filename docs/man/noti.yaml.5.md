% NOTI.YAML(5) noti 3.8.0 | Noti Configuration File Format
% variadico
% 2018/03/25

#  NAME

noti.yaml - noti configuration file

# SYNOPSIS

noti.yaml

# DESCRIPTION

File format is YAML.

If not explicitly set with \--file, then noti will check the following paths,
in the following order.

* ./.noti.yaml
* $XDG_CONFIG_HOME/noti/noti.yaml

If $XDG_CONFIG_HOME is empty, then $HOME/.config will be used as its default
value and noti will check $HOME/.config/noti/noti.yaml.

# BANNER

icon
: Path to notification icon image. On macOS, accepts PNG or JPEG. On Linux,
  accepts an image path or a freedesktop icon theme name. On Windows, accepts
  an .ico file path.

# NSUSER

soundName
: Banner success sound. Default is Ping. Possible options are Basso, Blow,
  Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
  Submarine, Tink. See /System/Library/Sounds for available sounds.

soundNameFail
: Banner failure sound. Default is Basso. Possible options are Basso,
  Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
  Submarine, Tink. See /System/Library/Sounds for available sounds.

# SAY

voice
: Name of voice used for speech notifications.

# ESPEAK

voiceName
: Name of voice used for speech notifications.

# SPEECHSYNTHESIZER

voice
: Name of voice used for speech notifications.

# BEARYCHAT

incomingHookURI
: BearyChat incoming URI.

# KEYBASE

conversation
: Keybase message destination. Can be either users (comma-separated) or team.

channel
: Keybase team's chat channel to send to. Conversation must be a team.
  If empty, the team's default channel will be used (typically "general").

explodingLifetime
: Keybase self-destructing message, after the specified time. Times are
  written like `30s` (30 seconds), `15m` (15 minutes), `24h` (24 hours).

public
: Enables broadcasting a message to everyone (when `conversation` is
  your username), or to teams (when `conversation` is your team name).

# PUSHBULLET

accessToken
: Pushbullet access token. Log into your Pushbullet account and retrieve a
  token from the Account Settings page.

deviceIden
: Pushbullet device iden of the target device, if sending to a single device.

# PUSHOVER

apiToken
: Pushover access token. Log into your Pushover account and create a
  token from the Create New Application/Plugin page.

userKey
: Pushover message destination. Should be your User Key.

# PUSHSAFER

key
: Pushsafer private or alias key. Log into your Pushsafer account and note
  your private or alias key.

# SIMPLEPUSH

key
: Simplepush key. Install the Simplepush app and retrieve your key there.

event
: Customize ringtone and vibration.

# SLACK

token
: Slack access token. Log into your Slack account and retrieve a token
  from the Slack Web API page.

channel
: Slack message destination. Can be either a #channel or a @username.

username
: Noti bot username.

# TWILIO

AuthToken
: Twilio access token. Log into your Twilio account and copy the AuthToken from your project dashboard

accountSid
: Twilio account id. Log into your Twilio account and copy the accountSid from your project dashboard.

numberTo
: This parameter determines the destination phone number for your SMS message. Format this number with a '+' and a country code, e.g., +16175551212

numberFrom
: From specifies the Twilio phone number, short code, or Messaging Service that sends this message. This must be a Twilio phone number that you own, formatted with a '+' and country code, e.g. +16175551212 (E.164 format)


# GCHAT

appurl
: This parameter defines the URL for the Google Chat webhook.

template
: This parameter defines the template combining the title and the message. The default is: '*{{.title}}*: {{.message}}'


# EXAMPLES

    ---
    banner:
      icon: /path/to/icon.png
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
      appurl: 'https://hooks.slack.com/services/xxx/yyy/zzz'
    twilio:
      numberto: +972542877978
      numberfrom: +18111119711
      accountsid: AC3cd135aa82XXXXXXXXf792ba23fc98
      authtoken: 74efd0bXXXXXXXXXXX32f7daca
	gchat:
	  appurl: 'https://chat.googleapis.com/v1/spaces/example/messages?key=keyexample'
	  template: '*{{.title}}*: {{.message}}'



# SEE ALSO

noti(1)
