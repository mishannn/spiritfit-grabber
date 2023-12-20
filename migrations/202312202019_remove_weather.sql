-- +goose Up
ALTER TABLE club_fullness DROP COLUMN `Temp`;
ALTER TABLE club_fullness DROP COLUMN `FeelsLike`;
ALTER TABLE club_fullness DROP COLUMN `WindSpeed`;
ALTER TABLE club_fullness DROP COLUMN `RainLevel`;
ALTER TABLE club_fullness DROP COLUMN `SnowLevel`;
ALTER TABLE club_fullness DROP COLUMN `Pressure`;
ALTER TABLE club_fullness DROP COLUMN `Humidity`;

-- +goose Down
ALTER TABLE club_fullness ADD COLUMN `Temp` Float32;
ALTER TABLE club_fullness ADD COLUMN `FeelsLike` Float32;
ALTER TABLE club_fullness ADD COLUMN `WindSpeed` Float32;
ALTER TABLE club_fullness ADD COLUMN `RainLevel` Float32;
ALTER TABLE club_fullness ADD COLUMN `SnowLevel` Float32;
ALTER TABLE club_fullness ADD COLUMN `Pressure` UInt16;
ALTER TABLE club_fullness ADD COLUMN `Humidity` UInt8;