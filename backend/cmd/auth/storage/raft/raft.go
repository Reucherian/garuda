package raft

import (
	"context"
	"errors"

	"garuda.com/m/model"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

//  the way we have to think about this now is one giant key value store so let us say that each user name is the key and the

type Raft struct {
	cli *clientv3.Client
}

func CreateNewRaft() *Raft {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   Endpoints,
		DialTimeout: DialTimeout,
	})
	if err != nil {
		panic(err)
	}
	// context with timeout
	return &Raft{cli: cli}
}

func (r *Raft) AddUser(username string, hashedPassword string) error {
	if _, err := r.GetUserHelper(username); err == nil {
		return errors.New("user already exists")
	} else if err != nil && err.Error() != "user not found" {
		return err
	}
	user := &model.UserStg{
		Username: username, HashPassword: hashedPassword,
	}
	userData, err := proto.Marshal(user)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	_, err = r.cli.KV.Put(ctx, username, string(userData))
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (r *Raft) GetUser(username string) (model.UserStg, error) {
	value, err := r.GetUserHelper(username)
	if err != nil {
		return model.UserStg{}, err
	}
	user := model.UserStg{}
	err = proto.Unmarshal(value, &user)
	if err != nil {
		return model.UserStg{}, err
	}
	return user, nil
}

func (r *Raft) UpdateUser(username, hash_password string) error {
	_, err := r.GetUserHelper(username)
	if err != nil {
		return err
	}
	user := &model.UserStg{
		Username: username, HashPassword: hash_password,
	}
	userData, err := proto.Marshal(user)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	_, err = r.cli.KV.Put(ctx, username, string(userData))
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (r *Raft) DeleteUser(username string) error {
	if _, err := r.GetUserHelper(username); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	_, err := r.cli.Delete(ctx, username)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (r *Raft) CreatePost(username string, title string, content string) error {
	if _, err := r.GetUserHelper(username); err != nil {
		return err
	}
	user, _ := r.GetUser(username)
	user.Posts = append(user.Posts, &model.PostStg{Title: title, Content: content})
	userData, err := proto.Marshal(&user)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	_, err = r.cli.KV.Put(ctx, username, string(userData))
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (r *Raft) GetPosts(username string) ([]*model.PostStg, error) {
	if _, err := r.GetUserHelper(username); err != nil {
		return nil, err
	}
	user, err := r.GetUser(username)
	if err != nil {
		return nil, err
	}
	return user.GetPosts(), nil
}

func (r *Raft) DeletePost(username string, title string) error {
	if _, err := r.GetUserHelper(username); err != nil {
		return err
	}
	user, err := r.GetUser(username)
	if err != nil {
		return err
	}
	for i, post := range user.GetPosts() {
		if post.GetTitle() == title {
			user.Posts = append(user.Posts[:i], user.Posts[i+1:]...)
			userData, err := proto.Marshal(&user)
			if err != nil {
				return err
			}
			r.cli.KV.Put(context.Background(), username, string(userData))
			return nil
		}
	}
	return errors.New("post not found")
}

func (r *Raft) UpdatePost(username string, title string, content string) error {
	if _, err := r.GetUserHelper(username); err != nil {
		return err
	}
	user, err := r.GetUser(username)
	if err != nil {
		return err
	}
	for _, post := range user.GetPosts() {
		if post.GetTitle() == title {
			post.Content = content
			userData, err := proto.Marshal(&user)
			if err != nil {
				return err
			}
			r.cli.KV.Put(context.Background(), username, string(userData))
			return nil
		}
	}
	return errors.New("post not found")
}

func (r *Raft) AddFollowing(follower, following string) error {
	if _, err := r.GetUserHelper(follower); err != nil {
		return err
	}
	if _, err := r.GetUserHelper(following); err != nil {
		return err
	}
	if follower == following {
		return errors.New("cannot follow yourself")
	}
	user, err := r.GetUser(follower)
	if err != nil {
		return err
	}
	followingMap := user.GetFollowing()
	if followingMap == nil {
		followingMap = make(map[string]int32)
	}
	if _, ok := followingMap[following]; ok {
		return errors.New("already following")
	}
	followingMap[following] = 1
	user.Following = followingMap
	userData, err := proto.Marshal(&user)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	_, err = r.cli.KV.Put(ctx, follower, string(userData))
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (r *Raft) GetFollowings(username string) ([]string, error) {
	if _, err := r.GetUserHelper(username); err != nil {
		return nil, err
	}
	user, err := r.GetUser(username)
	if err != nil {
		return nil, err
	}
	followings := make([]string, 0)
	followingMap := user.GetFollowing()
	for user := range followingMap {
		followings = append(followings, user)
	}
	return followings, nil
}

func (r *Raft) DeleteFollowing(follower, following string) error {
	if _, err := r.GetUserHelper(follower); err != nil {
		return err
	}
	if _, err := r.GetUserHelper(following); err != nil {
		return err
	}
	user, err := r.GetUser(follower)
	if err != nil {
		return err
	}
	followingMap := user.GetFollowing()
	if _, ok := followingMap[following]; !ok {
		return errors.New("following not found")
	}
	delete(followingMap, following)
	user.Following = followingMap
	userData, err := proto.Marshal(&user)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	_, err = r.cli.KV.Put(ctx, follower, string(userData))
	cancel()
	if err != nil {
		return err
	}
	return nil
}
