package db

import "fmt"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(t *Task) (int64, error) {
	res, err := Conn.Exec(
		"INSERT INTO scheduler(date,title,comment,repeat) VALUES(?,?,?,?)",
		t.Date, t.Title, t.Comment, t.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	var t Task
	err := Conn.QueryRow(
		"SELECT id,date,title,comment,repeat FROM scheduler WHERE id = ?",
		id,
	).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func UpdateTask(t *Task) error {
	res, err := Conn.Exec(
		"UPDATE scheduler SET date=?,title=?,comment=?,repeat=? WHERE id=?",
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
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := Conn.Exec(
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
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateDate(date, id string) error {
	res, err := Conn.Exec(
		"UPDATE scheduler SET date=? WHERE id=?",
		date, id,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := Conn.Query(
		"SELECT id,date,title,comment,repeat FROM scheduler ORDER BY date LIMIT ?",
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
	if list == nil {
		list = make([]*Task, 0)
	}
	return list, nil
}
