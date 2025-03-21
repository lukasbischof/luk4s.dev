package app

import (
	"database/sql"
	"github.com/lukasbischof/luk4s.dev/app/forum"
	"time"
)

const VisitorCountKey = "visitor-count"
const ForumKey = "forum-entries"

func IncreaseVisitorCount(db *sql.DB) (int, error) {
	rows, err := db.Query("UPDATE stats SET visitors = visitors + 1 RETURNING visitors")
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var visitorCount int
	for rows.Next() {
		if err = rows.Scan(&visitorCount); err != nil {
			return 0, err
		}
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}

	return visitorCount, nil
}

func SaveForumEntry(db *sql.DB, forumEntry *forum.Entry) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO forum_entries (content, author, captcha_response, created) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(forumEntry.Content, forumEntry.Author, forumEntry.CaptchaResponse, time.Now().Unix())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetForumEntries(db *sql.DB) ([]*forum.Entry, error) {
	count, err := getForumEntriesCount(db)
	if err != nil {
		return []*forum.Entry{}, err
	}

	rows, err := db.Query("SELECT id, content, author, created FROM forum_entries")
	defer rows.Close()

	if err != nil {
		return []*forum.Entry{}, err
	}

	entriesList := make([]*forum.Entry, count)
	i := 0
	for rows.Next() {
		var entry forum.Entry
		var created int
		err = rows.Scan(&entry.Id, &entry.Content, &entry.Author, &created)
		if err != nil {
			return []*forum.Entry{}, err
		}

		entry.Created = time.Unix(int64(created), 0)

		entriesList[i] = &entry
		i++
	}

	if err = rows.Err(); err != nil {
		return []*forum.Entry{}, err
	}

	return entriesList, nil
}

func DeleteForumEntry(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM forum_entries WHERE id = ?", id)
	return err
}

func getForumEntriesCount(db *sql.DB) (int, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM forum_entries")
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var count int
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}

	return count, nil
}
