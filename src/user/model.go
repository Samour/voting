package user

type User struct {
	UserId      string
	DisplayName string
}

type UsernamePasswordCredential struct {
	Username     string
	PasswordHash string
	UserId       string
}
