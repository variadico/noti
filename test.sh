#!/bin/bash
set -euo pipefail
set -x

TITLE="test title"
MESG="test message"

noti
sleep 3

noti -t "$TITLE"
sleep 3

noti -m "$MESG"
sleep 3

noti -t "$TITLE" -m "$MESG" ls
sleep 3

noti -b=0 -speech ls
sleep 3

export NOTI_DEFAULT="slack"
noti ls
sleep 3

unset NOTI_DEFAULT
export NOTI_SOUND="Hero"
noti ls
sleep 3

export NOTI_SOUND_FAIL="Funk"
noti _badcmd
sleep 3
