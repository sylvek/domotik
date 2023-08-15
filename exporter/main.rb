require 'sqlite3'
require 'httparty'

entrypoint = ENV['ENTRYPOINT']
db = SQLite3::Database.new ENV['DATABASE']

db.execute("SELECT ts,value FROM data WHERE name='daily_power_consumption' ORDER BY ts DESC LIMIT 1") do |row| 
    title = "#{row[1].to_i/1000.0} kW.h"
    response = HTTParty.get("#{entrypoint}?title=#{title}&timestamp=#{row[0]}")
    puts response.body
end