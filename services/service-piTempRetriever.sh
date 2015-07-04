#!/bin/sh
/opt/vc/bin/vcgencmd measure_temp | grep -o "[0-9.]*"
