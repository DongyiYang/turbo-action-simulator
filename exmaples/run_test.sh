#!/bin/bash
set -x

host=localhost
port=18087
url="http://$host:$port/api/discovery"

discovery=discovery_data.json
curl -X POST -H "Content-Type: application/json" $url -d @$discovery
ret=$?
if [ $ret -ne 0 ] ; then
    echo "failed to call server"
else
    echo "success"
fi
