# Surfer a simple implementation of Serfer

This repo is an example about how to use Serfer package by building a simple binary named Surfer (_just because sounds cool and simmilar_) to create a sandbox cluster.

## Requirements

* **Docker**: To build and use Surfer you need, at least, Docker.
* **Docker Compose**: You may also need a cluster orchestrator, in this example, it is Docker Compose.
* **Go**: (_optional_) Go is needed only if you like to build surfer for your system.

The image can be packed with **Serf** just to verify the integration between Surfer and Serf, but Serf is not started when the container starts. Uncomment the lines in the Dockerfile to install Serf if needed. You may want to install it also in your computer if you plan to do the same with a local Surfer binary.

## Quick Start

In a terminal #1, run:

    docker build --rm -t surfer .
    docker-compose up -d
    docker-compose logs -f leader

You'll see a node running surfer and sending events (**5secEvent**) and queries (**7secQuery**). It's also the only node that will collect the responses from the other nodes (including itself).

In a terminal #2, execute:

    docker-compose logs -f node

You'll see another node running surfer.

Both nodes will receive the events and print the event payload. They also will receive the queries, print the payload and reply to the queries with a response. When they receive the **join** event, both will print all the nodes in the cluster.

To scale up/down the number of nodes the cluster, execute:

    docker-compose scale node=3

Replace `3` for the number of nodes desired. You'll see in the terminal #2 the logs of all the nodes with different color.

When you are done, tear down the cluster execute in the terminal #1:

    docker-compose down

## Build the image

Before start the cluster, build the Docker image `surfer` that contain the surfer binary.

    docker build --rm -t surfer .

This container is required by the docker-compose. List the images to verify it's there:

    docker images surfer

To login into the container execute:

    docker run -it --rm --entrypoint="/bin/sh" surfer

## Build the binary

To build the surfer binary for your system, execute:

    go build -o surfer .

## Testing Serfer

If you are doing changes in Serfer, modify the `replace` section in the `go.mod` file to use the local version of your Serfer code.

## Build the cluster orchestrated with Docker Compose

If you made a change in Serfer or Surfer, make sure to update the vendors and re-build the image as explained above.

Start the cluster with:

    docker-compose up

Use the flag `-d` to start the cluster in background, this option won't display the Surfer logs unless you use:

    docker-compose logs -f leader
    docker-compose logs -f node

To scale up or down the cluster, execute this command replacing `N` by the number of nodes to have:

    docker-compose scale node=N

Tear down the cluster with:

    docker-compose down

Login into the leader with:

    docker-compose exec --entrypoint="/bin/sh" leader

Replace `leader` by `--index=N node` to login into the node `N`.
