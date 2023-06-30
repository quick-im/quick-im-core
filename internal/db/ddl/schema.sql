-- public.conversations definition

-- Drop table

-- DROP TABLE public.conversations;

CREATE TABLE IF NOT EXISTS public.conversations (
	conversation_id uuid NOT NULL,
	last_msg_id varchar NULL,
	last_send_time varchar NOT NULL,
	is_delete bool NOT NULL DEFAULT false,
	conversation_type int4 NOT NULL DEFAULT 0,
	last_send_session varchar NOT NULL,
	CONSTRAINT conversations_pk PRIMARY KEY (conversation_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS conversations_conversation_id_idx ON public.conversations USING btree (conversation_id);

-- public.conversation_session_id definition

-- Drop table

-- DROP TABLE public.conversation_session_id;

CREATE TABLE IF NOT EXISTS public.conversation_session_id (
	id serial4 NOT NULL,
	session_id int4 NOT NULL,
	last_recv_msg_id varchar NULL,
	is_kick_out bool NOT NULL DEFAULT false,
	convercation_id uuid NOT NULL,
	CONSTRAINT conversation_session_id_pk PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS conversation_session_id_convercation_id_idx ON public.conversation_session_id USING btree (convercation_id);
CREATE INDEX IF NOT EXISTS conversation_session_id_session_id_idx ON public.conversation_session_id USING btree (session_id);
-- public.sessions definition

-- Drop table

-- DROP TABLE public.sessions;

CREATE TABLE IF NOT EXISTS public.sessions (
	id serial4 NOT NULL,
	"session" varchar NOT NULL,
	CONSTRAINT sessions_pk PRIMARY KEY (id)
);

-- public.messages definition

-- Drop table

-- DROP TABLE public.messages;

CREATE TABLE IF NOT EXISTS public.messages (
	msg_id varchar NOT NULL,
	convercation_id uuid NOT NULL,
	from_session int4 NOT NULL,
	send_time timestamp NOT NULL,
	status int4 NOT NULL DEFAULT 0,
	"type" int4 NOT NULL DEFAULT 0,
	"content" varchar NULL,
	CONSTRAINT messages_pk PRIMARY KEY (msg_id)
);
CREATE INDEX IF NOT EXISTS messages_convercation_id_idx ON public.messages USING btree (convercation_id);
CREATE INDEX IF NOT EXISTS messages_type_idx ON public.messages USING btree (type);