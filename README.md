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


## Examples
### Registering and Getting a topic
This example registers a topic with the server and then makes a GET request for the topics URL.  
Building and Running (assumes the server is already running):
```
cd bulletin/examples/simple
go build .
./simple
```

### Pub/Sub Registration and Getting topic
This example registers a publisher that begins publishing on the auto generated URL, and then a subscriber gets the URL from the server and begins subscribing.
Building and Running the publisher (assumes the server is already running):
```
cd bulletin/examples/pubsub/pub
go build .
./pub
```
Building and Running the subscriber (assumes the server is already running):
```
cd bulletin/examples/pubsub/sub
go build .
./sub
```



