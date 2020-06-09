#!/bin/bash
# Run the docker container for pfh-reader against Local environment
# See also: dockerbuild.sh for building the docker image locally
docker run --restart always -d -v $HOME/.aws:/root/.aws -e app_environment=Local -e app_name=pfh_reader -e app_dir_path=/go/src/github.com/spacetimi/pfh_reader_server --publish 9000:9000 pfh-reader-server | xargs -I containerId docker logs -f containerId

