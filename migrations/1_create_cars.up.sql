CREATE TABLE IF NOT EXISTS Cars
(
    Id             SERIAL8 PRIMARY KEY,
    Producer       VARCHAR(255),
    Model          VARCHAR(255),
    Engine_capacity FLOAT,
    Power          FLOAT,
    Number         VARCHAR(6) UNIQUE,
    Images_count    SMALLINT CHECK (Images_count <= 15) DEFAULT 0,
    Description    TEXT
);

CREATE INDEX IF NOT EXISTS cars_producer_idx ON Cars (Producer);
CREATE INDEX IF NOT EXISTS cars_model_idx ON Cars (Model);
CREATE INDEX IF NOT EXISTS cars_number_idx ON Cars (Number);
CREATE INDEX IF NOT EXISTS cars_engine_capacity_idx ON Cars (Engine_capacity);
CREATE INDEX IF NOT EXISTS cars_power_idx ON Cars (Power);
CREATE INDEX IF NOT EXISTS cars_images_count_idx ON Cars (Images_count);