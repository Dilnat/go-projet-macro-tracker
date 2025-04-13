package data

import (
	"context"
	"fmt"
	"macro-tracker/data/config"
	"macro-tracker/models"
)

func SaveCurrentMealAs(name string, meal models.Meal) error {
	ctx := context.Background()

	tx, err := config.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var savedMealID int
	err = tx.QueryRow(ctx, `INSERT INTO saved_meals (name) VALUES ($1) RETURNING id`, name).Scan(&savedMealID)
	if err != nil {
		return err
	}

	for _, item := range meal.Items {
		_, _ = tx.Exec(ctx, `
			INSERT INTO foods (fdc_id, description, calories, proteins, carbs, fats, fibers)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (fdc_id) DO NOTHING
		`, item.Food.FdcID, item.Food.Description, item.Food.Calories, item.Food.Proteins,
			item.Food.Carbs, item.Food.Fats, item.Food.Fibers)
		_, err := tx.Exec(ctx, `
			INSERT INTO saved_meal_items (saved_meal_id, food_id, quantity)
			VALUES ($1, $2, $3)
		`, savedMealID, item.Food.FdcID, item.Quantity)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func ListSavedMeals() ([]string, error) {
	rows, err := config.DB.Query(context.Background(), `SELECT name FROM saved_meals ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		_ = rows.Scan(&name)
		names = append(names, name)
	}
	return names, nil
}

func LoadSavedMeal(name string) (models.Meal, error) {
	var meal models.Meal
	meal.Name = name

	ctx := context.Background()
	var savedMealID int
	err := config.DB.QueryRow(ctx, `
		SELECT id FROM saved_meals WHERE LOWER(name) = LOWER($1)
	`, name).Scan(&savedMealID)
	if err != nil {
		return meal, err
	}

	rows, err := config.DB.Query(ctx, `
		SELECT f.fdc_id, f.description, f.calories, f.proteins, f.carbs, f.fats, f.fibers, i.quantity
		FROM saved_meal_items i
		JOIN foods f ON f.fdc_id = i.food_id
		WHERE i.saved_meal_id = $1
	`, savedMealID)
	if err != nil {
		return meal, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.MealItem
		err := rows.Scan(
			&item.Food.FdcID,
			&item.Food.Description,
			&item.Food.Calories,
			&item.Food.Proteins,
			&item.Food.Carbs,
			&item.Food.Fats,
			&item.Food.Fibers,
			&item.Quantity,
		)
		if err != nil {
			return meal, err
		}
		meal.Items = append(meal.Items, item)
	}

	return meal, nil
}

func RemoveItemFromSavedMeal(name string, fdcID int) error {
	// Récupère l’ID du repas enregistré
	var mealID int
	err := config.DB.QueryRow(context.Background(), `SELECT id FROM saved_meals WHERE name = $1`, name).Scan(&mealID)
	if err != nil {
		return fmt.Errorf("repas '%s' introuvable : %w", name, err)
	}

	// Supprime l'aliment ciblé
	result, err := config.DB.Exec(context.Background(), `
		DELETE FROM saved_meal_items
		WHERE saved_meal_id = $1 AND food_id = $2
	`, mealID, fdcID)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("aliment %d non trouvé dans le repas '%s'", fdcID, name)
	}

	return nil
}
