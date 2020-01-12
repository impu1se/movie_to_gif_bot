package storage

import (
	"context"

	"github.com/impu1se/movie_to_gif_bot/configs"
	"github.com/jackc/pgx/v4"
)

type Database struct {
	connect *pgx.Conn
}

func NewDb(config *configs.Config) (*Database, error) {
	conn, err := pgx.Connect(context.Background(), config.Dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &Database{conn}, nil
}

func (db *Database) CreateUser(ctx context.Context, user *User) error {
	_, err := db.connect.Exec(ctx, "insert into users (chat_id, user_name) values ($1, $2) on conflict do nothing", user.ChatId, "")
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetText(ctx context.Context, message string) (string, error) {
	var text string
	err := db.connect.QueryRow(ctx, "select text from messages where name = $1", message).Scan(&text)
	if err != nil {
		return "", err
	}
	return text, nil
}

func (db *Database) GetUser(ctx context.Context, chatId int64) (*User, error) {
	var user User
	err := db.connect.QueryRow(ctx, `select id, chat_id, last_video, start_time, end_time, user_name from users where chat_id = $1`, chatId).
		Scan(&user.Id, &user.ChatId, &user.LastVideo, &user.StartTime, &user.EndTime, &user.UserName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Database) ClearTime(ctx context.Context, chatId int64) error {
	_, err := db.connect.Exec(ctx, "update users set start_time = null, end_time = null where chat_id = $1", chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateStartTime(ctx context.Context, chatId int64, startTime int) error {
	_, err := db.connect.Exec(ctx, "update users set start_time = $1 where chat_id = $2", startTime, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateEndTime(ctx context.Context, chatId int64, endTime int) error {
	_, err := db.connect.Exec(ctx, "update users set end_time = $1 where chat_id = $2", endTime, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateLastVideo(ctx context.Context, chatId int64, lastVideo string) error {
	_, err := db.connect.Exec(ctx, "update users set last_video = $1 where chat_id = $2", lastVideo, chatId)
	if err != nil {
		return err
	}
	return nil
}
