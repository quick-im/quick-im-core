package poll

import "net/http"

type PollProtoc struct {
}

func (p *PollProtoc) Handler(_ http.ResponseWriter, _ *http.Request) {
	panic("not implemented") // TODO: Implement
}
