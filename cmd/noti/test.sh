#!/bin/bash
set -xuo pipefail

go install
noti -v

noti
sleep 3

noti ls -al
sleep 3

export NOTI_SOUND_FAIL="Funk"
noti _badcmd
sleep 3
unset NOTI_SOUND_FAIL

export NOTI_DEFAULT="banner speech"
noti -t "hello" -m "world"
sleep 3
unset NOTI_DEFAULT
