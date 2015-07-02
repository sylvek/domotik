#!/usr/bin/python
import urllib2, urllib, json, argparse

parser = argparse.ArgumentParser(description='retrieve weather yahoo data from woeid')
parser.add_argument('woeid', metavar='woeid', help='woeid of your city')
parser.add_argument('sensor', metavar='sensor', help='wind or temp')
args = parser.parse_args()

baseurl = "https://query.yahooapis.com/v1/public/yql?"
yql_query = "select item.condition,wind from weather.forecast where woeid=" + args.woeid + " and u='c'"
yql_url = baseurl + urllib.urlencode({'q':yql_query}) + "&format=json"
result = urllib2.urlopen(yql_url).read()
data = json.loads(result)
temp = data['query']['results']['channel']['item']['condition']['temp']
wind = data['query']['results']['channel']['wind']['speed']
if "wind" == args.sensor: # wind or temp
    print wind
else:
    print temp
