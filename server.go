package main 

import (
	"context"
	"flag"
	"sync"
	"fmt"
	"net"
	"syscall"
	"os"
	// "io"
	"crypto/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	// "github.com/golang/protobuf/ptypes"
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
	ClientGroups map [string]string
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
		ClientTokens: make(map[string]string),
		ClientGroups: make(map[string]string),
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
	// write lock
	s.ClientLog.Mutex.Lock()

	if _, found := s.ClientLog.ClientChannels[req.Username]; found {
		s.ClientLog.Mutex.Unlock()
		return nil, fmt.Errorf("Username %s already exists", req.Username)
	}

	token := GenerateToken()

	s.ClientLog.ClientChannels[req.Username] = make(chan chat.StreamResponse)
	s.ClientLog.ClientTokens[token] = req.Username
	s.ClientLog.ClientGroups[req.Username] = req.Group

	s.ClientLog.Mutex.Unlock()

	return &chat.LoginResponse{Token: token}, nil
}

func (s *Server) Logout(ctx context.Context, req *chat.LogoutRequest) (*chat.LogoutResponse, error) {
	
	token, err := s.GetTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Missing token provided for %s", req.Username)
	}
	
	s.ClientLog.Mutex.Lock()

	if _, found := s.ClientLog.ClientTokens[token]; !found || s.ClientLog.ClientTokens[token] != req.Username {
		s.ClientLog.Mutex.Unlock()
		return nil, fmt.Errorf("Invalid Token for %s", req.Username)
	}

	delete(s.ClientLog.ClientChannels, req.Username)
	delete(s.ClientLog.ClientTokens, token)
	delete(s.ClientLog.ClientGroups, req.Username)
	
	s.ClientLog.Mutex.Unlock()

	return &chat.LogoutResponse{}, nil
}

// func (s *Server) Stream(stream chat.Chat_StreamServer) error {
// 	token, err := s.GetTokenFromContext(stream.Context())
	
// 	if err != nil {
// 		wg := sync.WaitGroup{}

// 		wg.Add(1)
// 		go reciever(s, stream, &wg, token)

// 		wg.Add(1)
// 		go sender(s, stream, &wg, token)

// 		wg.Wait()
// 	}
// 	return nil
// }

// func reciever(s *Server, stream chat.Chat_StreamServer, wg *sync.WaitGroup, token string) {
// 	defer wg.Done()
// 	for {
// 		req, err := stream.Recv()
// 		if err == io.EOF {
// 			continue
// 		} else if err != nil {
// 			break
// 		}

// 		username, group, message := req.Username, req.Group, req.Message

// 		s.ClientLog.Mutex.RLock()
// 		if _, found := s.ClientLog.ClientTokens[token]; !found || s.ClientLog.ClientTokens[token] != req.Username {
// 			s.ClientLog.Mutex.RUnlock()
// 			break
// 		}

// 		res := chat.StreamResponse{
// 			Timestamp: ptypes.TimestampNow(),
// 			Event: &chat.StreamResponse_ClientMessage{
// 				&chat.StreamResponse_Message{
// 					Username: username,
// 					Group: group,
// 					Message: message,
// 				},
// 			},
// 		}

// 		for username, stream := range s.ClientLog.ClientChannels {
// 			if s.ClientLog.ClientGroups[username] == group {
// 				stream <- res
// 			}
// 		}

// 		s.ClientLog.Mutex.RUnlock()
// 	}
// }

// func sender(s *Server, stream chat.Chat_StreamServer, wg *sync.WaitGroup, token string) {
// 	defer wg.Done()

// 	s.ClientLog.Mutex.RLock()
// 	username := s.ClientLog.ClientTokens[token]
// 	s.ClientLog.Mutex.RUnlock()
	
// 	for {
// 		res := <- s.ClientLog.ClientChannels[username]
// 		stream.Send(&res)
// 	}
// }

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:5000"))
	if err != nil {
		fmt.Errorf("failed to listen: %v", err)
    }

    server := grpc.NewServer()
	chat.RegisterChatServer(server, CreateChatServer())
	server.Serve(lis)
}

func (s *Server) GetTokenFromContext(ctx context.Context) (string, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if  value, found := meta["token"]; !found || len(value) == 0 || !ok {
		return "", fmt.Errorf("Missing Token in Request")
	}

	return meta["token"][0], nil
}
