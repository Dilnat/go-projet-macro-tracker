package test

import (
	"context"
	"macro-tracker/data/config"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load()
	if err := config.ConnectDB(); err != nil {
		panic("Erreur de connexion DB : " + err.Error())
	}
	defer config.DB.Close()
	os.Exit(m.Run())

	// Nettoyage total de toutes les données persistées (hard reset)
	config.DB.Exec(context.Background(), `TRUNCATE meal_items, meals, day_logs, users RESTART IDENTITY CASCADE`)

}
