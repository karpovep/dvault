-- USERS --
create sequence if not exists users_id_seq;

create table if not exists users
(
    id bigint primary key not null default nextval('users_id_seq'::regclass),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_pub_id text,
    username text,
    is_public boolean
);

create unique index if not exists users_user_pub_id_key on users using btree (user_pub_id);
create index if not exists idx_users_deleted_at on users using btree (deleted_at);
