# Change Log

All notable changes to this project will be documented in this file. This
project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

## Changed

* Noti has been refactored into a library and command.

## Added

* Package espeak for speech notifications on Linux and FreeBSD.
* Package libnotify for banner notifications on Linux and FreeBSD.
* Package nsspeechsynthesizer for speech notifications on OS X.
* Package nsuser for banner notifications on OS X.
* Package pushbullet for multi-device notifications.
* Contributing document
* Change log

## Removed

* `-f` flag for OS X. This caused unexpected behavior for people who use iTerm2.
* OS X-specific flags and usage text from Linux and FreeBSD.

[Unreleased]: https://github.com/variadico/noti/compare/master...v2
