USE csn_db;
CREATE TABLE IF NOT EXISTS hits (ip varchar(15), hit_count smallint, INDEX (ip));
