CREATE TABLE IF NOT EXISTS forum_entries
(
    id               INTEGER NOT NULL
        CONSTRAINT forum_entries_pk
            PRIMARY KEY AUTOINCREMENT,
    content          TEXT    NOT NULL,
    author           TEXT    NOT NULL,
    captcha_response TEXT    NOT NULL,
    created          INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS forum_replies
(
    id               INTEGER NOT NULL
        CONSTRAINT forum_replies_pk
            PRIMARY KEY AUTOINCREMENT,
    forum_entry_id   INTEGER NOT NULL,
    content          TEXT    NOT NULL,
    created          INTEGER NOT NULL,
    FOREIGN KEY (forum_entry_id) REFERENCES forum_entries(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS stats
(
    visitors INTEGER NOT NULL
);
