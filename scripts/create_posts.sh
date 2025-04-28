#!/bin/bash

# Replace this with your real JWT token
USER_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDU4NDc4NDIsInVzZXJfaWQiOjN9.OjE3uI3SH_1-CQ5JEKdLkvyZbx8RrASGyiyFpaqsAmc"

for i in {1..10}
do
  curl -X POST http://localhost:8080/posts \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"content\":\"Post number $i from user1\"}"
  echo ""  # for newline after each post
done
