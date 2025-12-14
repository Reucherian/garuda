package storage

import (
	"errors"

	"garuda.com/m/backend/cmd/auth/storage/memory"
	"garuda.com/m/backend/cmd/auth/storage/raft"
	"garuda.com/m/model"
)

type Storage interface {
	AddUser(string, string) error
	GetUser(string) (model.UserStg, error)
	UpdateUser(string, string) error
	DeleteUser(username string) error
	CreatePost(username string, title string, content string) error
	GetPosts(username string) ([]*model.PostStg, error)
	UpdatePost(username string, title string, content string) error
	DeletePost(username string, title string) error
	AddFollowing(follower, following string) error
	GetFollowings(username string) ([]string, error)
	DeleteFollowing(follower, following string) error
}

func CreateNewStorage(storageType string) (Storage, error) {
	if storageType == "memory" {
		// return nil, errors.New("not implemented")
		return memory.CreateNewMemory(), nil
	}
	if storageType == "raft" {
		return raft.CreateNewRaft(), nil
		// return nil, errors.New("Not implemented")
	}
	return nil, errors.New("unidentified storage type")
}

var OurStore, _ = CreateNewStorage(SelectedStorage)
