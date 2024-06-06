package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
)

var URL string
var Source = "migrations"

func loadEnv() {
	err := godotenv.Load(".lazymigrate")
	if err == nil {
		return
	}
	_ = godotenv.Load(".env")
}

func init() {
	loadEnv()
	loadSource()
	if loadDirect() {
		return
	}
	if err := loadPostgres(); err == nil {
		return
	}
	if err := loadMySQL(); err == nil {
		return
	}
	if err := loadSqlite(); err == nil {
		return
	}
}

func loadSource() {
	source, ok := os.LookupEnv("LAZYMIGRATE_SOURCE")
	if ok {
		Source = source
	}
}

func loadDirect() bool {
	url, ok := os.LookupEnv("LAZYMIGRATE_URL")
	if ok {
		URL = url
	}
	return ok
}

func loadMySQL() error {
	var mysql struct {
		User     string `envconfig:"USER" required:"true"`
		Password string `envconfig:"PASSWORD" required:"true"`
		Host     string `envconfig:"HOST" default:"localhost"`
		Port     uint16 `envconfig:"PORT" default:"3306"`
		Database string `envconfig:"DATABASE" required:"true"`
	}

	if err := envconfig.Process("MYSQL", &mysql); err != nil {
		return err
	}

	URL = fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%d)/%s?query",
		mysql.User,
		mysql.Password,
		mysql.Host,
		mysql.Port,
		mysql.Database,
	)
	return nil
}

func loadPostgres() error {
	var pg struct {
		User     string `envconfig:"USER" required:"true"`
		Password string `envconfig:"PASSWORD" required:"true"`
		Host     string `envconfig:"HOST" default:"localhost"`
		Port     uint16 `envconfig:"PORT" default:"5432"`
		Database string `envconfig:"DB" required:"true"`
	}

	if err := envconfig.Process("POSTGRES", &pg); err != nil {
		return err
	}

	URL = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pg.User,
		pg.Password,
		pg.Host,
		pg.Port,
		pg.Database,
	)
	return nil
}

func loadSqlite() error {
	var sqlite struct {
		Database string `envconfig:"DB" required:"true"`
	}

	if err := envconfig.Process("SQLITE", &sqlite); err != nil {
		return err
	}

	URL = fmt.Sprintf(
		"sqlite://%s",
		sqlite.Database,
	)
	return nil
}
