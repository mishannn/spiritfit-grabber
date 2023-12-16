-- +goose Up
CREATE TABLE club_fullness
(
    `DateTime` DateTime,
    `Fullness` UInt8,
    `Temp` Float32,
    `FeelsLike` Float32,
    `WindSpeed` Float32,
    `RainLevel` Float32,
    `SnowLevel` Float32,
    `Pressure` UInt16,
    `Humidity` UInt8
) ENGINE = MergeTree()
ORDER BY DateTime;

-- +goose Down
DROP TABLE club_fullness;