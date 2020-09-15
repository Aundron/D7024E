#!/bin/bash

go build -o ./bin/kademlia
docker build . -t kademlia