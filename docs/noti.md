# Noti manual

## Configuring command

Noti uses environment variables for configuration. Here are a couple of ways of
setting those.

### One run

Set a variable for only 1 execution of `noti`. The next time you run `noti` this
variable will no longer take effect.

```
NOTI_SOUND=hero noti
```

### One session

Set a variable for an entire shell session. Each time you run `noti` the
variable will be in effect. However, after you close your terminal, the variable
will be unset.

```
export NOTI_SOUND=hero
noti
```

### Forever

Set the variable in your shell `*rc` file. Whatever values you set here will
persist across all shell sessions, until you remove the variable from your shell
`rc` file.

```
# .bashrc, .zshrc, or whatever shell rc you use.
echo "export NOTI_SOUND=hero" >> ~/.bashrc
noti
```

### Config file

Noti can also read a `.noti.yaml` file. It checks your current and home
directory for this file. These are the current options that can be set in the
config file. Not all notification types are supported on every platform.

```yaml
---
nsuser:
  soundName: Ping
  soundNameFail: Basso
```

## Configuring cloud services

### HipChat

Log into your HipChat account. Go to My Account > Rooms > {pick a room} >
Tokens. Create a new token. Set the Scope to "Send Notification". That's what
you'll set `NOTI_HIPCHAT_TOK` to.

Next, go to My Account > Rooms > {pick a room} > Summary. Look for "API ID". You
can set `NOTI_HIPCHAT_DEST` to "API ID" or you can use the Room name, like
"MyRoom".

### Pushbullet

Log into your Pushbullet account. Next, click on [Settings] on the left sidebar.
Scroll down to "Access Tokens" and click "Create Access Token". The text that
appears will be what you'll set `NOTI_PUSHBULLET_TOK` to.

### Pushover

Log into your [Pushover] account. Next, look for the "User Key". That's what
you'll set `NOTI_PUSHOVER_DEST` to.

Next [create a new application]. Fill in the fields. Under "Type", select
"Script". Finally, go to the application page. Look for "API Token/Key". This is
what you'll set `NOTI_PUSHOVER_TOK` to.

### Pushsafer

Log into your [Pushsafer] account. Next, look for the "Private or Alias Key".
That's what you'll set `NOTI_PUSHSAFER_KEY` to.

### Simplepush

Install the Simplepush Android app to get your Simplepush key.
That's the key you'll set to `NOTI_SIMPLEPUSH_KEY`.
Simplepush requires no registration and sending notifications is completely free.

In the app you can create events to customize ringtone and vibration patterns for
different kinds of notifications.
The event id you can set in the app, translates to `NOTI_SIMPLEPUSH_EVENT` in noti.

### Slack

Log into your Slack account. Then go to the [OAuth Tokens for Testing and
Development] page. Create a token. This is what you'll set `NOTI_SLACK_TOK` to.

The variable `NOTI_SLACK_DEST` can be set to a channel like `#general` or
`#random`. You can also set it to someone's username, like `@juan` or
`@variadico`.

### BearyChat

Log into your BearyChat account. Then create an [incoming robot][bc-incoming].
Next, look for the "Hook Address" (or "Hook 地址" in Chinese), this is what
you'll set `NOTI_BC_INCOMING_URI` to.


[Settings]: https://www.pushbullet.com/#settings
[Pushover]: https://pushover.net
[create a new application]: https://pushover.net/apps/build
[Pushsafer]: https://www.pushsafer.com
[OAuth Tokens for Testing and Development]: https://api.slack.com/docs/oauth-test-tokens
[bc-incoming]: https://bearychat.com/integrations/incoming
