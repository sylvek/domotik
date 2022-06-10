# Domotik

![tv](extras/tv.png)

Domotik ~~is~~ was a very simple home-automation software based on micro-services principle.
~~Each service (services folder) does only one thing. A script (mosquitto_pub) pushes the data on MQTT.
A another process (mosquitto_sub) is in charge of dispatch the data somewhere.
An another service (services folder) could calculate the mean or the max value, etc and push it on MQTT, etc, etc.~~

Originally, I bought an [Envi CC128](http://www.currentcost.com/product-cc128.html), one Raspberry Pi 1 and used [measureit](https://github.com/lalelunet/measureit/wiki/The-end-of-the-measureit-project-%3F!) to display my power consumption. It was fun but reports were cheap and difficult to modify (think about a bunch of PHP and MySQL code). I think that i wrote the french translation but that's all i could do!

Early in 2015, i discovered [ESP8266](https://fr.wikipedia.org/wiki/ESP8266) chips, [Arduino UNO](https://en.wikipedia.org/wiki/Arduino_Uno) and IoT stuff (like MQTT). I decided to build my own devices…

I took a comm' module from an Oregon Scientific Clock in order to receive data from an external temperature sensor [thgr-511](https://www.disteo-sante.fr/accessoires/1759-thgr-511-sonde-thermo-hygro-.html) (and that worked!)

![clock](extras/oregonscientist.jpg)
![module](extras/communication_module.jpg)

After this first success, I decided to build my own temperature sensors (with or without battery but always with a certain art spirit..)

![sensor1](extras/temp_sensor_1.jpg)
![sensor2](extras/temp_sensor_2.jpg)

At last in 2018, i received a [Linky](https://fr.wikipedia.org/wiki/Linky) and decided to [extract data with an Raspberry Pi Zero](https://github.com/sylvek/linkiki).

![linky](extras/linky.jpg)

My main objective was to understand my power consumption in order to reduce bills but also to *save the Earth*. So i built a webpage to display a dashboard in my TV.

![tv2](extras/tv2.jpg)

This dashboard is displayed thanks to my original Raspberry Pi 1 (yes!) over an HDMI cable. ~~I built a [small web browser that fits well with an embedded device](github.com/sylvek/kiosk-browser/) (angularjs 1.4 works well on it, i did an unsuccesful test with reactjs)~~

Code side, i wrote several small scripts in Python using Mosquitto and Mongodb. Basically, everything ran on my raspberry pi. In 2016, i used Ansible to deploy it. In 2017, I decided to bought a real computer _(sorry for Raspberry pi fans)_. Indeed, raspberries are cool for thin devices but the IO stack and using an HDD disk _(to display movies for example)_ crashes too much my raspberry and my data. So I deciced to move on an [Intel NUC](https://en.wikipedia.org/wiki/Next_Unit_of_Computing) based on an x86 CPU (with 4GB of RAM) and an SSD. It's not so expansive comparing to a raspberry pi 2 or 3 with an SSD Disk on USB port. I decided to deploy my scripts using Docker and slowly, i rewrote the backend in Java. I also, decided to migrate my data from MongoDB to ~~InfluxDB _(much more suitable for timeseries data - i also tested Elastisearch but it took to much space, around 3 times more and too big for my needs)_~~ my dedicated datastore based on SQLite. So, i went to a micro-services architecture to a monolith. I learnt that begining by a micro-services architecture was cool to experiment some stuff, but using less code to maintain is better when you want to consolidate domain _(and reduce ops works)_.

In 2020, i replaced my raspberry-pi "TV" by a [Kindle TV Stick + Fully Kiosk](https://www.fully-kiosk.com/en/#download-box).

In 2022, i decided to give up influxdb and build a little datastore _(less footprint memory and more fun)_

## Build / Run it

```
> docker-compose up

> #open http://localhost:3000

> mosquitto_pub -h localhost -t sensors/esp8266/temp -m 4
> mosquitto_pub -h localhost -t sensors/linky/watt -m 4000
```

## Use it

Basically i use 3 temperature sensors and a "watt meter".

| name | topic | unit | usage |
|---------|------|------|-----|
| outside | sensors/esp12e/temp | xx.xx °c | basically, a value every 5 min |
| living room | sensors/esp8266/temp | xx.xx °c | a value every minute |
| rooms | sensors/esp32/temp | xx.xx °c | a value every minute |
| power consumption | sensors/linky/watt | xxxxx | a watt-hour value every second |

You can simulate data by using [`mosquitto_pub`](https://mosquitto.org/man/mosquitto_pub-1.html).

| service | link |
|---------|------|
| tv dashboard | http://your_ip:3000 |
| MQTT broker | tcp://your_ip:1883 |

I use Grafana _(plugged on SQLite)_ to display more dasboards. _(look extras/grafana-save folder)_

![today](extras/grafana_1.png)
![over 5y](extras/grafana_2.png)