#!/bin/bash
touch $HOME/domotik/services/fswebcam.conf
fswebcam -c $HOME/domotik/services/fswebcam.conf $1
