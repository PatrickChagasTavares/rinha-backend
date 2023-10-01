package people

import "errors"

var (
	ErrNicknameAlreadyUsed = errors.New("nickname used")
	ErrPersonNotFound      = errors.New("person not found")
)
