# hostsctl
A small tool to manage local /etc/hosts file and sync DNS records from various sources into it.

# TODO
* Implement backup/restore for the database
* Implement clear/clean for a block
* Update readme
* Create CI
* Publish to marketplaces

* cleanup hosts backend from not used dependencies
* implement GlobalOptions, move format determining logic there, make it customizable as per command
* How to implement build tasks https://www.alfusjaganathan.com/blogs/psake-build-automation-net-dotnet/


# Test Coverage

```
mkdir coverage
go test ./... -cover -coverpkg=./... -coverprofile=./coverage/prof.out 
go tool cover -html=./coverage/prof.out -o ./coverage/coverage.html
firefox ./coverage/coverage.html
```