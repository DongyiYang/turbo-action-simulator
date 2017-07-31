#!/usr/bin/env bash

curl -X POST -H "Content-Type: application/json" http://localhost:1234/api/actions -d @move_data.json
