#!/usr/bin/python
import urllib2, urllib, json, sys
baseurl = "https://query.yahooapis.com/v1/public/yql?"
yql_query = "select item.condition,wind from weather.forecast where woeid=55863490 and u='c'"
yql_url = baseurl + urllib.urlencode({'q':yql_query}) + "&format=json"
result = urllib2.urlopen(yql_url).read()
data = json.loads(result)
temp = data['query']['results']['channel']['item']['condition']['temp']
wind = data['query']['results']['channel']['wind']['speed']
if len(sys.argv) == 2 and "wind" == sys.argv[1]: # wind or temp 
    print wind
else:
    print temp
