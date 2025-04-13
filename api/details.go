package api

import (
	"encoding/json"
	"fmt"
	"io"
	"macro-tracker/models"
	"net/http"
	"sync"
)

func GetFoodDetails(fdcID int) models.Food {
	endpoint := fmt.Sprintf("https://api.nal.usda.gov/fdc/v1/food/%d?api_key=%s", fdcID, apiKey)

	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("Erreur HTTP:", err)
		return models.Food{}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var details FoodDetails
	if err := json.Unmarshal(body, &details); err != nil {
		fmt.Println("Erreur JSON:", err)
		return models.Food{}
	}

	// Conversion vers models.Food
	return models.Food{
		FdcID:       details.FdcID,
		Description: details.Description,
		Calories:    details.LabelNutrients.Calories.Value,
		Proteins:    details.LabelNutrients.Protein.Value,
		Carbs:       details.LabelNutrients.Carbs.Value,
		Fats:        details.LabelNutrients.Fat.Value,
		Fibers:      details.LabelNutrients.Fiber.Value,
	}
}

func GetMultipleFoodDetails(fdcIDs []int) []models.Food {
	var wg sync.WaitGroup
	results := make(chan models.Food, len(fdcIDs))

	for _, id := range fdcIDs {
		wg.Add(1)
		go func(fdcID int) {
			defer wg.Done()
			food := GetFoodDetails(fdcID)
			results <- food
		}(id)
	}

	wg.Wait()
	close(results)

	var foods []models.Food
	for f := range results {
		foods = append(foods, f)
	}

	return foods
}
