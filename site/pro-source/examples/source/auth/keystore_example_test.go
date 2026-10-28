// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package auth_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"echo.labstack.com/pro/auth"
)

// docs:start persistent-keystore
// exampleSQLKeys reads the latest record on every request. The application owns
// the pool and an api_keys table with id text and record jsonb columns.
// The record stores auth.StoredKey as JSON, including current RevokedAt.
type exampleSQLKeys struct{ db *sql.DB }

func (s exampleSQLKeys) Lookup(ctx context.Context, id string) (auth.StoredKey, error) {
	var raw []byte
	err := s.db.QueryRowContext(ctx,
		"SELECT record FROM api_keys WHERE id = $1", id,
	).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return auth.StoredKey{}, nil
	}
	if err != nil {
		return auth.StoredKey{}, err
	}
	var record auth.StoredKey
	if err := json.Unmarshal(raw, &record); err != nil {
		return auth.StoredKey{}, err
	}
	return record, nil
}

var _ auth.KeyStore = exampleSQLKeys{}

// ExampleKeyStore compiles persistent lookup wiring. Supply an opened pool and
// apply the application's table migration before serving; no database runs here.
func ExampleKeyStore() {
	var db *sql.DB // Your application's opened, caller-owned pool.
	scheme := auth.APIKey{Store: exampleSQLKeys{db: db}}
	// Also set Namespace, Prefix and Peppers before building auth.Config.
	_ = scheme
}

// docs:end persistent-keystore
