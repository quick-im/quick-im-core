package contant

const (
	CTX_POSTGRES_KEY        ContentKey = "__ctx.database.postgres.key__"
	CTX_RETHINK_DB_KEY      ContentKey = "__ctx.database.rethinkdb.key__"
	CTX_REDIS_KEY           ContentKey = "__ctx.cache.redis.key__"
	CTX_NATS_KEY            ContentKey = "__ctx.msg.nats.key__"
	CTX_SERVICE_MSGBORKER   ContentKey = "__ctx.service.msgborker.key__"
	CTX_SERVICE_PERSISTENCE ContentKey = "__ctx.service.persistence.key__"
)
