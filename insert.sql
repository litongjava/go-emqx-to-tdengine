INSERT INTO test.sensor_data
VALUES (now,
        ${payload.temperature},
        ${payload.humidity},
        ${payload.volume},
        ${payload.PM10},
        ${payload.pm25},
        ${payload.SO2},
        ${payload.NO2},
        ${payload.CO},
        '${payload.id}',
        ${payload.area},
        ${payload.ts});