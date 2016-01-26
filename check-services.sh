#!/bin/bash
ps aux | grep python | grep "\.py" | awk '{print $2,$12,$13}'
