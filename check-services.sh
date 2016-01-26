#!/bin/bash
ps aux | grep python | grep "\.py" | awk '{print $2,$13}'
