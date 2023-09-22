package quickerr

import "errors"

var (
	ErrConversationTypeRange     = warpErr(90010, "conversation id ranges from 0 to 15")
	ErrConversationNumberRange   = warpErr(90011, "conversation contains at least one user")
	ErrTraceClosed               = warpErr(90012, "trace closed")
	ErrDriveNotSupport           = warpErr(90013, "unsupported driver")
	ErrToken                     = warpErr(10101, "invalid token")
	ErrHttpInvaildParam          = warpErr(10001, "invalid parameter")
	ErrInternalServiceCallFailed = warpErr(10002, "internal service call failed")
	ErrNotAllowedRequest         = warpErr(10003, "not allowed request")
)

type responseWarp[T any] struct {
	Code   int32  `json:"code"`
	ErrStr string `json:"error,omitempty"`
	Err    error  `json:"-"`
	Data   T      `json:"data,omitempty"`
}

func warpErr(code int32, err string) responseWarp[struct{}] {
	return responseWarp[struct{}]{Err: errors.New(err), ErrStr: err, Code: code}
}

func (e responseWarp[T]) Error() string {
	return e.Err.Error()
}

func (e responseWarp[T]) GetCode() int32 {
	return e.Code
}

func HttpResponeWarp[T any](data T) responseWarp[T] {
	return responseWarp[T]{
		Code: 0,
		Data: data,
	}
}
