CREATE EXTENSION IF NOT EXISTS PG_TRGM;

-- CREATE EXTENSION unaccent WITH SCHEMA public;

-- CREATE TEXT SEARCH CONFIGURATION unaccent_dict ( COPY = simple );
-- ALTER TEXT SEARCH CONFIGURATION unaccent_dict
--     ALTER MAPPING FOR hword, hword_part, word
--     WITH unaccent, simple;

-- -- Create a user-defined function to generate the tsvector
-- CREATE OR REPLACE FUNCTION generate_fts_tokens(name varchar, nick_name varchar, stack varchar)
-- RETURNS tsvector AS $$
-- BEGIN
--     RETURN to_tsvector('unaccent_dict', name) || ' ' || 
--            to_tsvector('unaccent_dict', nick_name) || ' ' || 
--            to_tsvector('unaccent_dict', stack);
-- END;
-- $$ LANGUAGE plpgsql IMMUTABLE;

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