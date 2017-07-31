#!/usr/bin/env bash

#curl -X POST -H "Content-Type: application/json" http://localhost:18087/api/actions -d @move_data.json

host="localhost"
port=18087

jsonfile=./move_data.json

function do_move {
    if [ "X$1" != "X" ] ; then 
        jsonfile=$1
    fi
    action="actions"
    url="http://$host:$port/api/$action"
    curl -X POST -H "Content-Type: application/json" --data @$jsonfile $url 
    ret=$?
    if [ $ret -ne 0 ] ; then
        echo "fail"
    else
        echo "success"
    fi
}

do_move

do_move ./move2_data.json

