package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/golang/protobuf/proto"
	pb "github.com/lucidity-dev/bulletin/protobuf"

	"nanomsg.org/go-mangos/protocol/rep"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func getUrl(mc *memcache.Client) string {
	port := rand.Intn(59999) + 1
	_, miss := mc.Get(string(port))

	for miss == nil {
		port = rand.Intn(59999) + 1
		_, miss = mc.Get(string(port))
	}
	mc.Set(&memcache.Item{Key: string(port), Value: []byte("test")})

	return fmt.Sprintf("tcp://127.0.0.1:%d", port)
}

func main() {
	fmt.Println("Waiting for memcached connection")
	mc := memcache.New("127.0.0.1:11211")
	fmt.Println("Connected to memcached")
	url := "tcp://127.0.0.1:40899"

	sock, err := rep.NewSocket()
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

		if err != nil {
			die("bad request: %s", err)
		}

		proto.Unmarshal(msg, &body)

		switch body.Cmd {
		case pb.Message_HELP:
			sock.Send([]byte("hello"))
		case pb.Message_REGISTER:
			fmt.Println("registering")
			if _, hit := mc.Get(string(body.Args)); hit != nil {
				link := getUrl(mc)
				mc.Set(&memcache.Item{Key: string(body.Args), Value: []byte(link)})
				//TODO: autogenerate URL with random socket

				result := &pb.Topic{
					Name: string(body.Args),
					Url:  link,
					Err:  "",
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
					Url:  "",
					Err:  "ERROR: Topic already registered",
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
					Url:  string(item.Value),
					Err:  "",
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
					Url:  "",
					Err:  "ERROR: Topic not registered",
				}
				var res []byte
				res, err = proto.Marshal(result)

				sock.Send(res)
			}
		case pb.Message_FLUSH_ALL:
			fmt.Println("flushing all")
			err = mc.FlushAll()
			if err != nil {
				die("memecached error: %v", err)
			}
			sock.Send([]byte("FLUSHED ENTIRE CACHE"))
		default:
			sock.Send([]byte("INVALID REQUEST"))
		}
		//fmt.Println(body.Args)
	}
}
