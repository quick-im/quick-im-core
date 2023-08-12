package quickerr

import "errors"

var (
	ErrConversationTypeRange   = errors.New("conversation id ranges from 0 to 15")
	ErrConversationNumberRange = errors.New("conversation contains at least one user")
	ErrTraceClosed             = errors.New("trace closed")
)
