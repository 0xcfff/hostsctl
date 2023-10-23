# hostsctl
A small tool to manage aliases records in local /etc/hosts file (or %SYSTEM32%\Drivers\etc\hosts for Windows). The purpose of creation of this tool was to simplify manual aliases managements as well as making it possible to entirely automate this process. Once installed, the tool can:
* Manage individual DNS aliases records
* Manage blocks of DNS aliases records
* Format /etc/hosts file
* Backup/restore /etc/hosts config

# Installation
The only option at the moment is to build the tool localy and copy it to a folder included into system path.
```
go build -o out/bin/hostsctl ./cmd/hostsctl
mv out/bin/hostsctl /usr/local/bin/
```

# Usage
The simplest option to use the tool is to manage individual IP aliases manually:
```
# backup database file
hostsctl database backup

# add an alias record into the alias database
hostsctl alias add 127.0.0.1 pet-project.local

# print the database
hostsctl database print

# revert the database changes
hostsctl database restore
```

A more sophisticated example of the tool usage might be syncing aliases from a K8S cluster directly into /etc/hosts
```
# backup database file
hostsctl database backup

# ensure a special block to keep local clusted aliases exists and it is empty
hostsctl block add --name k8s-local --force
hostsctl block clear --name k8s-local

# import k8s nodes aliases to the aliases database
kubectl get nodes \
    -o custom-columns='"IP":.status.addresses[?(@.type=="InternalIP")].address,ALIAS:.status.addresses[?(@.type=="Hostname")].address' \
    --no-headers \
    | hostsctl alias add --block k8s-local

# import aliases for all LoadBalancer services registered in k8s
kubectl get svc -A \
    -o jsonpath='{range .items[?(@.status.loadBalancer.ingress[0].ip)]}{.status.loadBalancer.ingress[0].ip} {.metadata.name}{"\n"}' \
    | hostsctl alias add --block k8s-local

# print the database
hostsctl database print

# revert the database changes
hostsctl database restore
```

# Known issues
No known issues at this point.

# Getting help
Report issues in the repository [issue tracker](https://github.com/0xcfff/hostsctl/issues)

# Getting involved
General instructions on how to contribute can be found in [CONTRIBUTING](CONTRIBUTING.md).

# TODO
* Implement clear/clean for a block
* Implement check on deletion of a system alias
* Update readme
* Create CI
* Publish to marketplaces

* cleanup hosts backend from not used dependencies
* implement GlobalOptions, move format determining logic there, make it customizable as per command
* How to implement build tasks https://www.alfusjaganathan.com/blogs/psake-build-automation-net-dotnet/