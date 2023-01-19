package auth

import (
	"errors"

	"garuda.com/m/web/cmd/auth/storage"
	"golang.org/x/crypto/bcrypt"
)

func Register(username, password string) error {
	// TODO: check if password is hashed successfully
	hash_password, _ := HashPassword(password)
	return storage.OurStore.AddUser(username, hash_password)

}

func Login(username, password string) error {
	user, ok := storage.OurStore.GetUser(username)
	if ok != nil {
		return errors.New("user does not exist")
	}
	if CheckPasswordHash(password, user.HashPassword) {
		// expirationTime := time.Now().Add(5 * time.Minute)
		// claims := &Claims{
		// 	Username: username,
		// 	RegisteredClaims: jwt.RegisteredClaims{
		// 		ExpiresAt: jwt.NewNumericDate(expirationTime),
		// 	},
		// }
		// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// // TODO: take the value from environment variable instead of using config file
		// tokenString, err := token.SignedString([]byte(AuthConfig.SecretKey))
		// if err != nil {
		// 	return "", time.Time{}, err
		// }
		// return tokenString, expirationTime, nil
		// TODO: move the above code to the client side
		return nil

	}
	return errors.New("wrong password")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
