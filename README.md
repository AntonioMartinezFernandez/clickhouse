# Playing with Clickhouse Database

This repository contains examples and scripts for working with ClickHouse, a fast open-source columnar database management system.

## Prerequisites

- Docker
- Golang
- NodeJS
- Basic knowledge of SQL and databases

## Resources

- [ClickHouse Official Documentation](https://clickhouse.com/docs/en/)
- [ClickHouse TS package](https://www.npmjs.com/package/clickhouse)
- [ClickHouse Go Client](https://pkg.go.dev/github.com/ClickHouse/clickhouse-go)

## Getting Started

1. Clone the repository:
   ```bash
    git clone https://github.com/AntonioMartinezFernandez/clickhouse-playground.git
    cd clickhouse-playground
   ```
2. Start ClickHouse using Docker:
   ```bash
    docker compose up -d
   ```
3. Access ClickHouse client:
   ```bash
    docker exec -it clickhouse clickhouse-client
   ```
4. Create a sample database and table:
   ```sql
    CREATE DATABASE playground_db;
    USE playground_db;
    CREATE TABLE temperatures (id UInt32, read_time DateTime, sensor_location String, read_value Float32) ENGINE = MergeTree() ORDER BY read_time;
   ```
5. Insert sample data:
   ```sql
    INSERT INTO temperatures VALUES (1, '2019-01-01 00:00:01', 'Murcia', 12.1), (2, '2019-01-01 00:00:02', 'Barcelona', 9.8), (3, '2019-01-01 00:00:03', 'Heilbronn', 7.4);
   ```
6. Query the data:
   ```sql
     SELECT * FROM temperatures;
   ```

## Commands

1. Run ClickHouse

```bash
make up
```

2. Stop ClickHouse

```bash
make down
```

3. Open ClickHouse client

```bash
make client
```

4. Run Golang script to store data in ClickHouse

```bash
make go-run-storer
```

5. Run Golang script to query data from ClickHouse

```bash
make go-run-asker
```

6. Run NodeJS script to interact with ClickHouse

```bash
make ts-run-storer
```

7. Run NodeJS script to query data from ClickHouse

```bash
make ts-run-asker
```

## Interesting topics

- Try engine alternatives: MergeTree (static data -OLAP-), ReplacingMergeTree (dynamic data -OLTP-) and AgregatingMergeTree (aggregated data)
- Analyze different data types. Keep an eye in the cardinality
- Check how efficiently order the columns
