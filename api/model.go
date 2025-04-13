package api

import (
	"encoding/json"
	"fmt"
	"io"
	"macro-tracker/models"
	"net/http"
)

func GetFoodAsModel(fdcID int) *models.Food {
	endpoint := fmt.Sprintf("https://api.nal.usda.gov/fdc/v1/food/%d?api_key=%s", fdcID, apiKey)

	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("Erreur HTTP:", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data FoodDetails
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Erreur JSON:", err)
		return nil
	}

	return &models.Food{
		FdcID:       data.FdcID,
		Description: data.Description,
		Calories:    data.LabelNutrients.Calories.Value,
		Proteins:    data.LabelNutrients.Protein.Value,
		Carbs:       data.LabelNutrients.Carbs.Value,
		Fats:        data.LabelNutrients.Fat.Value,
		Fibers:      data.LabelNutrients.Fiber.Value,
	}
}
