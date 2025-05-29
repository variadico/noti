% NOTI(1) noti 3.8.0 | Noti Manual
% variadico
% 2018/03/25

#  NAME

noti - monitor a process and trigger a notification

# SYNOPSIS

noti [flags] [utility [args...]]

# DESCRIPTION

Never sit and wait for some long-running process to finish. Noti can alert you
when it's done. You can receive messages on your computer or phone.

# OPTIONS

-t \<string\>, \--title \<string\>
: Set notification title. Default is utility name.

-m \<string\>, \--message \<string\>
: Set notification message. Default is "Done!". Read from stdin with "-".

-b, \--banner
: Trigger a banner notification. This is enabled by default. To disable this
  service, set this flag to false. This will be either `nsuser`, `freedesktop`,
  or `notifyicon` notification, depending on the OS.

-s, \--speech
: Trigger a speech notification. This will be either `say`, `espeak`, or
  `speechsynthesizer` notification, depending on the OS.

-c, \--bearychat
: Trigger a BearyChat notification. This requires `bearychat.incomingHookURI` to
  be set.

--keybase
: Trigger a Keybase notification. This requires `keybase.conversation` to
  be set.

-p, \--pushbullet
: Trigger a Pushbullet notification. This requires `pushbullet.accessToken` to
  be set.

-o, \--pushover
: Trigger a Pushover notification. This requires `pushover.apiToken` and
  `pushover.userKey` to be set.

-u, \--pushsafer
: Trigger a Pushsafer notification. This requires `pushsafer.key` to be set.

-l, \--simplepush
: Trigger a Simplepush notification. This requires `simplepush.key` to be set.

-n,\--gchat
: Trigger a Google Chat notification. This requires `gchat.appurl` to be set.

-i,\--chanify
: Trigger a Chanify notification. This requires `chanify.channelURL` to be set.

-k, \--slack
: Trigger a Slack notification. This requires `slack.appurl` (for Slack apps)
  or `slack.token` and `slack.channel` (for legacy tokens) to be set.

--twilio
: Trigger a Twilio notification. This requires `twilio.authToken`, `twilio.accountSid`, `twilio.numberFrom` and `twilio.numberTo` to be set.

--ntfy
: Trigger a ntfy notification.  This requires `ntfy.topic` be set.  Optionally, `ntfy.url` can also be set to use a different Ntfy server. For private Ntfy topics, access token authentication can be provided via `ntfy.token`. 

-w <pid>, \--pwatch <pid>
: Monitor a process by PID and trigger a notification when the pid disappears.

-f, \--file
: Path to `noti.yaml` configuration file.

\--verbose
: Enable verbose mode.

-v, \--version
: Print `noti` version and exit.

-h, \--help
: Print `noti` help and exit.

# ENVIRONMENT

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
* `NOTI_NTFY_TOPIC`
* `NOTI_NTFY_URL`
* `NOTI_NTFY_TOKEN`
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


# FILES

If not explicitly set with \--file, then noti will check the following paths,
in the following order.

* ./.noti.yaml
* $XDG_CONFIG_HOME/noti/noti.yaml

# EXAMPLES

Display a notification when `tar` finishes compressing files.

    noti tar -cjf music.tar.bz2 Music/

Add noti after a command, in case you forgot at the beginning.

    clang foo.c -Wall -lm -L/usr/X11R6/lib -lX11 -o bizz; noti

If you already started a command, but forgot to use `noti`, then you can do
this to get notified when that process' PID disappears.

    noti --pwatch $(pgrep docker-machine)

Receive your message from stdin with `-`.

    rsync -az --stats ~/  server:/backups/homedir | noti -t "backup stats" -m -

# REPORTING BUGS

Report bugs on GitHub at https://github.com/variadico/noti/issues.

# SEE ALSO

noti.yaml(5)
