package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/quick-im/quick-im-core/application/gateway/gateway/access"
	"github.com/quick-im/quick-im-core/application/gateway/gateway/middleware"
	"github.com/quick-im/quick-im-core/application/gateway/protoc"
	"github.com/quick-im/quick-im-core/application/gateway/protoc/poll"
	"github.com/quick-im/quick-im-core/application/gateway/protoc/sse"
	"github.com/quick-im/quick-im-core/application/gateway/protoc/ws"
	"github.com/quick-im/quick-im-core/internal/contant"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var Data = []byte{113, 117, 105, 99, 107, 45, 105, 109}

func (a *apiServer) InitAndStartServer(ctx context.Context) {
	ctx = context.WithValue(ctx, contant.CTX_LOGGER_KEY, a.logger)
	router := httprouter.New()
	// 注册长连接支持的多种协议
	protoc.RegisterDrive("ws", ws.InitProtoc())
	protoc.RegisterDrive("sse", sse.InitProtoc())
	protoc.RegisterDrive("poll", poll.InitProtoc())
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
	// 公开接口
	// conversation部分
	router.HandlerFunc("GET", "/notify", middleware.JwtAuth(ctx, access.NotifyHandler))
	router.HandlerFunc("POST", "/send_msg", middleware.JwtAuth(ctx, access.SendMsgHandler))
	router.HandlerFunc("POST", "/check_joined_conversation", middleware.JwtAuth(ctx, access.CheckJoinedConversation))
	router.HandlerFunc("GET", "/get_joined_conversation", middleware.JwtAuth(ctx, access.GetJoinedConversation))
	router.HandlerFunc("POST", "/get_conversation_detail", middleware.JwtAuth(ctx, access.GetConversationInfo))
	router.HandlerFunc("POST", "/join_conversation", middleware.JwtAuth(ctx, access.JoinConversation))
	router.HandlerFunc("POST", "/create_conversation", middleware.JwtAuth(ctx, access.CreateConversation))
	router.HandlerFunc("POST", "/leave_conversation", middleware.JwtAuth(ctx, access.LeaveConversation))
	router.HandlerFunc("POST", "/create_conversation", middleware.JwtAuth(ctx, access.CreateConversation))
	// 受保护的接口
	//conversation部分
	router.HandlerFunc("POST", "/create_conversation", middleware.ProtectApi(ctx, access.CreateConversationInner))
	router.HandlerFunc("POST", "/kickout_conversation", middleware.ProtectApi(ctx, access.KickoutConversationInner))
	router.HandlerFunc("POST", "/join_conversation", middleware.ProtectApi(ctx, access.JoinConversationInner))
	router.HandlerFunc("GET", "/", middleware.AllowCros(ctx, func(w http.ResponseWriter, r *http.Request) {
		w.Write(Data)
	}))
	if err := h1.ListenAndServe(); err != nil {
		panic(err)
	}
}
