package data

import (
	"context"
	"macro-tracker/data/config"
	"macro-tracker/models"
	"time"
)

func InsertMeasurement(m models.Measurement) (wasUpdate bool, err error) {
	// Vérifie si une mesure existe déjà pour ce user à cette date
	var exists bool
	err = config.DB.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM measurements WHERE user_id = $1 AND date = $2
		)
	`, m.UserID, m.Date).Scan(&exists)
	if err != nil {
		return false, err
	}

	// Effectue l’upsert
	_, err = config.DB.Exec(context.Background(), `
		INSERT INTO measurements (user_id, date, weight, body_fat, note)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, date) DO UPDATE SET
			weight = EXCLUDED.weight,
			body_fat = EXCLUDED.body_fat,
			note = EXCLUDED.note
	`, m.UserID, m.Date, m.Weight, m.BodyFat, m.Note)
	if err != nil {
		return false, err
	}

	return exists, nil // si existed déjà → c’est un update
}

func GetMeasurements(userID int) ([]models.Measurement, error) {
	rows, err := config.DB.Query(context.Background(), `
		SELECT id, date, weight, body_fat, note
		FROM measurements
		WHERE user_id = $1
		ORDER BY date DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Measurement
	for rows.Next() {
		var m models.Measurement
		m.UserID = userID
		if err := rows.Scan(&m.ID, &m.Date, &m.Weight, &m.BodyFat, &m.Note); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

func GetLastMeasurementDate(userID int) (time.Time, error) {
	var lastDate time.Time
	err := config.DB.QueryRow(context.Background(), `
		SELECT date FROM measurements
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT 1
	`, userID).Scan(&lastDate)
	return lastDate, err
}
