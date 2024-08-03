package signup

type signUpModel struct {
	Username         string
	ErrorMessage     string
	ValidationErrors []string
	LogInUrl         string
}
