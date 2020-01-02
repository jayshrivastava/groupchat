package main 

import (
	"context"
	"flag"
	"sync"
	"fmt"
	"net"
	"syscall"
	"os"
	"crypto/rand"
	"google.golang.org/grpc"
	chat "./chat"
)

type Server struct {
	chat.UnimplementedChatServer
	Password *string
	Host *string
	ClientLog *ClientLog
}

type ClientLog struct {
	ClientChannels map[string]chan chat.StreamResponse 
	ClientTokens map [string]string
	Mutex sync.RWMutex
}

func ServerError(e error) {
	fmt.Printf("%s\n", e.Error())
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
}

func CreateChatServer() *Server {
	server := Server{
		Password: flag.String("p", "", "Server Password"),
		Host: flag.String("h", "", "Host"),
	}

	flag.Parse();

	if (*server.Password == "" || *server.Host == "") {
		ServerError(fmt.Errorf("Missing Flags"))
	}

	server.ClientLog = &ClientLog{
		ClientChannels: make(map [string]chan chat.StreamResponse),
		ClientTokens: make(map [string]string),
	}
	
	return &server
}

func GenerateToken() string {
	b := make([]byte, 5)
   	rand.Read(b)
	s := fmt.Sprintf("%X", b)
	return s
}

func (s *Server) Login(ctx context.Context, req *chat.LoginRequest) (*chat.LoginResponse, error) {
	s.ClientLog.Mutex.Lock()

	if _, found := s.ClientLog.ClientChannels[req.Username]; found {
		s.ClientLog.Mutex.Unlock()
		return nil, fmt.Errorf("Username %s already exists", req.Username)
	}

	token := GenerateToken()

	s.ClientLog.ClientChannels[req.Username] = make(chan chat.StreamResponse)
	s.ClientLog.ClientTokens[req.Username] = token

	s.ClientLog.Mutex.Unlock()

	return &chat.LoginResponse{Token: token}, nil
}

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:5000"))
	if err != nil {
		fmt.Errorf("failed to listen: %v", err)
    }

    server := grpc.NewServer()
	chat.RegisterChatServer(server, CreateChatServer())
	server.Serve(lis)
}

// API METHODS
