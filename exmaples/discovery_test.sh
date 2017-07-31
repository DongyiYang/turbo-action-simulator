#!/usr/bin/env bash

curl -X POST -H "Content-Type: application/json" http://localhost:1234/api/discovery -d @discovery_data.json
