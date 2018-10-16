curl -s 'https://api.co2signal.com/v1/latest?countryCode=FR' -H 'auth-token: e12e7866c5c8a779' | jq .data.carbonIntensity | mosquitto_pub -h $1 -t $2 -s
