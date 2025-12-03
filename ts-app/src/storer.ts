import { ClickHouse } from 'clickhouse';
import { faker } from '@faker-js/faker';

console.log('Storing data ...');
const totalEntries = 100000;
const batchSize = 1000;

const clickhouse = new ClickHouse({
  url: 'http://localhost',
  port: 8123,
  debug: true,
  basicAuth: {
    username: 'default',
    password: 'supersecret',
  },
  isUseGzip: false,
  format: 'json', // "json" || "csv" || "tsv"
  config: {
    session_timeout: 60,
    output_format_json_quote_64bit_integers: 0,
    enable_http_compression: 0,
    database: 'default',
  },
});

const queries = [
  `DROP TABLE IF EXISTS sensor_data`,

  `CREATE TABLE sensor_data (
    id UInt32,
		read_time DateTime,
		sensor_location String,
		read_value Float32
  )
  ENGINE = MergeTree()
  PARTITION BY toYYYYMM(read_time)
  ORDER BY (sensor_location, read_time)
  SETTINGS index_granularity = 8192`,
];

for (const query of queries) {
  const res = await clickhouse.query(query).toPromise();
  console.log(query, res);
}

for (let start = 0; start <= totalEntries; start += batchSize) {
  const batch = [];
  for (let i = start; i < Math.min(start + batchSize, totalEntries); i++) {
    const timestamp = faker.date
      .between({ from: '2022-01-01T00:00:00Z', to: '2022-12-31T23:59:59Z' })
      .toISOString();
    batch.push({
      id: i,
      read_time: timestamp.replace(/\..*$/g, '').replace('T', ' '),
      sensor_location: faker.location.continent(),
      read_value: Number(
        faker.number.float({ min: -10, max: 45, fractionDigits: 2 }).toFixed(2),
      ),
    });
  }

  const insertQuery =
    'INSERT INTO sensor_data (id, read_time, sensor_location, read_value) VALUES';
  const values = batch
    .map(
      (entry) =>
        `(${entry.id}, '${entry.read_time}', '${entry.sensor_location}', ${entry.read_value})`,
    )
    .join(',');

  try {
    await clickhouse.query(`${insertQuery} ${values}`).toPromise();
    console.log(`Inserted batch starting at index ${start}`);
  } catch (error) {
    console.error('Error inserting batch:', error);
  }
}

console.log('Data insertion completed.');
