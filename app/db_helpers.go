package app

import (
	"database/sql"
	"github.com/lukasbischof/luk4s.dev/app/forum"
	"os"
	"time"
)

func IncreaseVisitorCount(db *sql.DB) (int, error) {
	visitorCount := 0
	err := db.QueryRow("UPDATE stats SET visitors = visitors + 1 RETURNING visitors").Scan(&visitorCount)

	return visitorCount, err
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

	rows, err := db.Query("SELECT id, content, author, created FROM forum_entries ORDER BY created DESC")
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

func InitializeDatabase(db *sql.DB) error {
	schemaFile, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schemaFile))

	return err
}

func SaveForumReply(db *sql.DB, reply *forum.Reply) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO forum_replies (forum_entry_id, content, created) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(int64(reply.ForumEntryId), reply.Content, time.Now().Unix())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetForumReplies(db *sql.DB, forumEntryId int) ([]*forum.Reply, error) {
	rows, err := db.Query("SELECT id, forum_entry_id, content, created FROM forum_replies WHERE forum_entry_id = ? ORDER BY created ASC", forumEntryId)
	defer rows.Close()

	if err != nil {
		return []*forum.Reply{}, err
	}

	var replies []*forum.Reply
	for rows.Next() {
		var reply forum.Reply
		var created int
		err = rows.Scan(&reply.Id, &reply.ForumEntryId, &reply.Content, &created)
		if err != nil {
			return []*forum.Reply{}, err
		}

		reply.Created = time.Unix(int64(created), 0)
		replies = append(replies, &reply)
	}

	if err = rows.Err(); err != nil {
		return []*forum.Reply{}, err
	}

	return replies, nil
}

func GetAllForumReplies(db *sql.DB) (map[int][]*forum.Reply, error) {
	rows, err := db.Query("SELECT id, forum_entry_id, content, created FROM forum_replies ORDER BY created ASC")
	defer rows.Close()

	if err != nil {
		return map[int][]*forum.Reply{}, err
	}

	repliesMap := make(map[int][]*forum.Reply)
	for rows.Next() {
		var reply forum.Reply
		var created int
		err = rows.Scan(&reply.Id, &reply.ForumEntryId, &reply.Content, &created)
		if err != nil {
			return map[int][]*forum.Reply{}, err
		}

		reply.Created = time.Unix(int64(created), 0)
		repliesMap[reply.ForumEntryId] = append(repliesMap[reply.ForumEntryId], &reply)
	}

	if err = rows.Err(); err != nil {
		return map[int][]*forum.Reply{}, err
	}

	return repliesMap, nil
}

func getForumEntriesCount(db *sql.DB) (int, error) {
	count := 0
	err := db.QueryRow("SELECT COUNT(*) FROM forum_entries").Scan(&count)

	return count, err
}
