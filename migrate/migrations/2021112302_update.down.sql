CREATE TEMPORARY TABLE workers_backup(
    worker_id text not null, 
    config text not null
    );

INSERT INTO workers_backup SELECT * FROM workers; 

DROP TABLE workers;

CREATE TABLE workers(
    worker_id text not null, 
    config text not null
    );

INSERT INTO workers SELECT worker_id,config FROM workers_backup;

DROP TABLE workers_backup;
