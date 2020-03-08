package application

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes"

	chat "github.com/jayshrivastava/groupchat/proto"
	authentication "github.com/jayshrivastava/groupchat/server/authentication"
	repositories "github.com/jayshrivastava/groupchat/server/repositories"
)

type Server struct {
	chat.UnimplementedChatServer
	ServerPassword    string
	Port              string
	ChannelRepository repositories.ChannelRepository
	GroupRepository   repositories.GroupRepository
	UserRepository    repositories.UserRepository
	Authenticator     authentication.Authenticator
}

func (s *Server) Login(ctx context.Context, req *chat.LoginRequest) (*chat.LoginResponse, error) {

	if s.ServerPassword != req.ServerPassword {
		return nil, fmt.Errorf("Invalid server password: %s", req.ServerPassword)
	}

	token := s.Authenticator.GenerateToken()
	if s.UserRepository.DoesUserExist(req.Username) {
		fmt.Println(s.UserRepository.IsLoggedIn(req.Username))
		if loggedIn, _ := s.UserRepository.IsLoggedIn(req.Username); loggedIn {
			return nil, fmt.Errorf("Invalid password for existing user %s", req.Username)
		}
		if valid, _ := s.UserRepository.CheckPassword(req.Username, req.UserPassword); valid {
			s.UserRepository.SetUserData(req.Username, token, req.Group)
		} else {
			return nil, fmt.Errorf("Invalid password for existing user %s", req.Username)
		}
	} else {
		s.UserRepository.Create(req.Username, token, req.Group, req.UserPassword, true)
	}
	s.ChannelRepository.Open(req.Username)
	s.GroupRepository.CreateIfNotExists(req.Group)
	s.GroupRepository.AddUserToGroup(req.Username, req.Group)

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

	groupMembers, _ := s.GroupRepository.GetGroupMembers(req.Group, req.Username)

	// Notify existing users about the new user
	for _, username := range groupMembers {
		stream, _ := s.ChannelRepository.Get(username)
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
		stream, _ := s.ChannelRepository.Get(req.Username)
		stream <- existing_user_res
	}

	return &chat.LoginResponse{Token: token}, nil
}

func (s *Server) Logout(ctx context.Context, req *chat.LogoutRequest) (*chat.LogoutResponse, error) {

	token, err := s.GetTokenFromContext(ctx)
	if err != nil || !s.Authenticator.AuthenticateToken(token, req.Username) {
		return nil, fmt.Errorf("Invalid token provided for %s", req.Username)
	}

	group, _ := s.UserRepository.GetGroup(req.Username)

	s.UserRepository.DeleteToken(req.Username)
	s.UserRepository.DeleteGroup(req.Username)
	s.UserRepository.LogOut(req.Username)

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

	groupMembers, _ := s.GroupRepository.GetGroupMembers(group, req.Username)
	for _, username := range groupMembers {
		stream, _ := s.ChannelRepository.Get(username)
		stream <- res
	}
	s.GroupRepository.RemoveUserFromGroup(req.Username, group)
	s.ChannelRepository.Close(req.Username)

	return &chat.LogoutResponse{}, nil
}

func (s *Server) Stream(stream chat.Chat_StreamServer) error {
	token, err := s.GetTokenFromContext(stream.Context())

	if err != nil || !s.Authenticator.IsTokenValid(token) {
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
		if err == io.EOF || err != nil {
			break
		}

		username, group, message := req.Username, req.Group, req.Message

		token, err := s.GetTokenFromContext(stream.Context())
		if err != nil || !s.Authenticator.AuthenticateToken(token, req.Username) {
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

		groupMembers, _ := s.GroupRepository.GetGroupMembers(group, username)
		for _, reciever_username := range groupMembers {
			if username != reciever_username {
				clientChannel, _ := s.ChannelRepository.Get(reciever_username)
				clientChannel <- res
			}
		}

	}
}

func sender(s *Server, stream chat.Chat_StreamServer, wg *sync.WaitGroup, token string) {
	defer wg.Done()

	username, _ := s.UserRepository.GetUsername(token)
	clientChannel, _ := s.ChannelRepository.Get(username)

	for res := range clientChannel {
		stream.Send(&res)
	}
}

func (s *Server) GetTokenFromContext(ctx context.Context) (string, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if value, found := meta["token"]; !found || len(value) == 0 || !ok {
		return "", fmt.Errorf("Missing Token in Request")
	}

	return meta["token"][0], nil
}
