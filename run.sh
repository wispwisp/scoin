#!/bin/bash

docker run --rm -it -v$(pwd -P)/src:/home/user/src -p127.0.0.1:8090:8090 golang:1.19 /bin/bash

