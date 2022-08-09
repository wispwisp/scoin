#!/bin/bash

clear && curl -X POST 127.0.0.1:8090/transaction \
              -H 'Content-Type: application/json' \
              -d '{"from":"me","to":"him","amount":55}'

