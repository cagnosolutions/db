#!/usr/bin/env bash

IP=$1
LOC=$(curl -s ipinfo.io/$IP | sed -r 's/[{}]|("[a-z]+":\ )|("|",)|//g') #| tac | sed '1,4d' | tac | sed '1,3d')
echo -e "\n  $IP\n$LOC\n"
