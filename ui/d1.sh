#!/usr/bin/env bash
../client/client -port 6892 &
sleep 0.5
./ui 6892
