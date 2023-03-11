# exdhcp

DHCP lease exhausting test tool

# Usage

1. determine the target interface

```
$ ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
13: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
```

2. Run

- 10 times (default)

  ```
  $ sudo ./exdhcp -if eth0
  ```

  CSV result file will be saved.

- n times

  ```
  $ sudo ./exdhcp -if eth0 -num 100
  ```

- run until lease exhausts

  ```
  $ sudo ./exdhcp -if eth0 -num 0
  ```

# Options

```
$ ./exdhcp -help
Usage of ./exdhcp:
  -csv
        Export leases to CSV file (default true)
  -if string
        Interface name
  -num int
        Number of tries. Set 0 to do exhaustion attack. (default 10)
  -release
        Release DHCP lease (default true)
  -server string
        DHCP server IP
  -timeout int
        Timeout seconds to wait DHCP offer (default 10)
  -verbose
        Print debug logs
```

# Build

```
$ go build ./cmd/dhcpcd/dhcpcd.go
```

# v1 milestone

- [ ] support Windows (use pcap)

# Reference

https://support.huawei.com/enterprise/en/doc/EDOC1100055047/c8e6549d/dhcp-dos-attack-by-changing-chaddr
