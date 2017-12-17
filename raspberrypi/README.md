# Raspberry client

## Requirements

- last rapsbian ISO
- at least 512 Mo RAM _(but should work with Raspberry pi 1A - 256Mo)_

## Install browser

```
$> sudo apt-get update
$> sudo apt-get install x11-xserver-utils libwebkitgtk-dev xserver-xorg xinit matchbox-window-manager
```

```
$> sudo apt-get update
$> sudo apt-get install git
$> git clone https://github.com/sylvek/kiosk-browser.git
kiosk-browser> make
```

```
change on `/etc/X11/Xwrapper.config` for
allowed_users=anybody
```

```
add before `exit 0` on /etc/rc.local
`startx -- -nolisten tcp -nocursor -dpms -s 0 &`
```

```
$> DISPLAY=:0 /home/pi/kiosk-browser/browser http://192.168.0.4:3000/internal
```

## Install temperature sensor

Ensure to activate Wire-1 with raspi-config.

```
$> git clone https://github.com/sylvek/domotik.git
$> sudo apt-get install python-pip
$> sudo pip install paho-mqtt
/home/pi/domotik/raspberrypi/domotik/mosquitto_pub> python temperature.py '/sys/bus/w1/devices/28-0315820a56ff/w1_slave' 192.168.0.4
```


## Install led sensor

```
$> git clone https://github.com/sylvek/domotik.git
$> sudo apt-get install python-pip
$> sudo pip install paho-mqtt RPi.GPIO
/home/pi/domotik/raspberrypi/domotik/mosquitto_sub> python led.py 192.168.0.4
```
