package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/wfunc/util/xnet"
)

var waiter = sync.WaitGroup{}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf(`Usage: transport <local> <remote> <local> <remote> ...`)
		return
	}
	mappings := os.Args[1:]
	n := len(mappings) / 2
	InfoLog("transport starting %v mapping by %v", n, mappings)
	for i := 0; i < n; i++ {
		waiter.Add(1)
		go runMapping(mappings[i*2], mappings[i*2+1])
	}
	waiter.Wait()
	InfoLog("transport is done")
}

func runMapping(local, remote string) {
	defer waiter.Done()
	InfoLog("start mapping %v to %v", local, remote)
	var transporter xnet.Transporter
	if strings.HasPrefix(remote, "tcp://") {
		transporter = xnet.RawDialerF(net.Dial)
	} else if strings.HasPrefix(remote, "ws://") || strings.HasPrefix(remote, "wss://") {
		transporter = xnet.NewWebsocketDialer()
	} else {
		err := fmt.Errorf("not supported remote %v", remote)
		ErrorLog("mapping %v to %v fail with %v", local, remote, err)
		return
	}
	ln, err := net.Listen("tcp", local)
	if err != nil {
		ErrorLog("mapping %v to %v fail with %v", local, remote, err)
		return
	}
	var conn net.Conn
	for {
		conn, err = ln.Accept()
		if err != nil {
			break
		}
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			InfoLog("start transport %v to %v", conn.RemoteAddr(), remote)
			xerr := transporter.Transport(conn, remote)
			InfoLog("stop transport %v to %v by %v", conn.RemoteAddr(), remote, xerr)
		}()
	}
	InfoLog("mapping %v to %v is done with %v", local, remote, err)
}
