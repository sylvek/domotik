require 'mqtt'
require 'json'

MQTT::Client.connect(ENV["MQTT_HOST"]) do |c|
  c.get('tele/+/SENSOR') do |topic,message|
    device = topic.split('/')[1]
    if JSON.parse(message)["DS18B20"]
        c.publish("sensors/#{device}/temp", "#{JSON.parse(message)["DS18B20"]["Temperature"]}", retain=true)
    end
    if JSON.parse(message)["AM2301"]
        c.publish("sensors/#{device}/temp", "#{JSON.parse(message)["AM2301"]["Temperature"]}", retain=true)
        c.publish("sensors/#{device}/humidity", "#{JSON.parse(message)["AM2301"]["Humidity"]}", retain=true)
        c.publish("sensors/#{device}/dewpoint", "#{JSON.parse(message)["AM2301"]["DewPoint"]}", retain=true)
    end
  end
end