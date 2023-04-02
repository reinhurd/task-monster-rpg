-- init.sql

-- Create the 'topic' table
CREATE TABLE topic
(
    id bigserial PRIMARY KEY NOT NULL,
    main_theme text,
    topic text
);

-- Create the 'players' table
CREATE TABLE players
(
    id bigserial PRIMARY KEY NOT NULL,
    name text,
    token text,
    currentTask text,
    level int,
    xp int,
    health int
);