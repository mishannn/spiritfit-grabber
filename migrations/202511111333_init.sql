-- +goose Up
CREATE TABLE club_fullness
(
    id        BIGSERIAL PRIMARY KEY,
    ts        TIMESTAMPTZ NOT NULL DEFAULT now(),
    fullness  SMALLINT NOT NULL CHECK (fullness BETWEEN 0 AND 100)
);

CREATE INDEX idx_club_fullness_ts ON club_fullness (ts);

-- +goose Down
DROP TABLE IF EXISTS club_fullness;