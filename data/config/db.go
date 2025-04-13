package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		getEnv("DB_USER", "macro"),
		getEnv("DB_PASSWORD", "secret"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "macrotracker"),
	)

	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("échec connexion DB : %w", err)
	}

	return DB.Ping(context.Background())
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func InitSchema() error {
	path := getProjectRoot() + "/data/config/schema.sql"
	sql, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = DB.Exec(context.Background(), string(sql))
	return err
}

func getProjectRoot() string {
	dir, _ := os.Getwd()
	for dir != "/" {
		if _, err := os.Stat(dir + "/go.mod"); err == nil {
			return dir
		}
		dir = dir + "/.."
	}
	return "." // fallback
}
