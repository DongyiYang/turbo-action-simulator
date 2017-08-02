#!/usr/bin/env bash

port=18087
curl -X POST -H "Content-Type: application/json" http://localhost:$port/api/discovery -d @discovery_data.json
