require 'sqlite3'
require 'httparty'

entrypoint = ENV['ENTRYPOINT']
db = SQLite3::Database.new ENV['DATABASE']
name = ENV['NAME'] || 'daily_power_consumption'

db.execute("SELECT ts,value,unit FROM data WHERE name='#{name}' and ts=strftime('%s','now', 'start of day', '-1 day') ORDER BY ts DESC LIMIT 1") do |row|
    
    timestamp = row[0]
    value = row[1]
    unit = row[2]
    
    title = "#{value.to_i/1000.0} kW.h" if unit == "indice"
    title = "#{sprintf('%.2f', value)} C" if unit == "temp"
    title = "#{value} L" if unit == "liter"
    
    response = HTTParty.get("#{entrypoint}?title=#{title}&timestamp=#{timestamp}")
    
    puts "#{Time.at(timestamp).utc.to_datetime} - #{response.body}"
end
