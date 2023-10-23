# /bin/bash

mkdir -p out/coverage
go test ./... -cover -coverpkg=./... -coverprofile=./out/coverage/prof.out 
go tool cover -html=./out/coverage/prof.out -o ./out/coverage/coverage.html
firefox ./out/coverage/coverage.html