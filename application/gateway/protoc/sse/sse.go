package sse

import "net/http"

type SSEProtoc struct {
}

func (s *SSEProtoc) Handler(_ http.ResponseWriter, _ *http.Request) {
	panic("not implemented") // TODO: Implement
}
