package storage

type User struct {
	Id        int64  `db:"id"`
	ChatId    int64  `db:"chat_id"`
	LastVideo string `db:"last_video"`
	StartTime *int   `db:"start_time"`
	EndTime   *int   `db:"end_time"`
	UserName  string `db:"user_name"`
}

type Message struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
	Text string `db:"text"`
}
