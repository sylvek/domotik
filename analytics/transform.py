#!/usr/bin/python
import csv
import datetime
import argparse

parser = argparse.ArgumentParser(description='transform mongoexport to panda data')
parser.add_argument('data', metavar='data', help='data to export', nargs='?', default="sumPerDay")
parser.add_argument('csv', metavar='csv', help='csv input file', nargs='?', default="export.csv")
args = parser.parse_args()

def output(row):
    print '{0},{1},{2},{3},{4}'.format(round(row[2], 2), row[4], row[5], row[6], row[7])

with open(args.csv, 'rb') as csvfile:
    reader = csv.reader(csvfile)
    first_line = True
    for row in reader:
        if first_line:
            first_line = False
            print "value,hour,day,month,year"
        else:
            row[2] = float(row[2])
            row[3] = float(row[3])
            date = datetime.datetime.fromtimestamp(row[3])
            row.append(date.time().hour)
            row.append(date.isoweekday())
            row.append(date.month)
            row.append(date.year)
            row.append(date.day)

        if row[1] == "meanPerHour" == args.data and 50 < row[2] < 10000:
            output(row)

        if row[1] == "sumPerDay" == args.data and 4000 < row[2] < 100000:
            output(row)

        if row[1] == "livingRoomPerHour" == args.data and 0.009 < row[2] < 1.0:
            output(row)

        if row[1] == "tankHotWaterPerDay" == args.data and 0 < row[2] < 300:
            output(row)
