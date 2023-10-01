CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION unaccent WITH SCHEMA public;

CREATE TEXT SEARCH CONFIGURATION unaccent_dict ( COPY = simple );
ALTER TEXT SEARCH CONFIGURATION unaccent_dict
    ALTER MAPPING FOR hword, hword_part, word
    WITH unaccent, simple;

-- Create a user-defined function to generate the tsvector
CREATE OR REPLACE FUNCTION generate_fts_tokens(name varchar, nick_name varchar, stack varchar[])
RETURNS tsvector AS $$
BEGIN
    RETURN to_tsvector('unaccent_dict', name) || ' ' || 
           to_tsvector('unaccent_dict', nick_name) || ' ' || 
           to_tsvector('unaccent_dict', ARRAY_TO_STRING(stack, ','));
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE TABLE IF NOT EXISTS people
(
    id              varchar(40)     PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            varchar(32)     NOT NULL,
    nick_name       varchar(100)    NOT NULL,
    birth_date      varchar(10)     NOT NULL,
    stack           varchar[],
    fts_tokens      tsvector        GENERATED ALWAYS AS (
        generate_fts_tokens(name, nick_name, stack)
    ) STORED,
    created_at      TIMESTAMP       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS people_pkey on people USING btree (id);
CREATE INDEX IF NOT EXISTS people_fts_idx ON people USING GIN (fts_tokens);
