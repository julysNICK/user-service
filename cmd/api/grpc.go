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
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *UserServer) GetAllUsers(context.Context, *emptypb.Empty) (*user.UserResponseGetAll, error) {
	users, err := s.Models.User.GetAll()
	if err != nil {
		return nil, err
	}

	var userResponses []*user.User

	for _, userDatabase := range users {
		userResponses = append(userResponses, &user.User{
			Email:     userDatabase.Email,
			Password:  userDatabase.Password,
			FirstName: userDatabase.FirstName,
		})
	}

	return &user.UserResponseGetAll{
		Users: userResponses,
	}, nil

}

func (s *UserServer) GetOneUser(c context.Context, req *user.UserRequestGetOne) (*user.UserResponseGetOne, error) {
	userDatabase, err := s.Models.User.GetOne(int(req.Id))
	if err != nil {
		return nil, err
	}

	return &user.UserResponseGetOne{
		Email:     userDatabase.Email,
		FirstName: userDatabase.FirstName,
		Message:   "User found",
	}, nil

}

func (s *UserServer) UpdateUser(ctx context.Context, req *user.UserRequestUpdate) (*user.UserResponseUpdate, error) {
	_, err := strconv.Atoi(req.Email)
	if err != nil {
		return nil, err
	}
	err = s.Models.User.Update()
	if err != nil {
		return nil, err
	}

	return &user.UserResponseUpdate{
		Email:   req.Email,
		Message: "User updated",
	}, nil
}
