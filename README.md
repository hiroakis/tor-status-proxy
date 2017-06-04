# tor-status-proxy

It caches the exit node list and the all node list from https://torstatus.blutmagie.de/

# Demo

https://tor-status-proxy.herokuapp.com/

# Installation

Build and run. See Makefile.

* options

```
-h string
      The listen IP. (default "0.0.0.0")
-i int
      The polling interval in sec. (default 3600)
-p int
      The listen port. (default 9000)
```

# Example

```
$ curl https://tor-status-proxy.herokuapp.com/all
 => all nodes list

$ curl https://tor-status-proxy.herokuapp.com/exit
 => exit nodes list

$ curl -XPOST -i -d "ip=IPADDR" https://tor-status-proxy.herokuapp.com/istor
 => If your IP is tor node, the response code will be 200.

$ curl -XPOST -i -d "ip=IPADDR" https://tor-status-proxy.herokuapp.com/isexit
 => If your IP is tor exit node, the response code will be 200.
```

# License

MIT
