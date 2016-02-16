#!/bin/bash
 if ping -c 1 192.168.0.1 ; then
      echo "Network connection up"
 else
      ps cax | grep ifup > /dev/null
      if [ $? -eq 1 ]; then
        echo "Network connection down! Attempting reconnection."
        /sbin/ifup --force wlan0
      else
        echo "wlan0 is already waiting for up."
      fi
 fi
