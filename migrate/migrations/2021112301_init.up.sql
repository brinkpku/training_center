

CREATE TABLE IF NOT EXISTS workers(
    worker_id text not null, 
    config text not null
    );

INSERT INTO `workers` VALUES
    ("test_1","config_1"),
    ("test_1", "config_2");
