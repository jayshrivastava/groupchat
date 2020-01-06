# terminal-chat

A terminal-based chat server and client written in Go using gRPC. 

This project leverages concurrency constructs in Go as well as socket-like bidirectional streaming in gRPC to exchange messages between clients.

## Usage
<b> Install Dependencies </b>
`go get -d ./...`

<b> Build Executable </b>
`go build`

<b> Start Server </b>
`go run server.go -p <password> -h <hostname>`
ex. `./gRPC-terminal-chat -s -p "password" -h "0.0.0.0:5000"`

<b> Start Client </b>
`go run client.go -u jayants -g beans -p test -h abcd`
ex. `./gRPC-terminal-chat -u "jay" -g "my group" -p "password" -h "0.0.0.0:5000"`


### Useful Commands
<b> Add protoc to path </b>
`export PATH=$PATH:$HOME/go/bin`

<b> Generate Protocol Bufffer </b>
`protoc -I chat/ chat/chat.proto --go_out=plugins=grpc:chat`

### References 

Rodaine, “grpc-chat,” GitHub, 13-Oct-2017. [Online]. Available: github.com/rodaine/grpc-chat. [Accessed: 03-Jan-2020].

