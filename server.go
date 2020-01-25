package main 

import (
	"context" 
	"crypto/rand" 
	"fmt" 
	"io" 
	"net" 
	"sync" 

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes"
	
	chat "github.com/jayshrivastava/groupchat/proto"
)


type ServerProps struct {
	Password string
	Port string
}

type Server struct {
	chat.UnimplementedChatServer
	Password string
	Port string
	ClientLog *ClientLog
}

type ClientLog struct {
	ClientChannels map[string]chan chat.StreamResponse 
	ClientTokens map [string]string
	ClientGroups map [string]string
	GroupMembers map [string](map [string]bool)
	// Readers/Writers lock on all maps
	Mutex sync.RWMutex 
}

func CreateChatServer(props ServerProps) *Server {
	server := Server{
		Password: props.Password,
		Port: props.Port,
	}

	server.ClientLog = &ClientLog{
		ClientChannels: make(map [string]chan chat.StreamResponse),
		ClientTokens: make(map[string]string),
		ClientGroups: make(map[string]string),
		GroupMembers: make(map[string](map[string]bool)),
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

	if s.Password != req.Password {
		return nil, fmt.Errorf("Invalid server password: %s", req.Username)
	}

	// write lock
	s.ClientLog.Mutex.Lock()

	if _, found := s.ClientLog.ClientChannels[req.Username]; found {
		s.ClientLog.Mutex.Unlock()
		return nil, fmt.Errorf("Username %s already exists", req.Username)
	}

	token := GenerateToken()

	s.ClientLog.ClientChannels[req.Username] = make(chan chat.StreamResponse, 5)
	s.ClientLog.ClientTokens[token] = req.Username
	s.ClientLog.ClientGroups[req.Username] = req.Group

	if _, found := s.ClientLog.GroupMembers[req.Group]; !found {
		s.ClientLog.GroupMembers[req.Group] = make(map[string]bool)
	}
	s.ClientLog.GroupMembers[req.Group][req.Username] = true

	s.ClientLog.Mutex.Unlock()

	// Login notification 
	new_user_res := chat.StreamResponse{
		Timestamp: ptypes.TimestampNow(),
		Event: &chat.StreamResponse_ClientLogin{
			&chat.StreamResponse_Login{
				Username: req.Username,
				Group: req.Group,
			},
		},
	}

	s.ClientLog.Mutex.RLock()

	// Notify existing users about the new user
	for username, stream := range s.ClientLog.ClientChannels {
		if username != req.Username && s.ClientLog.ClientGroups[username] == req.Group {
			stream <- new_user_res
		}
	}

	// List existing users for the new user
	for username, _ := range s.ClientLog.GroupMembers[req.Group] { 
		if username != req.Username && s.ClientLog.ClientGroups[username] == req.Group {
			existing_user_res := chat.StreamResponse{
				Timestamp: ptypes.TimestampNow(),
				Event: &chat.StreamResponse_ClientExisting{
					&chat.StreamResponse_Existing{
						Username: username,
						Group: req.Group,
					},
				},
			}
			s.ClientLog.ClientChannels[req.Username] <- existing_user_res
		}
	}
	s.ClientLog.Mutex.RUnlock()

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

	group := s.ClientLog.ClientGroups[req.Username]

	delete(s.ClientLog.ClientChannels, req.Username)
	delete(s.ClientLog.ClientTokens, token)
	delete(s.ClientLog.ClientGroups, req.Username)
	
	s.ClientLog.Mutex.Unlock()

	// Send logout notification
	res := chat.StreamResponse{
		Timestamp: ptypes.TimestampNow(),
		Event: &chat.StreamResponse_ClientLogout{
			&chat.StreamResponse_Logout{
				Username: req.Username,
				Group: group,
			},
		},
	}

	s.ClientLog.Mutex.RLock()
	for username, stream := range s.ClientLog.ClientChannels {
		if s.ClientLog.ClientGroups[username] == group {
			stream <- res
		}
	}
	s.ClientLog.Mutex.RUnlock()

	return &chat.LogoutResponse{}, nil
}

func (s *Server) Stream(stream chat.Chat_StreamServer) error {
	token, err := s.GetTokenFromContext(stream.Context())

	if err == nil {
		wg := sync.WaitGroup{}

		wg.Add(1)
		go reciever(s, stream, &wg, token)

		wg.Add(1)
		go sender(s, stream, &wg, token)

		wg.Wait()
	}
	return nil
}

func reciever(s *Server, stream chat.Chat_StreamServer, wg *sync.WaitGroup, token string) {
	defer wg.Done()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			break
		}

		username, group, message := req.Username, req.Group, req.Message

		s.ClientLog.Mutex.RLock()
		if _, found := s.ClientLog.ClientTokens[token]; !found || s.ClientLog.ClientTokens[token] != req.Username {
			s.ClientLog.Mutex.RUnlock()
			break
		}

		res := chat.StreamResponse{
			Timestamp: ptypes.TimestampNow(),
			Event: &chat.StreamResponse_ClientMessage{
				&chat.StreamResponse_Message{
					Username: username,
					Group: group,
					Message: message,
				},
			},
		}

		for reciever_username, clientChannel := range s.ClientLog.ClientChannels {
			if username != reciever_username && s.ClientLog.ClientGroups[reciever_username] == group {
				clientChannel <- res
			}
		}

		s.ClientLog.Mutex.RUnlock()
	}
}

func sender(s *Server, stream chat.Chat_StreamServer, wg *sync.WaitGroup, token string) {
	defer wg.Done()

	s.ClientLog.Mutex.RLock()
	username := s.ClientLog.ClientTokens[token]
	clientChannel := s.ClientLog.ClientChannels[username]
	s.ClientLog.Mutex.RUnlock()
	
	for {
		res := <- clientChannel
		stream.Send(&res)
	}
}

func ServerMain(props ServerProps) {

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", props.Port))
	if err != nil {
		fmt.Errorf("failed to listen: %v", err)
    }

    server := grpc.NewServer()
	chat.RegisterChatServer(server, CreateChatServer(props))
	server.Serve(lis)
}

func (s *Server) GetTokenFromContext(ctx context.Context) (string, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if  value, found := meta["token"]; !found || len(value) == 0 || !ok {
		return "", fmt.Errorf("Missing Token in Request")
	}

	return meta["token"][0], nil
}
