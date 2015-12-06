# Contributing

Thanks for being interested in contributing! Here's how you can help!

## Create an issue

Before you create a new issue, check to see if someone else has opened a similar
issue. There might already be some discussion about the issue you were thinking
of opening.

Create an issue if you notice some unexpected behavior with Noti or if you think
the documentation (readme, code comments, etc) is lacking. New notification type
requests are also fine. Try to tag your issue appropriately. If the issue is a
bug report, use the bug tag; if it's a feature request use the enhancement tag.

## Submit a pull request

Wow. Thanks! You're awesome!

But before spending hours writing some code, make sure you voice your proposal
with an issue. After you get the green light there, then you can start coding
away!

### Set up dev environment

The following steps are for hackers only. If you're a regular user, then follow
the [installation instructions on the readme][1].

#### Fork Noti

Create a fork of Noti on GitHub.

#### Install Go

Since this project is small, the version of Go you install hopefully won't
matter too much. However, Noti has been developed on Go 1.5.2.

On OS X, you can install Go with this command.

```
brew install go
```

Make sure you installed Go correctly and don't forget to set your `GOPATH`.

#### Install dependencies

Noti uses cgo for most notification types. Packages that use cgo have a build
tag to protect against compiling on the wrong operating system. In addition,
you'll need to have the right C dependencies installed to compile Noti.


##### OS X

The following notification types are specific to OS X.

* package nsspeechsynthesizer
* package nsuser

Lucky for you, Cocoa comes installed with OS X. So, you won't have to install
anything  extra.

##### Linux and FreeBSD

The following notification types are specific to Linux and FreeBSD.

* package espeak
* package libnotify

On Ubuntu 14.04, you can run this command to install the required dependencies.

```
sudo apt-get install libespeak-dev
sudo apt-get install libnotify-dev
```

#### Download source

Clone the repo into your Go path.

```
cd $GOPATH/src/github.com/your_username
git clone git@github.com:your_username/noti.git
cd noti
```

### Hack!

Project Noti has two parts. First, Noti is a library of independent notification
types. Second, `noti` is a command that uses the Noti library. When I write
"Noti", I'm talking about the Noti library and packages like `nsuser` or
`espeak`. When I write `noti`, I'm talking about the code at `noti/cmd/noti/`.

If you want to add a new notification type to `noti`, then you'll first have to
create a new Noti package. Then, you can import your package in `noti` and
integrate it there.

### Contribute

Before submitting a pull request, make sure you run these tools on your code.
Clean up any warnings or errors these tools generate.

* `gofmt`
* `golint`
* `go vet`
* `gofmt -s *.go`

After that, submit a pull request on GitHub! Thanks for helping make Noti
better!

## Reading List

* [NSUserNotification Class Reference][2]
* [Libnotify Reference Manual][3]
* [Pushbullet API][4]
* [espeak.cpp][5]

[1]: https://github.com/variadico/noti/blob/master/README.md#installation
[2]: https://developer.apple.com/library/mac/documentation/Foundation/Reference/NSUserNotification_Class/#//apple_ref/doc/constant_group/NSUserNotificationDefaultSoundName
[3]: https://developer.gnome.org/libnotify/0.7/
[4]: https://docs.pushbullet.com/
[5]: https://fossies.org/dox/espeak-1.48.04-source/espeak_8cpp_source.html
