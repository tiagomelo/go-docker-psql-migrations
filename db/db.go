// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	InvalidTextRepresentationErrCode = "22P02"
	CheckViolationErrCode            = "23514"
)

// For ease of unit testing.
var (
	sqlOpen = sql.Open
)

// Connect establishes a connection to a PostgreSQL database.
func Connect(user, pass, host, schema string) (*sql.DB, error) {
	db, err := sqlOpen("postgres", dsn(user, pass, host, schema))
	if err != nil {
		return nil, errors.Wrap(err, "connecting to PostgreSQL")
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "pinging PostgreSQL")
	}
	return db, nil
}

// dsn returns the Data Source Name (DSN) for the PostgreSQL connection.
func dsn(customer, pass, host, schema string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", customer, pass, host, schema)
}
