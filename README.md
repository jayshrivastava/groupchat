# terminal-chat 
A terminal-based chat server and client written in Go using gRPC. 

This project leverages concurrency constructs in Go as well as socket-like bidirectional streaming in gRPC to exchange messages between clients.

### Requirements  
- go v1.13.5
- Network must support HTTP/2

### Usage
<b> Install Dependencies </b>  
`go get -d ./...`

<b> Build Executable </b>  
`go build -o groupchat`

<b> Start Server </b>  
`go run server.go -spass <password> -h <hostname>`  
ex. `./groupchat -s -spass "password" -port "5000"`

<b> Start Client </b>   
`./gRPC-terminal-chat -u <username> -g <groupname> -p <password> -h <hostname>`  
ex. `./groupchat -u "jay" -pass "abc" -g "my group" -spass "password"  -url "localhost:5000"`

### Development  

<b> Generate Protocol Bufffer </b>  
`protoc -I proto/ proto/groupchat.proto --go_out=plugins=grpc:proto`

<b> Formatting </b>
`go fmt ./...`


