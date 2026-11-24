package cli

type Session struct {
	UserID   int64
	loggedIn bool
}

func NewSession(userID int64) *Session {
	return &Session{UserID: userID}
}
