package backend

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Thread struct {
	Id   string
	Name string
}

type Message struct {
	Sender string
	Time   string
	Text   string
}

var (
	db        *sql.DB
	connected bool
)

func OpenDatabase(file string) (err error) {
	if connected {
		db.Close()
		connected = false
	}

	db, err = sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		fmt.Println(err)
		db.Close()
		return
	}

	connected = true
	return
}

func Threads() (threads []Thread) {
	if !connected {
		return
	}

	rows, err := db.Query("select threadId from threads order by lastEventTimestamp DESC")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err == nil {

			rows, err := db.Query("select normalizedId from thread_participants where threadId = ?", id)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var p string
				if err = rows.Scan(&p); err == nil {
					threads = append(threads, Thread{id, p})
				}
			}
		}
	}
	return
}

func Messages(id string) (messages []Message) {
	if !connected {
		return
	}

	rows, err := db.Query("select senderId, timestamp, message from text_events where threadId = ? order by timestamp ASC", id)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var sender, time, text string
		if err = rows.Scan(&sender, &time, &text); err == nil {
			messages = append(messages, Message{sender, time, text})
		}
	}
	return
}
