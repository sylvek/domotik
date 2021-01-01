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

This dashboard is displayed thanks to my original Raspberry Pi 1 (yes!) over an HDMI cable. I built a [small web browser that fits well with an embedded device](github.com/sylvek/kiosk-browser/) (angularjs 1.4 works well on it, i did an unsuccesful test with reactjs)

Code side, i wrote several small scripts in Python using Mosquitto and Mongodb. Basically, everything ran on my raspberry pi. In ~2016, i used Ansible to deploy it. In 2017, I decided to bought a real computer (sorry for Raspberry pi fans). Indeed, raspberries are cool for thin devices but the IO stack and using an HDD disk (to display movies for example) crashes too much my raspberry and my data. So I deciced to move on an [Intel NUC](https://en.wikipedia.org/wiki/Next_Unit_of_Computing) based on an x86 CPU (with 4GB of RAM) and an SSD. It's not so expansive comparing to a raspberry pi 2 or 3 with an SSD Disk on USB port.  I decided to depoy my script with Docker and slowly, i rewrote the backend in Java (VertX). I decided to migrate my data from Mongo to InfluxDB (after tested during 1y Elasticsearch). So, i went to a micro-services architecture to a monolith. I learnt that begining with a micro-services architecture was cool to experiment some stuff, but using a VertX code with only one project to maintain is better when you want to consolidate domain _(and reduce ops works)_.

## Build it

```
back> docker build -t domotik-back .
front> docker build -t domotik-front .
bridge-to-infludb> docker build -t domotik-bridge-to-influxdb .
```

## Run it

```

$> docker run -d --name influxdb -p 8086:8086 influxdb
$> docker run -d --name mosquitto -p 1883:1883 -p 9883:9883 jllopis/mosquitto:v1.4.14 mosquitto

$> docker run -d --name domotik-back --link mosquitto:mosquitto domotik-back app.jar mosquitto
$> docker run -d --name domotik-front --link influxdb:influxdb -p 3000:3000 domotik-front

$> docker run -d --name domotik-bridge-to-influxdb --link mosquitto:mosquitto --link influxdb:influxdb domotik-bridge-to-influxdb
```

## Use it

| service | link |
|---------|------|
| tv dashboard | http://your_ip:3000 |
| MQTT broker | tcp://your_ip:1883 |
| influxdb | http://your_ip:8086/domotik |

I use Grafana to display more dasboards. It compatibles with influxdb ;)

*Note, I reduced by 23% my consumption!*

![today](extras/grafana_1.png)
![over 5y](extras/grafana_2.png)

## influxdb tips

```
> show continuous queries;
name: _internal
name query
---- -----

name: domotik
name      query
----      -----
outside   CREATE CONTINUOUS QUERY outside ON domotik BEGIN SELECT mean(value) AS value INTO domotik.infinite.daily_temp_outside FROM domotik.autogen.sensors WHERE "name" = 'esp12e' GROUP BY time(1d) END
inside    CREATE CONTINUOUS QUERY inside ON domotik BEGIN SELECT mean(value) AS value INTO domotik.infinite.daily_temp_inside FROM domotik.autogen.sensors WHERE ("name" = 'esp8266' OR "name" = 'esp32') GROUP BY time(1d) END
sumPerDay CREATE CONTINUOUS QUERY sumPerDay ON domotik BEGIN SELECT sum(value) AS value INTO domotik.infinite.daily_power_consumption FROM domotik.autogen.measures WHERE "name" = 'meanPerHour' GROUP BY time(1d) END
```

```
> show retention policies;
name     duration  shardGroupDuration replicaN default
----     --------  ------------------ -------- -------
autogen  2160h0m0s 168h0m0s           1        true
infinite 0s        168h0m0s           1        false
```
