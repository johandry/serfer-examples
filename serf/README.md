# Serf cluster

This directory is used to start a Serf cluster without Surfer or any Serfer implementation, just to verify that what works for Serf also works for Serfer.

## Build image

    docker build --rm -t serf .

## Run Serf from container

Login to the container and start the agent:

    docker run --name serf --rm -it --entrypoint="/bin/sh" serf
    / # serf agent

## Run a Serf cluster with Docker Compose

It's required to build the image as explained above.

    docker-compose up

Use the flag `-d` to start the cluster in background and you can see the logs with:

    docker-compose logs -f leader 
    docker-compose logs -f node

Login to the leader:

    docker-compose exec leader /bin/sh

Or to a node (i.e. node #1):

    docker-compose exec --index=1 node /bin/sh

Tear down the cluster:

    docker-compose down
