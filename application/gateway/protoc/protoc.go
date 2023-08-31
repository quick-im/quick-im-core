package protoc

import "net/http"

type ProtocHandler interface {
	Handler(http.ResponseWriter, *http.Request)
}
