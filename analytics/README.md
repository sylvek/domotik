# data analysis

This page is a collection of reflection about data analysis. It is an introduction
of a new feature of this sofware.

Currently, services detection are based on an analysis of the received data between
yesterday an now. If this value is greater than x percent, there is a problem.

My two years of data collection are a funny sandbox which i'm going to use to
experiment some features like :

- detect abnormal power consumption
- detect abnormal human intrusion
- detect the best moment to turn on my heaters
- detect a heat wave
- etc.

first of all, how to extract data from mongodb?
```
$> mongoexport --host localhost --db domotik --collection measures --csv --out export.csv --fields type,sensor,value,timestamp
```

the second step is focused on transforming data to generate several files
```
$> python transform.py tankHotWaterPerDay > tankHotWaterPerDay.csv
```

using panda via docker is a good way to save your time!
```
$> go to https://github.com/ContinuumIO/docker-images/tree/master/anaconda
$> docker pull continuumio/anaconda
$> docker run -i -t -p 8888:8888 continuumio/anaconda /bin/bash -c "/opt/conda/bin/conda install jupyter -y --quiet && mkdir /opt/notebooks && /opt/conda/bin/jupyter notebook --notebook-dir=/opt/notebooks --ip='*' --port=8888 --no-browser"
$> go to http://localhost:8888
```

enter something like that
```
import pandas as pd
df = pd.read_csv('sumPerDay.csv')
min = df.groupby(['month','day']).mean() - df.groupby(['month','day']).std()
max = df.groupby(['month','day']).mean() + df.groupby(['month','day']).std()
print min.value
print max.value
min.to_csv('sumPerDay.min')
max.to_csv('sumPerDay.max')
```

retrieve this files from jupyter and store it into the folder 'services/data'

at this point, we could use min and max to check if a new value is "correct"

# image analysis

first of all, we need to extract images data. Images are stored on mongodb but
every 24h this database is erased (obviously data are backuped on an hard drive disk)

to retrieve data, the simpliest way is to "zgrep" backup files.

```
$ zgrep -h camera sensors* > ~/Downloads/camera.json
```

data is a list of json object containing a value field. This field is a base64
jpg encoded.

To extract data to jpg files, we can use :

```
# camera.json should be in the same folder of extract_images.py
$ python extract_images.py
```
