package memory

import (
	"errors"
	"sync"

	"garuda.com/m/model"
)

type Memory struct {
	users sync.Map
}

func CreateNewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) AddUser(username string, hashedPassword string) error {
	if _, ok := m.users.Load(username); ok {
		return errors.New("user already exists")
	}
	m.users.Store(username, model.UserStg{Username: username, HashPassword: hashedPassword})
	return nil
}

func (m *Memory) GetUser(username string) (model.UserStg, error) {
	if _, ok := m.users.Load(username); !ok {
		return model.UserStg{}, errors.New("user does not exist")
	}
	user, _ := m.users.Load(username)
	return user.(model.UserStg), nil
}

func (m *Memory) UpdateUser(username, hash_password string) error {
	if _, ok := m.users.Load(username); !ok {
		return errors.New("user does not exist")
	}
	m.users.Store(username, model.UserStg{Username: username, HashPassword: hash_password})
	return nil
}

func (m *Memory) DeleteUser(username string) error {
	if _, ok := m.users.Load(username); !ok {
		return errors.New("user does not exist")
	}
	m.users.Delete(username)
	return nil
}

func (m *Memory) CreatePost(username string, title string, content string) error {
	if _, ok := m.users.Load(username); !ok {
		return errors.New("user does not exist")
	}
	user, _ := m.GetUser(username)
	user.Posts = append(user.Posts, &model.PostStg{Title: title, Content: content})
	m.users.Store(username, user)
	return nil
}

func (m *Memory) GetPosts(username string) ([]*model.PostStg, error) {
	if _, ok := m.users.Load(username); !ok {
		return nil, errors.New("user does not exist")
	}
	user, _ := m.GetUser(username)
	return user.Posts, nil
}

func (m *Memory) DeletePost(username string, title string) error {
	if _, ok := m.users.Load(username); !ok {
		return errors.New("user does not exist")
	}
	user, _ := m.GetUser(username)
	for i, post := range user.Posts {
		if post.Title == title {
			user.Posts = append(user.Posts[:i], user.Posts[i+1:]...)
			m.users.Store(username, user)
			return nil
		}
	}
	return errors.New("post not found")
}

func (m *Memory) UpdatePost(username string, title string, content string) error {
	if _, ok := m.users.Load(username); !ok {
		return errors.New("user does not exist")
	}
	user, _ := m.GetUser(username)
	for i, post := range user.Posts {
		if post.Title == title {
			user.Posts[i].Content = content
			m.users.Store(username, user)
			return nil
		}
	}
	return errors.New("post not found")
}

func (m *Memory) AddFollowing(follower, following string) error {
	if _, ok := m.users.Load(follower); !ok {
		return errors.New("user does not exist")
	}
	if _, ok := m.users.Load(following); !ok {
		return errors.New("user does not exist")
	}
	if follower == following {
		return errors.New("cannot follow yourself")
	}
	user, _ := m.GetUser(follower)
	followingMap := user.GetFollowing()
	if followingMap == nil {
		followingMap = make(map[string]int32)
	}
	if _, ok := followingMap[following]; ok {
		return errors.New("already following")
	}
	followingMap[following] = 1
	user.Following = followingMap
	m.users.Store(follower, user)
	return nil
}

func (m *Memory) GetFollowings(username string) ([]string, error) {
	if _, ok := m.users.Load(username); !ok {
		return nil, errors.New("user does not exist")
	}
	user, _ := m.GetUser(username)
	followings := make([]string, 0)
	followingMap := user.GetFollowing()
	for user := range followingMap {
		followings = append(followings, user)
	}
	return followings, nil
}

func (m *Memory) DeleteFollowing(follower, following string) error {
	if _, ok := m.users.Load(follower); !ok {
		return errors.New("user does not exist")
	}
	if _, ok := m.users.Load(following); !ok {
		return errors.New("user does not exist")
	}
	if follower == following {
		return errors.New("cannot unfollow yourself")
	}
	user, _ := m.GetUser(follower)
	followingMap := user.GetFollowing()
	if _, ok := followingMap[following]; !ok {
		return errors.New("following not found")
	}
	delete(followingMap, following)
	user.Following = followingMap
	return nil
}
