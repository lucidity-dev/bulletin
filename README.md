# bulletin

A simple service discovery server  for nanomsg pub/sub sockets using go, nanomsg, and protobuf.

## Building and running the server
Building:
```
git clone https://github.com/apache8080/bulletin.git
cd bulletin
go build .
```

Running:
```
./bulletin
```

In a seperate shell run memcached: `memcached -p 8000`
