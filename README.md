# noti
Display a notification in OS X after a program finishes running.

## Installation
If you have Go installed, then you can do this.
````
go get github.com/pi241a/noti
````

## Usage
Basically, just prefix your regular commands with `noti`.
````
noti [-tms] <utility> [args...]

    -t    Title of notification. Default is the utility name.

    -m    Message notification will display. Default is "Done!"

    -s    Sound to play when notified. Default is Ping. Possible options
          are Basso, Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping,
          Pop, Purr, Sosumi, Submarine, Tink. Check /System/Library/Sounds
          for more info.

    -h    Display usage information and exit.
````

## Examples
Get notified when `curl` finishes downloading files.
````
noti curl -O https://wordpress.org/latest.tar.gz
````

Get notified when `tar` finishes compressing files.
````
noti tar -cjf music.tar.bz2 Music/
````

Get notified when `brew` finishes updating. Set the notification title to
"Homebrew" and the message to "Up to date."
````
noti -t "Homebrew" -m "Up to date." brew update
````
