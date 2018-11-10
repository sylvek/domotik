require 'mqtt'

MQTT_HOSTNAME     = ARGV[0] || "localhost"

class Co2Reporter
  GREEN     = "#009900"
  YELLOW    = "#ffa500"
  RED       = "#ff0000"

  attr_accessor :co2, :watt

  def initialize host
    MQTT::Client.connect(host) do |c|
      c.get('sensors/co2/gperkwh') do |topic,message|
        perform c, message.to_i
      end
    end
  end

  def perform client, message
    # puts "current co2: #{@co2}, mean watt/min: #{@watt}"
    # color depends on co2 indice
    # lower than 74gCO2/kwh means that CO2 production is clean
    # around 74gCO2/kwh is the annual production mean in France
    # greater than 85gCO2/kwh should alert us!
    case message || 74
    when 0..70 then color_is = GREEN
    when 71..85 then color_is = YELLOW
    else color_is = RED
    end

    client.publish 'triggers/led/update', color_is
  end
end

trap "SIGINT" do
  puts "Exiting"
  exit 130
end

Co2Reporter.new ARGV[0]
