#!/bin/bash

# 8090

curl -X POST 127.0.0.1:8090/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8091"}'

curl -X POST 127.0.0.1:8090/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8092"}'

# 8091

curl -X POST 127.0.0.1:8091/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8090"}'

curl -X POST 127.0.0.1:8091/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8092"}'

# 8092

curl -X POST 127.0.0.1:8092/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8090"}'

curl -X POST 127.0.0.1:8092/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8091"}'
