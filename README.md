# Kademlia - Distributed Hash Table
This repository contains an implementation of Kademlia, a distributed hash table. Built in the course D7024E.

## Prerequisites
To be able to run this you need:
* Docker (https://www.docker.com)
* Docker-compose (https://docs.docker.com/compose/install/)
* Git-bash or similar shell if using Windows

## Build

First, we need to build the docker image. The easiest way is the run the `build.sh` shell script:

```console
$ ./build.sh
```

## Run

Now, to spin up the network of nodes, run the `start.sh` script. This will by default spin up a network of 50 nodes.

```console
$ ./start.sh
```

If you want to change the amount of nodes that are created, change the `--scale` flag in the `start.sh` script.

```shell
...

docker-compose up -d --scale kademliaNodes=<Amount>
```

To stop and tear down the network of nodes, just run the `stop.sh` script:

```console
$ ./stop.sh
```