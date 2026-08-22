package commands

type Session struct {
	UserID   int
	loggedIn bool
}

func NewSession(userID int) *Session {
	return &Session{}
}
