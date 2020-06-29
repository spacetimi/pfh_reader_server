#!/bin/bash
# Run the docker container for pfh-reader against Local environment
# See also: dockerbuild.sh for building the docker image locally
docker run --restart always -d -v $HOME/Library/Containers/com.spacetimi.pfh-daemon/Data/Documents:/root/Library/Containers/com.spacetimi.pfh-daemon/Data/Documents -e app_environment=Local -e app_name=pfh_reader -e app_dir_path=/go/src/github.com/spacetimi/pfh_reader_server -e shared_dir_path=/go/src/github.com/spacetimi/timi_shared_server --publish 9001:9001 pfh-reader-server | xargs -I containerId docker logs -f containerId

