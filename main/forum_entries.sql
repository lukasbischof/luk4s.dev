CREATE TABLE forum_entries
(
    id               INTEGER NOT NULL
        CONSTRAINT forum_entries_pk
            PRIMARY KEY AUTOINCREMENT,
    content          TEXT    NOT NULL,
    author           TEXT    NOT NULL,
    captcha_response TEXT    NOT NULL,
    created          INTEGER NOT NULL
);

