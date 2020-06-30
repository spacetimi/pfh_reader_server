#!/bin/bash
# Run the docker container for pfh-reader against Local environment
# See also: dockerbuild.sh for building the docker image locally
docker run --restart always -d -e "TZ=Asia/Kolkata" -e app_environment=Local -v $HOME/Library/Containers/com.spacetimi.pfh-daemon/Data/Documents:/root/Library/Containers/com.spacetimi.pfh-daemon/Data/Documents --publish 9001:9001 pfh-reader-server | xargs -I containerId docker logs -f containerId
