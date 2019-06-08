#!/usr/bin/env bash
rm robot.tar.gz robot-docker.zip
cp ../robot .
mkdir data
echo "VOLUME " >data/readme.txt
tar cvfz robot.tar.gz robot data config.json
zip -r9 robot-docker.zip robot.tar.gz Dockerfile
#docker build -t powerchatrobot .

