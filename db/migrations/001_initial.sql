CREATE TABLE anime (
    id         SERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    genre      TEXT NOT NULL,
    episodes   INT NOT NULL DEFAULT 0,
    status     TEXT NOT NULL DEFAULT 'ongoing',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE studios (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE anime_studios (
    anime_id   INT NOT NULL REFERENCES anime(id) ON DELETE CASCADE,
    studio_id  INT NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
    PRIMARY KEY (anime_id, studio_id)
);
