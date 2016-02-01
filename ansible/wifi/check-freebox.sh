#!/bin/bash
 if ping -c 1 192.168.0.1 ; then
      echo "Network connection up"
 else
      echo "Network connection down! Attempting reconnection."
      sudo /sbin/ifup --force wlan0
 fi
