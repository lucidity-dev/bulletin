package main

import (
	"fmt"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	
	"github.com/golang/protobuf/proto"
	pb "github.com/apache8080/bulletin/protobuf"

	"nanomsg.org/go-mangos/protocol/rep"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func main() {
	mc := memcache.New("127.0.0.1:8000")
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
			fmt.Println("registering")
			if _, hit := mc.Get(string(body.Args)); hit != nil {
				mc.Set(&memcache.Item{Key: string(body.Args), Value: []byte("url")})
				//TODO: autogenerate URL with random socket
				result := &pb.Topic{
					Name: string(body.Args),
					Url: "url",
					Err: "",
				}
				var res []byte
				res, err = proto.Marshal(result)
				if err != nil {
					sock.Send([]byte("registering failed: data marshalling error"))
				}
				sock.Send(res)
			} else {
				result := &pb.Topic{
					Name: "",
					Url: "",
					Err: "ERROR: Topic already registered",
				}
				var res []byte
				res, err = proto.Marshal(result)

				sock.Send(res)
			}
		case pb.Message_GET:
			fmt.Println("getting")
			if item, hit := mc.Get(string(body.Args)); hit == nil {
				result := &pb.Topic{
					Name: string(body.Args),
					Url: string(item.Value),
					Err: "",
				}
				var res []byte
				res, err = proto.Marshal(result)
				if err != nil {
					sock.Send([]byte("getting failed: data marshalling error"))
				}
				sock.Send(res)
			} else {
				result := &pb.Topic{
					Name: "",
					Url: "",
					Err: "ERROR: Topic not registered",
				}
				var res []byte
				res, err = proto.Marshal(result)

				sock.Send(res)
			}
		default:
			sock.Send([]byte("INVALID REQUEST"))
		}
		//fmt.Println(body.Args)
	}
}
