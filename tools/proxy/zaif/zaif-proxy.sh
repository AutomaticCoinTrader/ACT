#!/bin/bash
rm log
while true
do
    ./zaif-proxy -confdir ./config >> log 2>&1
done