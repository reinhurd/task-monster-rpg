CREATE TABLE topic
(
    id bigserial PRIMARY KEY NOT NULL,
    main_theme text,
    topic text,
);
CREATE TABLE players
(
    id bigserial PRIMARY KEY NOT NULL,
    name text,
    token text,
    currentTask text,
    level int,
    xp int,
    health int,
);