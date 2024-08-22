// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	testCases := []struct {
		name          string
		mockSqlOpen   func(driverName string, dataSourceName string) (*sql.DB, error)
		expectedError error
	}{
		{
			name: "happy path",
			mockSqlOpen: func(driverName, dataSourceName string) (*sql.DB, error) {
				db, _, _ := sqlmock.New()
				return db, nil
			},
		},
		{
			name: "error when opening conn",
			mockSqlOpen: func(driverName, dataSourceName string) (*sql.DB, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New("connecting to PostgreSQL: random error"),
		},
		{
			name: "error pinging",
			mockSqlOpen: func(driverName, dataSourceName string) (*sql.DB, error) {
				db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
				mock.ExpectPing().WillReturnError(errors.New("random error"))
				return db, nil
			},
			expectedError: errors.New("pinging PostgreSQL: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlOpen = tc.mockSqlOpen
			db, err := Connect("customer", "pass", "host", "schema")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
				require.NotNil(t, db)
			}
		})
	}
}
