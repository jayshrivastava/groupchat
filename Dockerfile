FROM golang:alpine

# for go.mod and go.sum
ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

WORKDIR /go/src/github.com/jayshrivastava/groupchat
COPY . .

RUN go mod download

RUN go build -o groupchat .

CMD ./groupchat