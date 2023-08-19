require 'sqlite3'
require 'httparty'

entrypoint = ENV['ENTRYPOINT']
db = SQLite3::Database.new ENV['DATABASE']
name = ENV['NAME'] || 'daily_power_consumption'

db.execute("SELECT ts,value,unit FROM data WHERE name='#{name}' and ts=strftime('%s','now', 'start of day', '-1 day') ORDER BY ts DESC LIMIT 1") do |row|
    unit = row[2]
    title = "#{row[1].to_i/1000.0} kW.h" if unit == "indice"
    title = "#{sprintf('%.2f', row[1])} C" if unit == "temp"
    title = "#{row[1]} L" if unit == "liter"
    response = HTTParty.get("#{entrypoint}?title=#{title}&timestamp=#{row[0]}")
    puts response.body
end