#!/bin/bash

## run turbo-simulator from docker
port=18087
image=beekman9527/simulator:latest
docker run -d -p $port:8087 $image --v 2 
