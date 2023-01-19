package auth

import (
	"testing"

	"garuda.com/m/web/cmd/auth/storage"
)

// Test to check if registering a user is working
func TestRegister(t *testing.T) {
	user := "testuser"
	password := "password"

	err := Register(user, password)
	if err != nil {
		t.Error(err)
	}
}

// Test to check if logging in a user is working
func TestLogin(t *testing.T) {
	user := "testuser"
	password := "password"

	err := Login(user, password)
	if err != nil {
		t.Error(err)
	}
	err = Login(user, "wrongpassword")
	if err == nil {
		t.Error("Login with wrong password should fail")
	}

	err = Login("wronguser", password)
	if err == nil {
		t.Error("Login with wrong user should fail")
	}

	err = Login("wronguser", "wrongpassword")
	if err == nil {
		t.Error("Login with wrong user and wrong password should fail")
	}
}

// Test to if password hashing is working
func TestHashPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}

	if hashedPassword == "" {
		t.Error("Hashed password is empty")
	}
}

// Test to check if checking password is working
func TestCheckPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}

	if hashedPassword == "" {
		t.Error("Hashed password is empty")
	}

	ok := CheckPasswordHash(password, hashedPassword)
	if !ok {
		t.Error(err)
	}

	ok = CheckPasswordHash("wrongpassword", hashedPassword)
	if ok {
		t.Error("CheckPassword with wrong password should fail")
	}
}

func TestEnder(t *testing.T) {
	// Test to make the storage empty
	storage, err := storage.CreateNewStorage("raft")
	if err != nil {
		t.Error(err)
	}
	if storage == nil {
		t.Error("Storage is nil")
	}
	storage.DeleteUser("testuser")
}
