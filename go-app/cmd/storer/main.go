package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-faker/faker/v4"
)

var (
	totalEntries = 10_000_000
	batchSize    = 1_0000
)

func main() {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "supersecret",
		},
		Debug: false,
	})

	if err := conn.Ping(); err != nil {
		fmt.Println("error pinging clickhouse", err)
		os.Exit(1)
	}

	defer conn.Close()

	if err := Preparation(conn); err != nil {
		fmt.Println("error in preparation", err)
		os.Exit(1)
	}

	numBatches := totalEntries / batchSize
	for i := range numBatches {
		startID := i * batchSize
		if err := BatchInsert(conn, startID); err != nil {
			fmt.Println("error inserting batch data", err)
			os.Exit(1)
		}
	}

	fmt.Println("batch insert completed successfully")
}

func Preparation(conn *sql.DB) error {
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	/*
		Try engine alternatives: MergeTree (static data -OLAP-), ReplacingMergeTree (dynamic data -OLTP-) and AgregatingMergeTree (aggregated data)

		Analyze different data types. Keep an eye in the cardinality

		Check how efficiently order the columns
	*/
	queries := []string{
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
	}

	for _, query := range queries {
		if _, err := conn.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func BatchInsert(conn *sql.DB, startID int) error {
	scope, err := conn.Begin()
	if err != nil {
		return err
	}

	batch, err := scope.Prepare("INSERT INTO sensor_data (id, read_time, sensor_location, read_value) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	for i := range batchSize {
		var data SensorData
		if err := faker.FakeData(&data); err != nil {
			return err
		}
		if _, err := batch.Exec(
			startID+i,
			time.Now().Add(time.Duration(data.RandomNumber)*time.Hour),
			data.SensorLocation,
			data.ReadValue,
		); err != nil {
			return err
		}
	}

	return scope.Commit()
}

type SensorData struct {
	RandomNumber   int     `faker:"boundary_start=1, boundary_end=1000"`
	SensorLocation string  `faker:"word"`
	ReadValue      float64 `faker:"boundary_start=-10, boundary_end=50"`
}
