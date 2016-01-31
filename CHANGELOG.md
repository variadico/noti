# Change Log

All notable changes to this project will be documented in this file. This
project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

## Added

* Speech notifications.
* Slack notifications.
* Optionally set default notification type(s) through `NOTI_DEFAULT` env var.
* Multi-notification support.
* Other configuration through environment variables.
* Contributing document.
* Change log.

## Changed

* On OS X, the notification sound must now be set in the environment variable,
`NOTI_SOUND`.
* On OS X, instead of AppleScript banner notifications are triggered with
Object-C, which shows (nicer) Terminal icon.

## Removed

* `-f` flag for OS X. This caused unexpected behavior for people who use iTerm2.
* OS X-specific flags and usage text from Linux and FreeBSD help.

[Unreleased]: https://github.com/variadico/noti/compare/master...v2
