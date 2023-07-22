# /bin/bash

ORIG_DIR=$(pwd)
mkdir -p coverage
go test ./... -cover -coverpkg=./... -coverprofile=./coverage/prof.out 
go tool cover -html=./coverage/prof.out -o ./coverage/coverage.html
firefox ./coverage/coverage.html