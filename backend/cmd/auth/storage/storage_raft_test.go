package storage

// important to start raft endpoints before running tests

import (
	"fmt"
	"testing"
)

// Test to check if adding a single user to the storage of raft type is working
func TestCreateUsersRaft(t *testing.T) {
	storage, err := CreateNewStorage("raft")
	testUser := "testuser"
	testPassword := "password"
	if err != nil {
		t.Error(err)
	}
	if storage == nil {
		t.Error("Storage is nil")
	}
	err = storage.AddUser(testUser, testPassword)
	if err != nil {
		t.Error(err)
	}

	user, err := storage.GetUser(testUser)
	if err != nil {
		t.Error(err)
	}
	if user.Username != testUser {
		t.Error("User with username " + testUser + " not found in memory storage")
	}
	go func() {
		storage, err := CreateNewStorage("raft")
		testUser := "testuserparallel"
		testPassword := "password"
		if err != nil {
			t.Error(err)
		}
		if storage == nil {
			t.Error("Storage is nil")
		}
		err = storage.AddUser(testUser, testPassword)
		if err != nil {
			t.Error(err)
		}
		user, err := storage.GetUser(testUser)
		if err != nil {
			t.Error(err)
		}
		if user.Username != testUser {
			t.Error("User with username " + testUser + " not found in memory storage")
		}
	}()
}

// Test to UpdateUser to the storage of raft type
func TestUpdateUserRaft(t *testing.T) {
	storage, err := CreateNewStorage("raft")
	testUser := "testuser"
	if err != nil {
		t.Error(err)
	}
	if storage == nil {
		t.Error("Storage is nil")
	}
	prevUser, err := storage.GetUser(testUser)
	if err != nil {
		t.Error(err)
	}
	err = storage.UpdateUser(testUser, "newpassword")
	if err != nil {
		t.Error(err)
	}
	currUser, err := storage.GetUser(testUser)
	if err != nil {
		t.Error(err)
	}
	if prevUser.HashPassword == currUser.HashPassword {
		t.Error("User with username " + testUser + " not updated in raft storage")
	}
	go func() {
		storage, err := CreateNewStorage("raft")
		testUser := "testuserparallel"
		if err != nil {
			t.Error(err)
		}
		if storage == nil {
			t.Error("Storage is nil")
		}
		prevUser, err := storage.GetUser(testUser)
		if err != nil {
			t.Error(err)
		}
		err = storage.UpdateUser(testUser, "newpassword")
		if err != nil {
			t.Error(err)
		}
		currUser, err := storage.GetUser(testUser)
		if err != nil {
			t.Error(err)
		}
		if prevUser.HashPassword == currUser.HashPassword {
			t.Error("User with username " + testUser + " not updated in raft storage")
		}
	}()
}

// Test to CreatePost to the storage of raft type
func TestCreatePostRaft(t *testing.T) {
	storage, err := CreateNewStorage("raft")
	if err != nil {
		t.Error(err)
	}
	testUser := "testuser"
	err = storage.CreatePost(testUser, "title", "content")
	if err != nil {
		t.Error(err)
	}
	posts, err := storage.GetPosts(testUser)
	if err != nil {
		t.Error(err)
	}
	if len(posts) != 1 {
		t.Error("Post not created")
	}
	if len(posts) > 0 {
		if posts[0].Title != "title" {
			t.Error("Post title not created")
		}
		if posts[0].Content != "content" {
			t.Error("Post content not created")
		}
	}
}

// Test to UpdatePost to the storage of raft type
func TestUpdatePostRaft(t *testing.T) {
	user := "testuser"
	title := "newtitle"
	content := "content"
	storage, err := CreateNewStorage("raft")

	if err != nil {
		t.Error(err)
	}
	if storage == nil {
		t.Error("Storage is nil")
	}

	err = storage.CreatePost(user, title, content)
	if err != nil {
		t.Error(err)
	}

	err = storage.UpdatePost(user, title, "newcontent")
	if err != nil {
		t.Error(err)
	}

	posts, err := storage.GetPosts(user)
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 2 {
		t.Error("Post not created")
	}

	if len(posts) > 1 {
		if posts[1].Title != title {
			t.Error("Post title not created")
		}

		if posts[1].Content != "newcontent" {
			t.Error("Post content not created")
		}
	}
}

// Test to DeletePost to the storage of raft type
func TestDeletePostRaft(t *testing.T) {
	user := "testuser"
	title := "title"
	title2 := "newtitle"
	storage, err := CreateNewStorage("raft")

	if err != nil {
		t.Error(err)
	}

	if storage == nil {
		t.Error("Storage is nil")
	}

	err = storage.DeletePost(user, title)
	if err != nil {
		t.Error(err)
	}

	posts, err := storage.GetPosts(user)
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 1 {
		t.Error("Post not deleted")
	}

	err = storage.DeletePost(user, title2)
	if err != nil {
		t.Error(err)
	}

	posts, err = storage.GetPosts(user)
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 0 {
		t.Error("Post not deleted")
	}
}

// Test to AddFollowing to the storage of raft type
func TestAddFollowingsRaft(t *testing.T) {
	user := "testuser"
	password := "password"
	following := "following"

	storage, err := CreateNewStorage("raft")
	if err != nil {
		t.Error(err)
	}
	if storage == nil {
		t.Error("Storage is nil")
	}

	err = storage.AddUser(following, password)
	if err != nil {
		t.Error(err)
	}

	err = storage.AddFollowing(user, following)
	if err != nil {
		t.Error(err)
	}

	followings, err := storage.GetFollowings(user)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(followings)
	if len(followings) != 1 {
		t.Error("Following not added")
	}

	if len(followings) > 0 {
		if followings[0] != following {
			t.Error("Following not added")
		}
	}
}

// Test to DeleteFollowings to the source of raft type
func TestDeleteFollowingRaft(t *testing.T) {
	user := "testuser"
	follower := "following"

	storage, err := CreateNewStorage("raft")
	if err != nil {
		t.Error(err)
	}
	if storage == nil {
		t.Error("Storage is nil")
	}

	err = storage.DeleteFollowing(user, follower)
	if err != nil {
		t.Error(err)
	}

	followings, err := storage.GetFollowings(user)
	if err != nil {
		t.Error(err)
	}
	if len(followings) != 0 {
		t.Error("Following not deleted")
	}
}

// Test to DeleteUser to the storage of raft type
func TestDeleteUserRaft(t *testing.T) {
	storage, err := CreateNewStorage("raft")
	testUser := "testuser"
	following := "following"
	testUserP := "testuserparallel"
	if err != nil {
		t.Error(err)
	}
	if storage == nil {
		t.Error("Storage is nil")
	}
	err = storage.DeleteUser(testUser)
	if err != nil {
		t.Error(err)
	}
	_, err = storage.GetUser(testUser)
	if err == nil {
		t.Error("User with username " + testUser + " not deleted in memory storage")
	}
	err = storage.DeleteUser(following)
	if err != nil {
		t.Error(err)
	}
	_, err = storage.GetUser(following)
	if err == nil {
		t.Error("User with username " + testUser + " not deleted in memory storage")
	}
	err = storage.DeleteUser(testUserP)
	if err != nil {
		t.Error(err)
	}
	_, err = storage.GetUser(testUserP)
	if err == nil {
		t.Error("User with username " + testUser + " not deleted in memory storage")
	}
}
