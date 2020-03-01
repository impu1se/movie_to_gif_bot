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

INSERT INTO public.messages ( name, text) VALUES ('start', 'Привет, я делаю гиф из видео!Продолжительность видео должно быть не более 10 сек');
INSERT INTO public.messages ( name, text) VALUES ('Очистить время начала и конца', 'Время сбросилось, введите новое время');
INSERT INTO public.messages ( name, text) VALUES ('download error', 'Не получилось загрузить видео, попробуйте позднее');
INSERT INTO public.messages ( name, text) VALUES ('save video', 'Пожалуйста подождите, сохраняю видео...');
INSERT INTO public.messages ( name, text) VALUES ('successful download', 'Видео успешно загружено, укажите с какой секунды начать делать gif');
INSERT INTO public.messages ( name, text) VALUES ('invalid message', 'Извините, я не понимаю.
  				  Если вы хотите сделать gif, выберите нужный пункт из меню');
INSERT INTO public.messages ( name, text) VALUES ('end second', 'Теперь введите секунду окончания видео');
INSERT INTO public.messages ( name, text) VALUES ('end more start', 'Конец не должен быть меньше начала');
INSERT INTO public.messages ( name, text) VALUES ('video more 10s', 'Продолжительность gif не должно превышать 10 сек');
INSERT INTO public.messages ( name, text) VALUES ('create video', 'Начало и конец получены, начинаем обработку...');
INSERT INTO public.messages ( name, text) VALUES ('start create video', 'Обработка видео завершена...
Началось создание gif...');
INSERT INTO public.messages ( name, text) VALUES ('loading gif', 'Создание gif завершено
Загружаем gif в чат...

');
INSERT INTO public.messages ( name, text) VALUES ('Gif из нового видео', 'Отправьте любое видео, выберите время и я сделаю для вас gif');
INSERT INTO public.messages ( name, text) VALUES ('Gif из того же видео', 'Введите новое время ');