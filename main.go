package main

import (
	"fmt"
	"os"
	"strings"
	
	"nanomsg.org/go-mangos/protocol/rep"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"

	"github.com/golang/protobuf/proto"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func main() {
	url := "tcp://127.0.0.1:40899"
	
	sock, err := rep.NewSocket();
	if (err != nil) {
		die("can't start up server: %s", err)
	}

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Listen(url); err != nil {
		die("can't start up server: %s", err)
	}

	for {
		msg, err := sock.Recv()
		if (err != nil) {
			die("bad request: %s", err)
		}
		
		cmd := strings.Fields(string(msg))
		if len(cmd) == 0 {
			sock.Send([]byte("INVALID REQUEST"))
		}

		fmt.Println(cmd[0])
		
		switch cmd[0] {
		case "TEST":
			sock.Send([]byte("hello"))
		case "REGISTER":
			
		default:
			sock.Send([]byte("INVALID REQUEST"))
		}

	}
}
