package handler

type SignUpForm struct {
	GoogleOAuthHref string
	GithubOAuthHref string
	FirstName       string
	FirstNameError  string
	LastName        string
	LastNameError   string
	Email           string
	EmailError      string
}
