CREATE TABLE IF NOT EXISTS users
(
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS segments
(
    id   BIGSERIAL PRIMARY KEY,
    slug VARCHAR(512) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_segments
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users (id) ON DELETE CASCADE,
    segment_id BIGINT REFERENCES segments (id) ON DELETE CASCADE,
    delete_at  TIMESTAMP DEFAULT NULL,
    UNIQUE (user_id, segment_id)
);
