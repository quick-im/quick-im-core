package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/quick-im/quick-im-core/application/gateway/gateway/access"
	"github.com/quick-im/quick-im-core/application/gateway/gateway/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var Data = []byte{113, 117, 105, 99, 107, 45, 105, 109}

func (a *apiServer) InitAndStartServer(ctx context.Context) {
	router := httprouter.New()
	// cert, err := tls.LoadX509KeyPair(config.PublicCert, config.PriviteCert)
	// if err != nil {
	// 	panic(err)
	// }

	// if err := server.ListenAndServeTLS("", ""); err != nil {
	// 	panic(err)
	// }
	// http.ListenAndServe(":8088", router)
	// server := http.Server{
	// 	Addr:    ":8088",
	// 	Handler: router,
	// 	TLSConfig: &tls.Config{
	// 		NextProtos:   []string{"h2", "http/1.1"},
	// 		Certificates: []tls.Certificate{cert},
	// 	},
	// 	ErrorLog: log.New(nil, "", 0),
	// }
	h2s := http2.Server{}
	h1 := http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.ip, a.port),
		Handler: h2c.NewHandler(router, &h2s),
	}
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
	router.HandlerFunc("GET", "/notify", middleware.JwtAuth(ctx, access.NotifyHandler))
	router.HandlerFunc("GET", "/", middleware.AllowCros(ctx, func(w http.ResponseWriter, r *http.Request) {
		w.Write(Data)
	}))
	if err := h1.ListenAndServe(); err != nil {
		panic(err)
	}
}
