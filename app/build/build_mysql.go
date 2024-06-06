//go:build mysql
// +build mysql

package build

import (
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
)
