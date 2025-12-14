package activity

import (
	"garuda.com/m/backend/cmd/auth/storage"
)

func CreatePost(username, title, content string) error {
	return storage.OurStore.CreatePost(username, title, content)
}

func GetPosts(username string) ([]map[string]string, error) {
	finalPost := []map[string]string{}
	posts, err := storage.OurStore.GetPosts(username)
	if err != nil {
		return finalPost, err
	}
	for _, post := range posts {
		finalPost = append(finalPost, map[string]string{"title": post.Title, "content": post.Content})
	}
	return finalPost, nil
}

func AddFollowing(username, following string) error {
	return storage.OurStore.AddFollowing(username, following)
}

func GetFollowings(username string) ([]string, error) {
	return storage.OurStore.GetFollowings(username)
}

func DeleteFollowing(username, following string) error {
	return storage.OurStore.DeleteFollowing(username, following)
}
