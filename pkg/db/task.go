package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (db *DB) AddTask(t *Task) (int64, error) {
	res, err := db.conn.Exec(
		"INSERT INTO scheduler(date, title, comment, repeat) VALUES(?,?,?,?)",
		t.Date, t.Title, t.Comment, t.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DB) GetTask(id int64) (*Task, error) {
	var t Task
	err := db.conn.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id,
	).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return &t, nil
}

func (db *DB) UpdateTask(t *Task) error {
	res, err := db.conn.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		t.Date, t.Title, t.Comment, t.Repeat, t.ID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (db *DB) DeleteTask(id int64) error {
	res, err := db.conn.Exec(
		"DELETE FROM scheduler WHERE id = ?",
		id,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (db *DB) Tasks(limit int) ([]*Task, error) {
	rows, err := db.conn.Query(
		"SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		list = append(list, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
