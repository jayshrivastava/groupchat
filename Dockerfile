FROM golang:1.13.5

WORKDIR /go/src/github.com/jayshrivastava/gRPC-terminal-chat

COPY . .

RUN go get github.com/golang/protobuf/proto \ 
    github.com/golang/protobuf/ptypes \
    github.com/golang/protobuf/ptypes/timestamp \
    google.golang.org/grpc \ 
    google.golang.org/grpc/codes \
    google.golang.org/grpc/metadata \
    google.golang.org/grpc/status

RUN go build -o main 

ENTRYPOINT ["./main"]