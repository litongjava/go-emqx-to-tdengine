# go-emqx-to-tdengine

从emqx订阅消息并插入到td engine

## 主要功能
- 从emqx订阅消息并插入到td engine
- eqmx 自动重连
- 支持任意json格式
- td engine 自动重连
- 插入语句支持sql模块

## 注意
如果传输的报文是二进制格式,暂不支持,需要自行修改源码进行解码


## 如何使用
config.toml
```shell
emqxBroker = "192.168.3.9"
emqxPort = 1883
emqxUsername = "emqx"
emqxPassword = "public"
emqxTopic = "sensor/data"
tdDriverName ="taosSql"
tdDataSourceName= "root:taosdata@tcp(192.168.3.9:6030)/test"
createSql= "create.sql"
insertSql= "insert.sql"
```

create.sql
```shell
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
```
insert.sql
```shell
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
```
启动程序测试
```shell
go-emqx-to-tdengine
```

## 使用mockjs发送模拟消息
```shell script
yarn add mqtt mockjs
node mock.js
```

```javascript
// mock.js
const mqtt = require('mqtt');
const Mock = require('mockjs');

const EMQX_SERVER = 'mqtt://192.168.3.9:1883';
const topic = 'sensor/data';
const CLIENT_NUM = 100;
const STEP = 5000; // 模拟采集时间间隔 ms
const AWAIT = 5000; // 每次发送完后休眠时间，防止消息速率过快 ms
const CLIENT_POOL = [];

startMock();


function sleep(timer = 100) {
  return new Promise(resolve => {
    setTimeout(resolve, timer)
  })
}

async function startMock() {
  const now = Date.now();
  for (let i = 0; i < CLIENT_NUM; i++) {
    const client = await createClient(`mock_client_${i}`);
    CLIENT_POOL.push(client)
  }
  // last 24h every 5s
  const last = 24 * 3600 * 1000;
  for (let ts = now - last; ts <= now; ts += STEP) {
    for (const client of CLIENT_POOL) {
      const mockData = generateMockData();
      const data = {
        ...mockData,
        id: client.options.clientId,
        ts,
      };

      client.publish(topic, JSON.stringify(data))
    }
    const dateStr = new Date(ts).toLocaleTimeString();
    console.log(`${dateStr} send success.`);
    await sleep(AWAIT)
  }
  console.log(`Done, use ${(Date.now() - now) / 1000}s`)
}

/**
 * Init a virtual mqtt client
 * @param {string} clientId ClientID
 */
function createClient(clientId) {
  return new Promise((resolve, reject) => {
    const client = mqtt.connect(EMQX_SERVER, {
      clientId,
    });
    client.on('connect', () => {
      console.log(`client ${clientId} connected`);
      resolve(client)
    });
    client.on('reconnect', () => {
      console.log('reconnect')
    });
    client.on('error', (e) => {
      console.error(e);
      reject(e)
    })
  })
}

/**
 * Generate mock data
 */
function generateMockData() {
  return {
    "temperature": parseFloat(Mock.Random.float(22, 100).toFixed(2)),
    "humidity": parseFloat(Mock.Random.float(12, 86).toFixed(2)),
    "volume": parseFloat(Mock.Random.float(20, 200).toFixed(2)),
    "PM10": parseFloat(Mock.Random.float(0, 300).toFixed(2)),
    "pm25": parseFloat(Mock.Random.float(0, 300).toFixed(2)),
    "SO2": parseFloat(Mock.Random.float(0, 50).toFixed(2)),
    "NO2": parseFloat(Mock.Random.float(0, 50).toFixed(2)),
    "CO": parseFloat(Mock.Random.float(0, 50).toFixed(2)),
    "area": Mock.Random.integer(0, 20),
    "ts": 1596157444170,
  }
}
```

## docker
```shell
docker build -t litongjava/go-emqx-to-tdengine .
```