package contant

const (
	CTX_POSTGRES_KEY         ContentKey = "__ctx.database.postgres.key__"
	CTX_RETHINK_DB_KEY       ContentKey = "__ctx.database.rethinkdb.key__"
	CTX_REDIS_KEY            ContentKey = "__ctx.cache.redis.key__"
	CTX_CACHE_DB_KEY         ContentKey = "__ctx.cache.db.key__"
	CTX_NATS_KEY             ContentKey = "__ctx.msg.nats.key__"
	CTX_SERVICE_MSGBORKER    ContentKey = "__ctx.service.msgborker.key__"
	CTX_SERVICE_PERSISTENCE  ContentKey = "__ctx.service.persistence.key__"
	CTX_SERVICE_CONVERSATION ContentKey = "__ctx.service.conversation.key__"
	CTX_SERVICE_MSGHUB       ContentKey = "__ctx.service.msghub.key__"
	CTX_SERVICE_MSGID        ContentKey = "__ctx.service.msgid.key__"
)

const (
	HTTP_CTX_JWT_CLAIMS ContentKey = "__http.ctx.jwt.claims__"
)
