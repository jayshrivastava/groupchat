package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes"

	uuid "github.com/google/uuid"
	chat "github.com/jayshrivastava/groupchat/proto"
	application "github.com/jayshrivastava/groupchat/server/application"
	. "github.com/jayshrivastava/groupchat/server/context"
)

type ServerProps struct {
	Password string
	Port     string
}

type Server struct {
	chat.UnimplementedChatServer
	Password string
	Port     string
	Context  *Context
}

func CreateChatServer(props ServerProps) *Server {
	server := Server{
		Password: props.Password,
		Port:     props.Port,
		Context:  application.CreateApplicationContext(),
	}

	return &server
}

func GenerateToken() string {
	return uuid.New().String()
}

func (s *Server) Login(ctx context.Context, req *chat.LoginRequest) (*chat.LoginResponse, error) {

	if s.Password != req.Password {
		return nil, fmt.Errorf("Invalid server password: %s", req.Username)
	}

	token := GenerateToken()

	s.Context.ChannelRepository.Create(req.Username)
	s.Context.UserRepository.Create(req.Username, token, req.Group)
	s.Context.GroupRepository.CreateIfNotExists(req.Group)
	s.Context.GroupRepository.AddUserToGroup(req.Username, req.Group)

	// Login notification
	new_user_res := chat.StreamResponse{
		Timestamp: ptypes.TimestampNow(),
		Event: &chat.StreamResponse_ClientLogin{
			&chat.StreamResponse_Login{
				Username: req.Username,
				Group:    req.Group,
			},
		},
	}

	groupMembers, _ := s.Context.GroupRepository.GetGroupMembers(req.Group, req.Username)

	// Notify existing users about the new user
	for _, username := range groupMembers {
		stream, _ := s.Context.ChannelRepository.Get(username)
		stream <- new_user_res
	}

	// List existing users for the new user
	for _, username := range groupMembers {
		existing_user_res := chat.StreamResponse{
			Timestamp: ptypes.TimestampNow(),
			Event: &chat.StreamResponse_ClientExisting{
				&chat.StreamResponse_Existing{
					Username: username,
					Group:    req.Group,
				},
			},
		}
		stream, _ := s.Context.ChannelRepository.Get(req.Username)
		stream <- existing_user_res
	}

	return &chat.LoginResponse{Token: token}, nil
}

func (s *Server) Logout(ctx context.Context, req *chat.LogoutRequest) (*chat.LogoutResponse, error) {

	token, err := s.GetTokenFromContext(ctx)
	if err != nil || !s.Context.Authenticator.Authenticate(token, req.Username) {
		return nil, fmt.Errorf("Invalid token provided for %s", req.Username)
	}

	group, _ := s.Context.UserRepository.GetGroup(req.Username)

	s.Context.UserRepository.DeleteToken(req.Username)
	s.Context.UserRepository.DeleteGroup(req.Username)

	// Send logout notification
	res := chat.StreamResponse{
		Timestamp: ptypes.TimestampNow(),
		Event: &chat.StreamResponse_ClientLogout{
			&chat.StreamResponse_Logout{
				Username: req.Username,
				Group:    group,
			},
		},
	}

	groupMembers, _ := s.Context.GroupRepository.GetGroupMembers(group, req.Username)
	for _, username := range groupMembers {
		stream, _ := s.Context.ChannelRepository.Get(username)
		stream <- res
	}
	s.Context.GroupRepository.RemoveUserFromGroup(req.Username, group)

	return &chat.LogoutResponse{}, nil
}

func (s *Server) Stream(stream chat.Chat_StreamServer) error {
	token, err := s.GetTokenFromContext(stream.Context())

	if err != nil || !s.Context.Authenticator.IsTokenValid(token) {
		return fmt.Errorf("Invalid Token %s", token)
	}

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

		token, err := s.GetTokenFromContext(stream.Context())
		if err != nil || !s.Context.Authenticator.Authenticate(token, req.Username) {
			continue
		}

		res := chat.StreamResponse{
			Timestamp: ptypes.TimestampNow(),
			Event: &chat.StreamResponse_ClientMessage{
				&chat.StreamResponse_Message{
					Username: username,
					Group:    group,
					Message:  message,
				},
			},
		}

		groupMembers, _ := s.Context.GroupRepository.GetGroupMembers(group, username)
		for _, reciever_username := range groupMembers {
			if username != reciever_username {
				clientChannel, _ := s.Context.ChannelRepository.Get(reciever_username)
				clientChannel <- res
			}
		}

	}
}

func sender(s *Server, stream chat.Chat_StreamServer, wg *sync.WaitGroup, token string) {
	defer wg.Done()

	username, _ := s.Context.UserRepository.GetUsername(token)
	clientChannel, _ := s.Context.ChannelRepository.Get(username)

	for {
		res := <-clientChannel
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
	if value, found := meta["token"]; !found || len(value) == 0 || !ok {
		return "", fmt.Errorf("Missing Token in Request")
	}

	return meta["token"][0], nil
}
