CREATE TABLE IF NOT EXISTS demo
(
    id   INTEGER NOT NULL
        CONSTRAINT p_id
            PRIMARY KEY AUTOINCREMENT,
    name TEXT    NOT NULL
        CONSTRAINT u_name
            UNIQUE
                ON CONFLICT FAIL
);

INSERT INTO demo (name)
VALUES ('item1')