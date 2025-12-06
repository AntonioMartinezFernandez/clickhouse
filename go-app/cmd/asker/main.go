package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	ctx := context.Background()
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

	err := Query(ctx, conn)
	if err != nil {
		fmt.Println("error querying data", err)
		os.Exit(1)
	}
}

func Query(ctx context.Context, conn *sql.DB) error {
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	queries := []string{
		// `SELECT * FROM sensor_data WHERE read_value > 5 AND read_value < 10 LIMIT 5`,
		`SELECT AVG(read_value) AS avg_value FROM sensor_data`,
		`SELECT MAX(read_value) AS max_value FROM sensor_data`,
		`SELECT MIN(read_value) AS min_value FROM sensor_data`,
	}

	for _, query := range queries {
		start := time.Now()

		rows, err := conn.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return err
		}

		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {
			if err := rows.Scan(valuePtrs...); err != nil {
				return err
			}

			for i, col := range columns {
				fmt.Printf("%s: %v\t", col, values[i])
			}
			fmt.Println()
		}

		fmt.Println("query duration: ", time.Since(start))
	}

	return nil
}
