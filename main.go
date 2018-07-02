package main

import (
	"fmt"
	"os"
	
	"nanomsg.org/go-mangos/protocol/rep"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"

	"github.com/golang/protobuf/proto"
	pb "github.com/apache8080/bulletin/protobuf"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func main() {
	url := "tcp://127.0.0.1:40899"
	
	sock, err := rep.NewSocket();
	if err != nil {
		die("can't start up server: %s", err)
	}

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Listen(url); err != nil {
		die("can't start up server: %s", err)
	}

	for {
		msg, err := sock.Recv()
		var body pb.Message
		
		if (err != nil) {
			die("bad request: %s", err)
		}
		
		proto.Unmarshal(msg, &body)

		switch body.Cmd {
		case pb.Message_HELP:
			sock.Send([]byte("hello"))
		case pb.Message_REGISTER:
			sock.Send([]byte("registering"))
		case pb.Message_GET:
			sock.Send([]byte("getting"))
		default:
			sock.Send([]byte("INVALID REQUEST"))
		}
		fmt.Println(body.Args)
	}
}
