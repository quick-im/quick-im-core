package ws

import "net/http"

type WSProtoc struct {
}

func (w *WSProtoc) Handler(_ http.ResponseWriter, _ *http.Request) {
	panic("not implemented") // TODO: Implement
}
