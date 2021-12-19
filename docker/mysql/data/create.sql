CREATE TABLE IF NOT EXISTS hits (ip varchar(15), hit_count smallint );
CREATE UNIQUE INDEX IF NOT EXISTS hits_idx ON hits (ip);
