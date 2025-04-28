#!/bin/bash

for i in {1..10}
do
  curl -X POST http://localhost:8080/signup \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"user$i\", \"password\":\"password$i\"}"
  echo ""
done
