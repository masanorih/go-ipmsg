go-ipmsg
==============

IPMSG protocol implementation for golang

To install, simply issue a `go get`:

```
go get github.com/masanorih/go-ipmsg
```

By default importing `github.com/masanorih/go-ipmsg` will import package
`ipmsg`

```go
import (
    "log"
    "github.com/masanorih/go-ipmsg"
)

conf := ipmsg.NewIPMSGConf()
ipmsg, err := ipmsg.NewIPMSG(conf)
if err != nil {
   log.Fatalf("Failed to start ipmsg: %v", err)
}
defer ipmsg.Close()

```

`go-ipmsg` is a port of [Net::IPMessenger](https://metacpan.org/release/Net-IPMessenger)

When you create a new struct via `NewIPMSG()` a new ipmsg instance is
automatically setup and launched. Don't forget to call `Close()` on this
struct to close the ipmsg server

If you want to customize the configuration, create a new config and set each
field on the struct:

```go

conf := ipmsg.NewIPMSGConf()
conf.Port = 12425

// Starts ipmsg listening on port 12425
ipmsg, err := ipmsg.NewIPMSG(conf)
```

TODO
====
add encription support
