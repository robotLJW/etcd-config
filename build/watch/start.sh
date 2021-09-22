#!/bin/sh
set -e

cd /home

for i in $(seq 1 2)
do
    ./watch &
done

while true; do
    sleep 60
done
