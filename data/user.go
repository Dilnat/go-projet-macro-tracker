package data

import (
	"context"
	"macro-tracker/data/config"
	"macro-tracker/models"
)

func InsertTestUser() error {
	_, err := config.DB.Exec(context.Background(),
		`INSERT INTO users (first_name, last_name, age, gender, weight_kg, height_cm, body_fat, carb_ratio, protein_ratio, fat_ratio, goal)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		"Jean", "Dupont", 30, "homme", 75.5, 180.0, 18.5, 0.40, 0.30, 0.30, "perte",
	)
	return err
}

func InsertUser(u models.User) error {
	_, err := config.DB.Exec(context.Background(), `
		INSERT INTO users 
		(first_name, last_name, age, gender, weight_kg, height_cm, body_fat, carb_ratio, protein_ratio, fat_ratio, goal)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, u.FirstName, u.LastName, u.Age, u.Gender, u.WeightKg, u.HeightCm, u.BodyFat, u.CarbRatio, u.ProteinRatio, u.FatRatio, u.Goal)
	return err
}

func CreateOrGetUser() (models.User, error) {
	users, err := GetUsers()
	if err != nil {
		return models.User{}, err
	}

	if len(users) > 0 {
		return users[0], nil
	}

	// Crée un user par défaut (à compléter ensuite)
	defaultUser := models.User{
		FirstName:    "Inconnu",
		LastName:     "Utilisateur",
		Age:          30,
		Gender:       "homme",
		WeightKg:     70.0,
		HeightCm:     170.0,
		BodyFat:      20.0,
		CarbRatio:    0.40,
		ProteinRatio: 0.30,
		FatRatio:     0.30,
		Goal:         "perte",
	}

	err = InsertUser(defaultUser)
	if err != nil {
		return models.User{}, err
	}

	// Recharge pour récupérer l'ID auto-généré
	return GetUser()
}

func GetUsers() ([]models.User, error) {
	rows, err := config.DB.Query(context.Background(),
		`SELECT id, first_name, last_name, age, gender, weight_kg, height_cm, body_fat, carb_ratio, protein_ratio, fat_ratio, goal
		 FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Age,
			&u.Gender,
			&u.WeightKg,
			&u.HeightCm,
			&u.BodyFat,
			&u.CarbRatio,
			&u.ProteinRatio,
			&u.FatRatio,
			&u.Goal,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUser() (models.User, error) {
	var u models.User

	query := `
		SELECT id, first_name, last_name, age, gender, weight_kg, height_cm, body_fat,
		       carb_ratio, protein_ratio, fat_ratio, goal
		FROM users
		ORDER BY id LIMIT 1;
	`

	err := config.DB.QueryRow(context.Background(), query).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Age,
		&u.Gender,
		&u.WeightKg,
		&u.HeightCm,
		&u.BodyFat,
		&u.CarbRatio,
		&u.ProteinRatio,
		&u.FatRatio,
		&u.Goal,
	)

	return u, err
}

func UpdateUser(u models.User) error {
	_, err := config.DB.Exec(context.Background(), `
		UPDATE users
		SET 
			first_name = $1,
			last_name = $2,
			age = $3,
			gender = $4,
			weight_kg = $5,
			height_cm = $6,
			body_fat = $7,
			carb_ratio = $8,
			protein_ratio = $9,
			fat_ratio = $10,
			goal = $11
		WHERE id = $12
	`, u.FirstName, u.LastName, u.Age, u.Gender, u.WeightKg, u.HeightCm, u.BodyFat, u.CarbRatio, u.ProteinRatio, u.FatRatio, u.Goal, u.ID)

	return err
}
