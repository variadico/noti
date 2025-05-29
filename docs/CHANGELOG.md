# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased] - YYYY-MM-DD

* Add access token support for ntfy

## [3.8.0] - 2025-03-06

### Added

* Chanify support
* ntfy support
* Bark support

### Removed

* Deprecated 'as_user' parameter from Slack service

## [3.7.0] - 2023-02-12

### Added

* Add Telegram Topic support

## [3.6.0] - 2022-10-10

### Added

* Support for Google Chat through webhooks

## [3.5.0] - 2020-09-06

### Added

* `--time` option to show execution time in message
* Twilio support

### Fixed

* Notification title in Telegram

## [3.4.0] - 2019-12-07

### Added

* Support to read message body from stdin
* Support for keybase

### Removed

* Support for HipChat

### Changed

* Noti's exit status is now based on executed subcommand's status
* Shell login flag is no longer passed when executing an alias

## [3.3.0] - 2019-09-11

### Added

* Support for Zulip.
* Support for Telegram.

### Fixed

* The name "noti" was appearing twice in Slack.

### Deprecated

* Hipchat will be removed in next release.

## [3.2.0] - 2019-06-25

### Added

* Support for Mattermost.
* Support for Slack App URL.
* Support for pollPID for BSD and Windows.

### Removed

* Installation by `go get`, can revisit in Go 1.14.

### Fixed

* Incorrect deprecation mapping.

## [3.1.0] - 2018-03-25

### Added

* Pushbullet config key to specify a target device.

### Fixed

* `--message` let's you change the message again.

## [3.0.0] - 2018-01-06

### Added

* `.noti.yaml` config file supports current notifications.
* Speech notification support for Windows.
* Man pages for `noti` and `noti.yaml`.
* Configuration option for Slack bot name.

### Changed

* `-b=false` no longer has to be passed to disable banner when enabling
  multiple notifications.
* `--help` has been simplified.
* `--pwatch` now polls PID every 2 seconds instead of every 1 second.
* Version omits `v` prefix, e.g. `3.0.0` instead of `v3.0.0`.

### Deprecated

* In an effort to normalize and allow future configuration for different
  services, certain environment variables have been deprecated.

```
| Deprecated           | Current                        |
---------------------------------------------------------
| NOTI_SOUND           | NOTI_NSUSER_SOUNDNAME          |
| NOTI_SOUND_FAIL      | NOTI_NSUSER_SOUNDNAMEFAIL      |
| NOTI_VOICE           | NOTI_SAY_VOICE                 |
| NOTI_VOICE           | NOTI_ESPEAK_VOICENAME          |
| NOTI_VOICE           | NOTI_SPEECHSYNTHESIZER_VOICE   |
| NOTI_BC_INCOMING_URI | NOTI_BEARYCHAT_INCOMINGHOOKURI |
| NOTI_HIPCHAT_TOK     | NOTI_HIPCHAT_ACCESSTOKEN       |
| NOTI_HIPCHAT_DEST    | NOTI_HIPCHAT_ROOM              |
| NOTI_PUSHBULLET_TOK  | NOTI_PUSHBULLET_ACCESSTOKEN    |
| NOTI_PUSHOVER_TOK    | NOTI_PUSHOVER_TOKEN            |
| NOTI_PUSHOVER_DEST   | NOTI_PUSHOVER_USER             |
| NOTI_SLACK_TOK       | NOTI_SLACK_TOKEN               |
| NOTI_SLACK_DEST      | NOTI_SLACK_CHANNEL             |
```

### Removed

* Single-dash long options. Long flags must be passed with two dashes, e.g.
  `--version`.

## [2.7.0] - 2017-10-09

### Deprecated

* Single-dash long options. Any scripts using `-banner` or `-title` should be
  updated to use `--banner` or `--title` instead.

## [2.6.0] - 2017-07-08

### Added

* Support for Pushsafer.

## [2.5.0] - 2016-11-25

### Added

* Support for Simplepush.

### Fixed

* Formatting bug in help.

## [2.4.0] - 2016-10-01

### Added

* Support for BearyChat.

## [2.3.0] - 2016-09-09

### Added

* Banner support for Windows 10.

### Fixed

* Aliases in Bash and ZSH now work.

## [2.2.2] - 2016-08-28

### Fixed

* `noti` now compiles on Windows.

## [2.2.1] - 2016-08-27

### Added

* Add notification config docs to repo.

### Fixed

* Bug that caused noti to not work on macOS Sierra.

## [2.2.0] - 2016-05-26

### Added

* `-pwatch` flag to trigger notification after PID disappears.
* Check for updates during `-v` flag.

### Changed

* Install command changed to: `go get -u github.com/variadico/noti/cmd/noti`

## [2.1.1] - 2016-02-27

### Added

* Tests.

### Fixed

* Setting `-t` or `-m` will now take precedence over utility name.

### Changed

* Slackbot icon is now a rocket.
* If utility fails, noti will exit 1.
* Utility name default now includes subcommand too.
* Improved error handling for notifiers that shell out.

## [2.1.0] - 2016-02-23

### Added

* `NOTI_SOUND_FAIL` tells Noti which sound to play for banner notifications
  when a utility fails on OS X.
* `NOTI_SOUND` and `NOTI_SOUND_FAIL` can be set to `_mute` for silent
  notifications.
* HipChat notifications.
* Pushover notifications.
* Noti Wiki.

### Changed

* On OS X, banner notifications will play different sounds depending on the
  utility's success or failure, instead of the same sound for both.
* `NOTI_SLACK_DEST` no longer defaults to "#random". It must be manually set.

## [2.0.0] - 2016-01-31

### Added

* Speech notifications.
* Slack notifications.
* Optionally set default notification type(s) through `NOTI_DEFAULT` env var.
* Multi-notification support.
* Other configuration through environment variables.
* Contributing document.
* Change log.

### Changed

* On OS X, the notification sound must now be set in the environment variable,
  `NOTI_SOUND`.
* On OS X, instead of AppleScript, banner notifications are triggered with
  Object-C, which shows (nicer) Terminal icon.

### Removed

* `-f` flag for OS X. This caused unexpected behavior for people who use iTerm2.
* OS X-specific flags and usage text from Linux and FreeBSD help.


[Unreleased]: https://github.com/variadico/noti/compare/3.8.0...main
[3.8.0]: https://github.com/variadico/noti/compare/3.7.0...3.8.0
[3.7.0]: https://github.com/variadico/noti/compare/3.6.0...3.7.0
[3.6.0]: https://github.com/variadico/noti/compare/3.5.0...3.6.0
[3.5.0]: https://github.com/variadico/noti/compare/3.4.0...3.5.0
[3.4.0]: https://github.com/variadico/noti/compare/3.3.0...3.4.0
[3.3.0]: https://github.com/variadico/noti/compare/3.2.0...3.3.0
[3.2.0]: https://github.com/variadico/noti/compare/3.1.0...3.2.0
[3.1.0]: https://github.com/variadico/noti/compare/3.0.0...3.1.0
[3.0.0]: https://github.com/variadico/noti/compare/v2.7.0...3.0.0
[2.7.0]: https://github.com/variadico/noti/compare/v2.6.0...v2.7.0
[2.6.0]: https://github.com/variadico/noti/compare/v2.5.0...v2.6.0
[2.5.0]: https://github.com/variadico/noti/compare/v2.4.0...v2.5.0
[2.4.0]: https://github.com/variadico/noti/compare/v2.3.0...v2.4.0
[2.3.0]: https://github.com/variadico/noti/compare/v2.2.2...v2.3.0
[2.2.2]: https://github.com/variadico/noti/compare/v2.2.1...v2.2.2
[2.2.1]: https://github.com/variadico/noti/compare/v2.2.0...v2.2.1
[2.2.0]: https://github.com/variadico/noti/compare/v2.1.1...v2.2.0
[2.1.1]: https://github.com/variadico/noti/compare/v2.1.0...v2.1.1
[2.1.0]: https://github.com/variadico/noti/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/variadico/noti/compare/v1.3.0...v2.0.0
