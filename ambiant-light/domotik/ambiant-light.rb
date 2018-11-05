require 'yeelight'
require 'mqtt'
require 'json'
require 'rx'

MQTT_HOSTNAME     = ARGV[0] || "localhost"
YEELIGHT_HOSTNAME = ARGV[1] || "localhost"
WATT_MEAN_RATE    = ARGV[2] || "60"

class Array
  def average
      (size > 0) ? inject(&:+) / size : 0
  end
end

class MqttObservable
  def initialize(mqtt_hostname, mqtt_topics, subject)
    @mqtt_hostname = mqtt_hostname
    @mqtt_topics = mqtt_topics
    @subject = subject
  end
  def start
    @thr = Thread.new do
      @client = MQTT::Client.connect(@mqtt_hostname) do |c|
        c.subscribe(@mqtt_topics)
        c.get do |topic, message|
          @subject.on_next(message)
        end
      end
    end
  end
end

class YeelightAmbiant

  GREEN     = "9bf442"
  YELLOW    = "f4ee41"
  RED       = "f24126"

  LOW       = 1
  MEDIUM    = 3
  HIGH      = 8
  TOO_HIGH  = 20

  attr_accessor :co2, :watt

  def initialize host
    @yeelight = Yeelight::Client.new(host, 55443)
    @lamp = Yeelight::Lamp.new(@yeelight)
  end
  def on
    is_on = JSON.parse(yeelight.get_prop('"power"'))['data']['result'][0] == 'on'
    @yeelight.on unless is_on
  end

  def perform
    # puts "current co2: #{@co2}, mean watt/min: #{@watt}"
    # color depends on co2 indice
    # lower than 74gCO2/kwh means that CO2 production is clean
    # around 74gCO2/kwh is the annual production mean in France
    # greater than 85gCO2/kwh should alert us!
    case @co2 || 74
    when 0..70 then color_is = GREEN
    when 71..85 then color_is = YELLOW
    else color_is = RED
    end

    # animation speed depends on watt consumption
    # lower than 1000wh is the base
    # around 2000wh means that a device like hot tank or washing machine is on
    # around 3000wh means that hover or electric hater is working
    # on top of 5000wh, please do something!
    case @watt || 0
    when 0..1000 then speed_is = LOW
    when 1001..2000 then speed_is = MEDIUM
    when 2001..5000 then speed_is = HIGH
    else
      speed_is = TOO_HIGH
      @yeelight.on
    end

    # puts "color: #{color_is}, speed: #{speed_is}"
    @lamp.toggle_color(color_is, speed_is)

  end
end

yeelight_ambiant = YeelightAmbiant.new YEELIGHT_HOSTNAME
linky_subject = Rx::Subject.new
co2_subject = Rx::Subject.new
linky_subject
  .as_observable
  .buffer_with_time(WATT_MEAN_RATE.to_i)
  .subscribe do |x|
    yeelight_ambiant.watt = x.map{|a|a.to_i}.average
    yeelight_ambiant.perform
  end
co2_subject
  .as_observable
  .subscribe do |x|
    yeelight_ambiant.co2 = x.to_i
    yeelight_ambiant.perform
  end
linky = MqttObservable.new(MQTT_HOSTNAME, ["sensors/linky/watt"], linky_subject).start
co2 = MqttObservable.new(MQTT_HOSTNAME, ["sensors/co2/gperkwh"], co2_subject).start

trap "SIGINT" do
  #thr.exit
  linky.exit
  co2.exit
  puts "Exiting"
  exit 130
end

linky.join
co2.join
