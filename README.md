# noti

![Testing](https://github.com/variadico/noti/actions/workflows/testing.yaml/badge.svg)

Monitor a process and trigger a notification.

Never sit and wait for some long-running process to finish. Noti can alert you
when it's done. You can receive messages on your computer or phone.

![macOS Banner Notification]

## Services

Noti can send notifications on a number of services.

| Service    | macOS | Linux | Windows |
| ---------- | :---: | :---: | :-----: |
| Banner     |   ✔   |   ✔   |    ✔    |
| Speech     |   ✔   |   ✔   |    ✔    |
| BearyChat  |   ✔   |   ✔   |    ✔    |
| Keybase    |   ✔   |   ✔   |    ✔    |
| Mattermost |   ✔   |   ✔   |    ✔    |
| Pushbullet |   ✔   |   ✔   |    ✔    |
| Pushover   |   ✔   |   ✔   |    ✔    |
| Pushsafer  |   ✔   |   ✔   |    ✔    |
| Simplepush |   ✔   |   ✔   |    ✔    |
| Slack      |   ✔   |   ✔   |    ✔    |
| Telegram   |   ✔   |   ✔   |    ✔    |
| Zulip      |   ✔   |   ✔   |    ✔    |
| Twilio     |   ✔   |   ✔   |    ✔    |
| GChat      |   ✔   |   ✔   |    ✔    |
| Chanify    |   ✔   |   ✔   |    ✔    |
| Bark       |   ✔   |   ✔   |    ✔    |

Check the [screenshots] directory to see what the notifications look like on different platforms.

## Installation

Install the Go binary with these commands.

```shell
# macOS install with Brew
brew install noti

# macOS install with curl
curl -L $(curl -s https://api.github.com/repos/variadico/noti/releases/latest | awk '/browser_download_url/ { print $2 }' | grep 'darwin-amd64' | sed 's/"//g') | tar -xz

# Linux install with curl
curl -L $(curl -s https://api.github.com/repos/variadico/noti/releases/latest | awk '/browser_download_url/ { print $2 }' | grep 'linux-amd64' | sed 's/"//g') | tar -xz
```

Or download it with your browser from the [latest release] page.

### From source

If you want to build from the source, then build like this.

```shell
# build binary
make build
# build binary and move to Go bin dir
make install
```

## Examples

Just put `noti` at the beginning or end of your regular commands. For more details, check the [docs].

Display a notification when `tar` finishes compressing files.

```sh
noti tar -cjf music.tar.bz2 Music/
```

Add `noti` after a command, in case you forgot at the beginning.

```sh
clang foo.c -Wall -lm -L/usr/X11R6/lib -lX11 -o bizz; noti
```

If you already started a command but forgot to use `noti`, then you can do this to get notified when that process' PID disappears.

```sh
noti --pwatch 1234
```

You can also press `ctrl+z` after you started a process. This will temporarily suspend the process, but you can resume it with `noti`.

```
$ dd if=/dev/zero of=foo bs=1M count=2000
^Z
zsh: suspended  dd if=/dev/zero of=foo bs=1M count=2000
$ fg; noti
[1]  + continued  dd if=/dev/zero of=foo bs=1M count=2000
2000+0 records in
2000+0 records out
2097152000 bytes (2.1 GB, 2.0 GiB) copied, 12 s, 175 MB/s
```

Additionally, `noti` can send a message piped from stdin with `-`.

```sh
$ make test 2>&1 | tail --lines 5 | noti -t "Test Results" -m -
```

[macos banner notification]: https://raw.githubusercontent.com/variadico/noti/main/docs/screenshots/macos_banner.png
[screenshots]: https://github.com/variadico/noti/tree/main/docs/screenshots
[latest release]: https://github.com/variadico/noti/releases/latest
[docs]: https://github.com/variadico/noti/blob/main/docs/noti.md
