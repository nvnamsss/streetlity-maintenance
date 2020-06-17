#/bin/bash

docker image rm -f maintenance_container
docker build . -t maintenance_container

