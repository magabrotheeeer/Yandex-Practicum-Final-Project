package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type ReqTask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const (
	DateFormat = "20060102"
)

func AddTask(task *ReqTask) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is not initialized")
	}
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to add new task")
	}
	id, err = res.LastInsertId()
	return id, err
}

func Tasks(limit int, search string) ([]*ReqTask, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	const baseSelect = `
        SELECT id, date, title, comment, repeat
        FROM scheduler
    `
	var (
		query  string
		params []any
	)

	if search != "" {
		if parsed, err := time.Parse("02.01.2006", search); err == nil {
			query = baseSelect + " WHERE date = ? ORDER BY date ASC LIMIT ?"
			params = []any{parsed.Format(DateFormat), limit}
		} else {
			query = baseSelect + " WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?"
			pat := "%" + search + "%"
			params = []any{pat, pat, limit}
		}
	} else {
		query = baseSelect + " ORDER BY date ASC LIMIT ?"
		params = []any{limit}
	}

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]*ReqTask, 0, limit)
	for rows.Next() {
		var (
			idInt   int64
			dateInt int64
			title   sql.NullString
			comment sql.NullString
			repeat  sql.NullString
		)
		if err := rows.Scan(&idInt, &dateInt, &title, &comment, &repeat); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		t := &ReqTask{
			ID:      strconv.FormatInt(idInt, 10),
			Date:    strconv.FormatInt(dateInt, 10),
			Title:   title.String,
			Comment: comment.String,
			Repeat:  repeat.String,
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

func GetTask(id string) (*ReqTask, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	query := `
        SELECT
            id,
            date,
            title,
            comment,
            repeat
        FROM scheduler
        WHERE id = ?
    `
	var t ReqTask
	var idInt, dateInt int64

	err := db.QueryRow(query, id).Scan(
		&idInt, &dateInt, &t.Title, &t.Comment, &t.Repeat,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	t.ID = strconv.FormatInt(idInt, 10)
	t.Date = strconv.FormatInt(dateInt, 10)
	return &t, nil
}

func UpdateTask(t *ReqTask) error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	query := `
        UPDATE scheduler
        SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?
    `
	res, err := db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func DeleteTask(id string) error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check deleted rows: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return fmt.Errorf("failed to update date: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check updated rows: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
