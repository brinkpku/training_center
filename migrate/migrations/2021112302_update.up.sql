CREATE TEMPORARY TABLE workers_backup(
    worker_id text not null, 
    config text not null,
    primary key (worker_id)
    );

INSERT INTO workers_backup 
    SELECT * FROM workers 
    WHERE true
    ON CONFLICT(worker_id) DO UPDATE SET 
        config = excluded.config 
    ; 
    -- sqlite version >= 3.24.0

DROP TABLE workers;

CREATE TABLE workers(
    worker_id text not null, 
    config text not null,
    primary key (worker_id)
    );

INSERT INTO workers SELECT worker_id,config FROM workers_backup;

DROP TABLE workers_backup;
