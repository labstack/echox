// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package pgstore_test

import (
	"fmt"
	"io/fs"

	"echo.labstack.com/pro/execution/pgstore"
)

// docs:start ExampleMigrations
func ExampleMigrations() {
	entries, err := fs.ReadDir(pgstore.Migrations(), "migrations")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
	// Output:
	// 0001.sql
}

// docs:end ExampleMigrations
