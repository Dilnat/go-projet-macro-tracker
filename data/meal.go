package data

import (
	"context"
	"fmt"
	"macro-tracker/data/config"
	"macro-tracker/models"
	"time"
)

var currentDay models.DayLog = models.DayLog{
	Date:  time.Now().Format("2006-01-02"), // tu pourras automatiser ça
	Meals: []models.Meal{},
}

func AddToMeal(mealName string, item models.MealItem) {
	// Cherche le repas ou crée-le
	for i, meal := range currentDay.Meals {
		if meal.Name == mealName {
			currentDay.Meals[i].Items = append(currentDay.Meals[i].Items, models.MealItem{
				Food:     item.Food,
				Quantity: item.Quantity,
			})
			return
		}
	}

	// Nouveau repas
	currentDay.Meals = append(currentDay.Meals, models.Meal{
		Name: mealName,
		Items: []models.MealItem{
			{
				Food:     item.Food,
				Quantity: item.Quantity,
			},
		},
	})
}

func GetCurrentDayLog() models.DayLog {
	return currentDay
}

func ClearCurrentDay() {
	currentDay.Meals = []models.Meal{}
}

func GetDayMacros(day models.DayLog) models.Food {
	var total models.Food
	for _, meal := range day.Meals {
		for _, item := range meal.Items {
			total.Calories += item.Food.Calories
			total.Proteins += item.Food.Proteins
			total.Carbs += item.Food.Carbs
			total.Fats += item.Food.Fats
			total.Fibers += item.Food.Fibers
		}
	}
	return total
}

// Enregistre tous les repas d’une journée dans la base
func SaveDayLog(userID int, day models.DayLog) error {
	fmt.Printf("👤 Save → userID = %d, date = %s\n", userID, day.Date)

	tx, err := config.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// 1. Insérer la journée
	// Supprime l’ancienne journée si elle existe (et tout ce qui dépend)
	_, _ = config.DB.Exec(context.Background(), `
		DELETE FROM day_logs WHERE user_id = $1 AND date = $2
	`, userID, day.Date)
	var dayLogID int
	err = tx.QueryRow(context.Background(), `
		INSERT INTO day_logs (user_id, date) VALUES ($1, $2) RETURNING id
	`, userID, day.Date).Scan(&dayLogID)
	if err != nil {
		return fmt.Errorf("day_logs insert error: %w", err)
	}

	for _, meal := range day.Meals {
		// 2. Insérer le repas
		var mealID int
		err = tx.QueryRow(context.Background(), `
			INSERT INTO meals (day_log_id, name) VALUES ($1, $2) RETURNING id
		`, dayLogID, meal.Name).Scan(&mealID)
		fmt.Printf("🍽️ Repas inséré : %s (mealID: %d)\n", meal.Name, mealID)
		if err != nil {
			return fmt.Errorf("meal insert error: %w", err)
		}

		for _, item := range meal.Items {
			// 3. Insérer dans le cache des aliments (foods) s’il n’existe pas
			_, _ = tx.Exec(context.Background(), `
				INSERT INTO foods (fdc_id, description, calories, proteins, carbs, fats, fibers)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				ON CONFLICT (fdc_id) DO NOTHING
			`, item.Food.FdcID, item.Food.Description, item.Food.Calories, item.Food.Proteins,
				item.Food.Carbs, item.Food.Fats, item.Food.Fibers)
			fmt.Printf("   ↪ Aliment: %s (%.0fg)\n", item.Food.Description, item.Quantity)

			// 4. Insérer l’aliment dans le repas
			_, err = tx.Exec(context.Background(), `
				INSERT INTO meal_items (meal_id, food_id, quantity)
				VALUES ($1, $2, $3)
			`, mealID, item.Food.FdcID, item.Quantity)
			if err != nil {
				return fmt.Errorf("meal_items insert error: %w", err)
			}
		}
	}

	// 5. Commit
	return tx.Commit(context.Background())
}

func LoadDayLog(userID int, date string) (models.DayLog, error) {
	fmt.Printf("👤 Load → userID = %d, date = %s\n", userID, date)

	var day models.DayLog
	day.Date = date

	// 1. Récupérer l'ID de la journée
	var dayLogID int
	err := config.DB.QueryRow(context.Background(), `
		SELECT id FROM day_logs WHERE user_id = $1 AND date = $2
	`, userID, date).Scan(&dayLogID)
	if err != nil {
		return day, fmt.Errorf("day not found: %w", err)
	}

	// 2. Récupérer les repas liés à cette journée
	fmt.Printf("🔎 Requête repas pour day_log_id = %d\n", dayLogID)
	rows, err := config.DB.Query(context.Background(), `
		SELECT id, name FROM meals WHERE day_log_id = $1
	`, dayLogID)
	if err != nil {
		return day, err
	}
	defer rows.Close()

	for rows.Next() {
		var meal models.Meal
		var mealID int
		err := rows.Scan(&mealID, &meal.Name)
		if err != nil {
			return day, err
		}

		// 3. Récupérer les aliments pour chaque repas
		itemRows, err := config.DB.Query(context.Background(), `
			SELECT f.fdc_id, f.description, f.calories, f.proteins, f.carbs, f.fats, f.fibers, mi.quantity
			FROM meal_items mi
			JOIN foods f ON f.fdc_id = mi.food_id
			WHERE mi.meal_id = $1
		`, mealID)
		if err != nil {
			return day, err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item models.MealItem
			err := itemRows.Scan(
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
				return day, err
			}
			meal.Items = append(meal.Items, item)
		}

		day.Meals = append(day.Meals, meal)
	}

	return day, nil
}

func SetCurrentDayLog(day models.DayLog) {
	currentDay = day
}

func GetLastMeal() models.Meal {
	if len(currentDay.Meals) == 0 {
		return models.Meal{}
	}
	return currentDay.Meals[len(currentDay.Meals)-1]
}

func AddFullMealToCurrentDay(mealName string, meal models.Meal) {
	for _, item := range meal.Items {
		AddToMeal(mealName, item) // Ajoute chaque aliment à un repas du jour
	}
}

func RemoveMealByName(name string) bool {
	for i, meal := range currentDay.Meals {
		if meal.Name == name {
			currentDay.Meals = append(currentDay.Meals[:i], currentDay.Meals[i+1:]...)
			return true
		}
	}
	return false
}

func DeleteSavedMeal(name string) error {
	_, err := config.DB.Exec(context.Background(),
		`DELETE FROM saved_meals WHERE LOWER(name) = LOWER($1)`, name)
	return err
}

func UpdateMealItemQuantity(mealName string, fdcID int, newQuantity float64) bool {
	for i, meal := range currentDay.Meals {
		if meal.Name == mealName {
			for j, item := range meal.Items {
				if item.Food.FdcID == fdcID {
					currentDay.Meals[i].Items[j].Quantity = newQuantity
					return true
				}
			}
		}
	}
	return false
}

func RemoveItemFromMeal(mealName string, fdcID int) bool {
	for i, meal := range currentDay.Meals {
		if meal.Name == mealName {
			newItems := []models.MealItem{}
			found := false

			for _, item := range meal.Items {
				if item.Food.FdcID == fdcID {
					found = true
					continue
				}
				newItems = append(newItems, item)
			}

			if found {
				if len(newItems) == 0 {
					// 🔥 Supprime le repas entier s’il est vide
					currentDay.Meals = append(currentDay.Meals[:i], currentDay.Meals[i+1:]...)
					fmt.Printf("ℹ️ Le repas '%s' a été supprimé car il était vide.\n", mealName)
					return true
				}
				currentDay.Meals[i].Items = newItems
				return true
			}
		}
	}
	return false
}

func RenameMeal(oldName, newName string) bool {
	for i, meal := range currentDay.Meals {
		if meal.Name == oldName {
			currentDay.Meals[i].Name = newName
			return true
		}
	}
	return false
}
