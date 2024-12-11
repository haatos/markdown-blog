package internal

type RoleID int

func (r RoleID) ToString() string {
	switch r {
	case Superuser:
		return "superuser"
	case Member:
		return "member"
	default:
		return "any user"
	}
}

const (
	Member    RoleID = 1
	Superuser RoleID = 10000
)

const (
	CookieExpiresLayout = "2006-01-02T15:04:05Z"
	DBTimestampLayout   = "2006-01-02 15:04:05"
	SessionCookie       = "session"
	OAuthCookie         = "oauthstate"
	PasswordlessCookie  = "passwordless"
)

const (
	CommentsPerPage = 10
)
