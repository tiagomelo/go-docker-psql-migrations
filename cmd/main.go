// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/tiagomelo/go-docker-psql-migrations/config"
	"github.com/tiagomelo/go-docker-psql-migrations/db"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(errors.Wrap(err, "reading config"))
		os.Exit(1)
	}
	db, err := db.Connect(cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresDb)
	if err != nil {
		fmt.Println(errors.Wrap(err, "connecting to db"))
		os.Exit(1)
	}
	query := `
        SELECT
            table_name,
            column_name
        FROM
            information_schema.columns
        WHERE
            table_schema = 'public'
        ORDER BY
            table_name,
            ordinal_position;
    `
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(errors.Wrap(err, "executing query"))
		os.Exit(1)
	}
	defer rows.Close()
	fmt.Println("+----------------------+--------------------+")
	fmt.Println("| Table Name           | Column Name        |")
	fmt.Println("+----------------------+--------------------+")
	var tableName, columnName string
	for rows.Next() {
		err := rows.Scan(&tableName, &columnName)
		if err != nil {
			fmt.Println(errors.Wrap(err, "scanning query results"))
			os.Exit(1)
		}
		fmt.Printf("| %-20s | %-18s |\n", tableName, columnName)
	}
	fmt.Println("+----------------+--------------------------+")
	if err = rows.Err(); err != nil {
		fmt.Println(errors.Wrap(err, "looping through returned rows"))
		os.Exit(1)
	}
}
