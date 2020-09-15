#!/bin/bash

./stop.sh
docker-compose up -d --scale kademliaNodes=49
