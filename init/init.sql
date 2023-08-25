CREATE TABLE users
(
    id      BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL
);

CREATE TABLE segments
(
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE user_segments
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users (user_id) ON DELETE CASCADE,
    segment_id BIGINT REFERENCES segments (id) ON DELETE CASCADE
);