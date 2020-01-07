require 'rubygems'
require 'mqtt'
require 'influxdb'

mqtt = MQTT::Client.new ARGV[0]
influxdb = InfluxDB::Client.new 'domotik', host: ARGV[1], async: true

begin
    mqtt.connect do |c|
        c.subscribe('sensors/#', 'measures/#')
        c.get do |topic, message|
            elements = topic.split('/')
            _type = elements[0]
            _name = elements[1]
            _unit = elements[2]

            data = {
                values: { value: message.to_f, type: _type, name: _name, unit: _unit },
                timestamp: Time.now.to_i
            }
            puts data
            influxdb.write_point(_type, data)
        end
    end
rescue Exception
    puts "ciao."
    mqtt.disconnect if mqtt
end