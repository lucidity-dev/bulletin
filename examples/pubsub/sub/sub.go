package main

import (
	"fmt"
	"os"
	"time"

	"nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/sub"
	"nanomsg.org/go-mangos/protocol/req"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"

	"github.com/golang/protobuf/proto"
	pb "github.com/apache8080/bulletin/protobuf"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func main() {
	url := "tcp://127.0.0.1:40899"
	var err error
	var server mangos.Socket
	var data []byte
	var msg []byte
	
	server, err = req.NewSocket()
	if err != nil {
		die("error creating socket: %s", err.Error())
	}

	server.AddTransport(ipc.NewTransport())
	server.AddTransport(tcp.NewTransport())
	if err = server.Dial(url); err != nil {
		die("cant dial on socket: %s", err)
	}

	request := &pb.Message{
		Cmd: pb.Message_GET,
		Args: "test1",
	}

	data, err = proto.Marshal(request)

	if err != nil {
		die("error creating protobuf: %s", err)
	}

	if err = server.Send(data); err != nil {
		die("error sending message: %s", err)
	}

	if msg, err = server.Recv(); err != nil {
		die("didn't receive any data: %s", err)
	}

	var body pb.Topic
	proto.Unmarshal(msg, &body)

	var sock mangos.Socket
	
	if sock, err = sub.NewSocket(); err != nil {
		die("can't get new sub socket: %s", err.Error())
	}

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial(body.Url); err != nil {
		die("can't dial on sub socket: %s", err.Error())
	}
	// Empty byte array effectively subscribes to everything
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		die("cannot subscribe: %s", err.Error())
	}
	for {
		if msg, err = sock.Recv(); err != nil {
			die("Cannot recv: %s", err.Error())
		}
		fmt.Printf("CLIENT: RECEIVED %s\n", string(msg))
	}

}
