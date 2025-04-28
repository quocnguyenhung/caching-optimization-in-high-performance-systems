#!/bin/bash

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDU4NTk1MTYsInVzZXJfaWQiOjR9.FSO1bLdHLAKmZcrz1h3eNDmilF3amyzMhIjTvnh-Wvw"

for i in {1..5}
do
  curl -X POST http://localhost:8080/posts/$i/like \
    -H "Authorization: Bearer $TOKEN"
done
