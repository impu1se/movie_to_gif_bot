create table users (
    id bigserial primary key,
    chat_id bigint not null,
    last_video varchar(300),
    start_time int,
    end_time int,
    user_name varchar(300)
);

create unique index users_chat_id_uindex
    on users (chat_id);

create table messages (
    id   serial primary key,
    name varchar,
    text varchar
);

insert into messages (name, text)
values
('Новая Gif', 'Продолжительность видео должно быть не более 10 сек'),
('Очистить время начала и конца', 'Время сбросилось, введите новое время'),
('start', 'Привет, я делаю гиф из видео!Отравьте любое видео, выберите время и я сделаю для вас gif'),
('create video', 'Начало и конец получены, начинаем обработку...' ||
 'Ок старт: %v, конец: %v ' ||
  'Секунду идет обработка...')
