package main

import (
	"context"
	"net"

	"garuda.com/m/backend/cmd/activity"
	"garuda.com/m/backend/cmd/auth"
	"garuda.com/m/backend/config"
	"garuda.com/m/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type activityServer struct {
	model.UnimplementedActivityServiceServer
}
type authServer struct {
	model.UnimplementedAuthServiceServer
}

func (s *authServer) Register(ctx context.Context, in *model.UserRequest) (*model.Empty, error) {
	username, password := in.GetUsername(), in.GetPassword()
	return &model.Empty{}, auth.Register(username, password)
}

func (s *authServer) Login(ctx context.Context, in *model.UserRequest) (*model.Empty, error) {
	username, password := in.GetUsername(), in.GetPassword()
	return &model.Empty{}, auth.Login(username, password)
}

func (s *activityServer) CreatePost(ctx context.Context, in *model.PostRequest) (*model.Empty, error) {
	user, post := in.GetUser(), in.GetPost()
	return &model.Empty{}, activity.CreatePost(user.Username, post.Title, post.Content)
}

// TODO: Think about optimizations for this code
func (s *activityServer) GetPosts(ctx context.Context, in *model.User) (*model.Posts, error) {
	username := in.GetUsername()
	finalPost := []*model.Post{}
	posts, err := activity.GetPosts(username)
	if err != nil {
		return &model.Posts{Posts: finalPost}, err
	}
	for _, post := range posts {
		finalPost = append(finalPost, &model.Post{Title: post["title"], Content: post["content"]})
	}
	return &model.Posts{Posts: finalPost}, nil
}

func (s *activityServer) AddFollowing(ctx context.Context, in *model.FollowingRequest) (*model.Empty, error) {
	user, following := in.GetUser(), in.GetFollowing()
	return &model.Empty{}, activity.AddFollowing(user.Username, following.Username)
}

// TODO: Think about optimizations for this code
func (s *activityServer) GetFollowings(ctx context.Context, in *model.User) (*model.Users, error) {
	username := in.GetUsername()
	followings, err := activity.GetFollowings(username)
	if err != nil {
		return nil, err
	}
	finalFollowings := []*model.User{}
	for _, following := range followings {
		finalFollowings = append(finalFollowings, &model.User{Username: following})
	}
	return &model.Users{Users: finalFollowings}, nil
}

func (s *activityServer) DeleteFollowing(ctx context.Context, in *model.FollowingRequest) (*model.Empty, error) {
	user, following := in.GetUser(), in.GetFollowing()
	return &model.Empty{}, activity.DeleteFollowing(user.Username, following.Username)
}

func main() {
	go func() {
		activityListener, err := net.Listen("tcp", config.ActivityServerPort)
		if err != nil {
			panic(err)
		}
		activityGrpcServer := grpc.NewServer()
		model.RegisterActivityServiceServer(activityGrpcServer, &activityServer{})
		reflection.Register(activityGrpcServer)
		if err := activityGrpcServer.Serve(activityListener); err != nil {
			panic(err)
		}
	}()
	authListener, err := net.Listen("tcp", config.AuthServerPort)
	if err != nil {
		panic(err)
	}
	authGrpcServer := grpc.NewServer()
	model.RegisterAuthServiceServer(authGrpcServer, &authServer{})
	reflection.Register(authGrpcServer)
	if err := authGrpcServer.Serve(authListener); err != nil {
		panic(err)
	}

}
