package test

import (
	"macro-tracker/data"
	"macro-tracker/models"
	"testing"
)

func TestAddToMealAndGetDayLog(t *testing.T) {
	// Réinitialise pour éviter interférence
	data.ClearCurrentDay()

	// Aliment de test
	food := models.Food{
		FdcID:       123,
		Description: "Test Aliment",
		Calories:    100,
		Proteins:    5,
		Carbs:       10,
		Fats:        2,
		Fibers:      1,
	}
	mealItem := models.MealItem{
		Food:     food,
		Quantity: 100,
	}

	// Ajoute l’aliment dans un repas "déjeuner"
	data.AddToMeal("déjeuner", mealItem)

	// Récupère la journée
	day := data.GetCurrentDayLog()

	if len(day.Meals) != 1 {
		t.Fatalf("Nombre de repas attendu : 1, obtenu : %d", len(day.Meals))
	}

	meal := day.Meals[0]
	if meal.Name != "déjeuner" {
		t.Errorf("Nom du repas attendu : 'déjeuner', obtenu : '%s'", meal.Name)
	}

	if len(meal.Items) != 1 {
		t.Fatalf("Nombre d'aliments dans le repas attendu : 1, obtenu : %d", len(meal.Items))
	}

	item := meal.Items[0]
	if item.Food.Description != "Test Aliment" {
		t.Errorf("Nom d’aliment attendu : 'Test Aliment', obtenu : '%s'", item.Food.Description)
	}
	if item.Quantity != 100 {
		t.Errorf("Quantité attendue : 100, obtenue : %.0f", item.Quantity)
	}
}

func TestClearCurrentDay(t *testing.T) {
	// Ajoute un aliment fictif dans un repas "dîner"
	food := models.Food{
		FdcID:       456,
		Description: "Aliment à effacer",
		Calories:    150,
		Proteins:    10,
		Carbs:       15,
		Fats:        5,
		Fibers:      2,
	}
	mealItem := models.MealItem{
		Food:     food,
		Quantity: 100,
	}
	data.AddToMeal("dîner", mealItem)

	// Vérifie qu’il y a bien un repas
	day := data.GetCurrentDayLog()
	if len(day.Meals) == 0 {
		t.Fatal("Aucun repas trouvé avant reset (attendu au moins 1)")
	}

	// Clear
	data.ClearCurrentDay()

	// Vérifie que la liste est vide
	day = data.GetCurrentDayLog()
	if len(day.Meals) != 0 {
		t.Errorf("Repas encore présents après Clear (attendu 0, obtenu %d)", len(day.Meals))
	}
}

func TestGetDayMacros(t *testing.T) {
	data.ClearCurrentDay()

	food1 := models.Food{
		Description: "Pâtes",
		Calories:    350,
		Proteins:    12,
		Carbs:       70,
		Fats:        1.5,
		Fibers:      3,
	}
	food2 := models.Food{
		Description: "Poulet",
		Calories:    200,
		Proteins:    30,
		Carbs:       0,
		Fats:        5,
		Fibers:      0,
	}

	mealItem1 := models.MealItem{
		Food:     food1,
		Quantity: 100,
	}
	mealItem2 := models.MealItem{
		Food:     food2,
		Quantity: 100,
	}

	data.AddToMeal("déjeuner", mealItem1)
	data.AddToMeal("déjeuner", mealItem2)

	day := data.GetCurrentDayLog()
	total := data.GetDayMacros(day)

	if total.Calories != 550 {
		t.Errorf("Calories attendues : 550, obtenues : %.0f", total.Calories)
	}
	if total.Proteins != 42 {
		t.Errorf("Protéines attendues : 42, obtenues : %.1f", total.Proteins)
	}
}
