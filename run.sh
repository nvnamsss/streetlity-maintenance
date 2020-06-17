#/bin/bash
name=$1

echo "STOP RUNNING MAINTENANCE CONTAINER"
docker stop -t 30 ${name}_maintenance_container 
docker rm -f ${name}_maintenance_container

echo "DONE STOPPING"

docker run --name ${name}_maintenance_container -d\
            --network common-net \
            --restart always \
            -p 9000:9000 \
            maintenance_container

docker cp config.json ${name}_maintenance_container:/server/config/config.json    
