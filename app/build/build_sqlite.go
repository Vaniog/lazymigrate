//go:build sqlite
// +build sqlite

package build

import (
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
)
