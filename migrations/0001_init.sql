create table users (
    id bigserial primary key,
    chat_id bigint not null,
    last_video varchar(300),
    start_time int,
    end_time int,
    user_name varchar(300)
);

create table messages (
    id   serial primary key,
    name varchar,
    text varchar
);