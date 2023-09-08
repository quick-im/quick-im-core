package quickerr

import "errors"

var (
	ErrConversationTypeRange   = errors.New("conversation id ranges from 0 to 15")
	ErrConversationNumberRange = errors.New("conversation contains at least one user")
	ErrTraceClosed             = errors.New("trace closed")
	ErrDriveNotSupport         = errors.New("unsupported driver")
	ErrToken                   = warpErr(10101, "invalid token")
	ErrHttpInvaildParam        = warpErr(10001, "invalid parameter")
)

type errWarp struct {
	Code   int32  `json:"code"`
	ErrStr string `json:"error"`
	Err    error  `json:"-"`
}

func warpErr(code int32, err string) errWarp {
	return errWarp{Err: errors.New(err), ErrStr: err, Code: code}
}

func (e errWarp) Error() string {
	return e.Err.Error()
}

func (e errWarp) GetCode() int32 {
	return e.Code
}
