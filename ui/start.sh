#!/usr/bin/env bash
../client/client -port 6990 > log.txt &
sleep 0.5
./ui 6990
