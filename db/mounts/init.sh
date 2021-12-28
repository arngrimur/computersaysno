#!/bin/bash
set -e
echo "I've run" > $HOME/run.log
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS hits (ip varchar(15) PRIMARY KEY , hit_count smallint);
EOSQL