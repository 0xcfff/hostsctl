![CI Status](https://github.com/0xcfff/hostsctl/actions/workflows/ci.yaml/badge.svg)
___

# hostsctl
![](./demo.gif)

A small tool to manage alias records in local /etc/hosts file (or %SYSTEM32%\Drivers\etc\hosts for Windows). The purpose of creation of this tool was to simplify manual aliases managements as well as making it possible to entirely automate this process. Once installed, the tool can:
* Manage individual DNS alias records
* Manage blocks of DNS alias records
* Format /etc/hosts file
* Backup/restore /etc/hosts config

# Installation
There are several options to install the tool.

## Manual Installation
Download [latest release](https://github.com/0xcfff/hostsctl/releases/latest) from GitHub for your OS version. Unarchive the binaries and copy to a folder added to path environment variable.


## Building from Source
One of the options is to build the tool localy and copy it to a folder included into system path.
```
git checkout 
go build -o out/bin/hostsctl ./cmd/hostsctl
mv out/bin/hostsctl /usr/local/bin/
```

# Usage
The simplest option to use the tool is to manage individual IP aliases manually:
```
# backup database file
hostsctl database backup

# add an alias record into the alias database
hostsctl alias add 127.0.0.1 pet-project2.local

# print all aliases from /etc/hosts
hostsctl alias list

# revert the database changes
hostsctl database restore
```

Sample output:
```
GRP  SYS  IP              ALIAS
[1]  +    127.0.0.1       localhost
          127.0.1.1       home-laptop
[3]  +    ::1             ip6-localhost
     +    ::1             ip6-loopback
          fe00::0         ip6-localnet
          ff00::0         ip6-mcastprefix
          ff02::1         ip6-allnodes
          ff02::2         ip6-allrouters
[4]       192.168.100.64  user-3d2b.local
          192.168.100.64  chart-example.local
[5]       192.168.100.8   home-laptop
          192.168.100.64  zipkin
          127.0.0.1       pet-project.local
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

# print all aliases from /etc/hosts
hostsctl alias list

# revert the database changes
hostsctl database restore
```

# Known issues
No known issues at this point.

# Getting help
Report issues in the repository [issue tracker](https://github.com/0xcfff/hostsctl/issues)

# Getting involved
General instructions on how to contribute can be found in [CONTRIBUTING](CONTRIBUTING.md).
