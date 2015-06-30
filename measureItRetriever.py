#!/usr/bin/python
import urllib, json, sys
url = "http://192.168.0.2/measureit/measureit_public_html/php/measureit_functions.php?do=summary_start"
response = urllib.urlopen(url);
data = json.loads(response.read())
print data[0][sys.argv[1]] #tmpr or watt

