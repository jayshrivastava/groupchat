# groupchat
A multithreaded groupchat server/client implementation made with Go and gRPC.

This project leverages concurrency constructs in Go as well as socket-like bidirectional streaming in gRPC to exchange messages.

The server uses a [hexagonal architecture](https://medium.com/sciforce/another-story-about-microservices-hexagonal-architecture-23db93fa52a2). Components like storage, authentication, and inbound/outbound ports are abstracted away from the core server logic using interfaces. This makes the server "agnostic to the outside world" \, and it will indivdual components easier to test in isolation (when finally get around to writing tests :sunglasses:)

## Usage
### Without Docker
Install Dependencies  
`go mod download`

Build Executable  
`go build -o groupchat`

Start Server  
`go run server.go -spass <password> -h <hostname>`  
ex. `./groupchat -s -spass "password" -port "5000"`

Start Client   
`./gRPC-terminal-chat -u <username> -g <groupname> -p <password> -h <hostname>`  
ex. `./groupchat -u "jay" -pass "abc" -g "my group" -spass "password" -url "localhost:5000"`

### With Docker
Build container 
`docker build --tag=groupchat .`

Run Server:  
`docker run -it -p 5000:5000 groupchat:latest ./groupchat -s -spass "password" -port "5000"`

Run Some Clients:   
`docker run -it -i --net=host groupchat:latest ./groupchat -u "user 1" -pass "abc" -g "my group" -spass "password"  -url "localhost:5000"` 

`docker run -it -i --net=host groupchat:latest ./groupchat -u "user 2" -pass "abc" -g "my group" -spass "password"  -url "localhost:5000"` 

## Development  

Generate Protocol Bufffer    
`protoc -I proto/ proto/groupchat.proto --go_out=plugins=grpc:proto`

Formatting  
`go fmt ./...`



