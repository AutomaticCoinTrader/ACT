#!/bin/bash
rm log
while true
do
    ./ACT -confdir ./config >> log 2>&1
done