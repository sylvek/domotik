import json
import base64
import sys

def process(line):
    sensor = json.loads(line)
    file_name = sensor['timestamp']
    image_base64 = sensor['value']
    image_raw = base64.standard_b64decode(image_base64)
    f = open(str(file_name) + '.jpg', 'w')
    f.write(image_raw)
    f.close()

with open('camera.json') as f:
    for line in f:
        sys.stdout.write('.')
        process(line)
        sys.stdout.flush()
