import { ClickHouse } from 'clickhouse';

console.log('Querying data ...');

const clickhouse = new ClickHouse({
  url: 'http://localhost',
  port: 8123,
  debug: false,
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
  `SELECT * FROM sensor_data WHERE read_value > 5 and read_value < 10 LIMIT 5`,
];

for (const query of queries) {
  const res = await clickhouse.query(query).toPromise();
  console.log(query, res);
}
