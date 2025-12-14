package auth

// this will have token validity duration and will include hash difficulties etc.
type Config struct {
	SecretKey string
}

var AuthConfig = Config{}
