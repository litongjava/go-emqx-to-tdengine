CREATE TABLE IF NOT EXISTS sensor_data
(
    ts          timestamp,
    temperature float,
    humidity    float,
    volume      float,
    PM10        float,
    pm25        float,
    SO2         float,
    NO2         float,
    CO          float,
    sensor_id   NCHAR(255),
    area        TINYINT,
    coll_time   timestamp
);