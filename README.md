# terminal-chat 
A multithreaded groupchat server/client implementation made with Go and gRPC.

This project leverages concurrency constructs in Go as well as socket-like bidirectional streaming in gRPC to exchange messages.

The server uses a [hexagonal architecture](https://medium.com/sciforce/another-story-about-microservices-hexagonal-architecture-23db93fa52a2). Components like storage, authentication, and inbound/outbound ports are abstracted away from the core server logic using interfaces. This makes the server "agnostic to the outside world" \, and it will indivdual components easier to test in isolation (when finally get around to writing tests :sunglasses:)

### Requirements  
- at least `golang v1.13.5`
- Network must support HTTP/2

### Usage
<b> Install Dependencies </b>  
`go get -d ./...`

<b> Build Executable </b>  
`go build -o groupchat`

<b> Start Server </b>  
`go run server.go -p <password> -h <hostname>`  
ex. `./groupchat -s -p "password" -port "5000"`

<b> Start Client </b>   
`./gRPC-terminal-chat -u <username> -g <groupname> -p <password> -h <hostname>`  
ex. `./groupchat -u "jay" -g "my group" -p "password" -h "localhost:5000"`

### Development  

<b> Generate Protocol Bufffer </b>  
`protoc -I proto/ proto/groupchat.proto --go_out=plugins=grpc:proto`

<b> Formatting </b>
`go fmt ./...`


