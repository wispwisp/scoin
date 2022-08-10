#!/bin/bash

curl -X POST 127.0.0.1:8090/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8888"}'
