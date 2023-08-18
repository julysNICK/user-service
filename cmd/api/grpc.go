package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"user/data"
	user "user/users"

	"google.golang.org/grpc"
)

type UserServer struct {
	user.UnimplementedUserServiceServer
	Models data.Models
}

func (s *UserServer) DeleteUser(ctx context.Context, req *user.UserRequestDelete) (*user.UserResponseDelete, error) {
	_, err := strconv.Atoi(req.Email)
	if err != nil {
		return nil, err
	}
	err = s.Models.User.Delete()
	if err != nil {
		return nil, err
	}

	return &user.UserResponseDelete{
		Email: req.Email,
	}, nil
}

func (app *Config) ListenGrpc() {
	lis, er := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))

	if er != nil {
		log.Fatal("Listener error", er)
	}

	s := grpc.NewServer()

	user.RegisterUserServiceServer(s, &UserServer{
		Models: app.Models,
	})

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
