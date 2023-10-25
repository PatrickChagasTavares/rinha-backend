ALTER SYSTEM SET max_connections = 1000;

-- ALTER DATABASE rinhabackend set synchronous_commit=OFF;

-- ALTER SYSTEM SET shared_buffers TO "425MB";

CREATE EXTENSION IF NOT EXISTS PG_TRGM;
CREATE TABLE IF NOT EXISTS people
(
    id              UUID            PRIMARY KEY NOT NULL,
    name            VARCHAR(100)    NOT NULL,
    nick_name       VARCHAR(32)     UNIQUE  NOT NULL,
    birth_date      CHAR(10)        NOT NULL,
    stack           VARCHAR,
    search          TEXT,
    created_at      TIMESTAMP       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_people_trigram ON people USING gist (
        search GIST_TRGM_OPS
    );
    