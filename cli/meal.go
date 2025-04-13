package cli

import (
	"fmt"
	"macro-tracker/api"
	"macro-tracker/data"
	"macro-tracker/models"
	"strconv"
	"strings"
)

func HandleMealCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Commandes : add, meals, clear")
		return
	}

	switch args[0] {
	case "add":
		handleAddToMeal(args[1:])
	case "meals":
		handleShowMeals()
	case "clear":
		handleClearMeals()
	default:
		fmt.Println("Commande repas inconnue.")
	}
}

func handleAddToMeal(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage : add [fdcId] [grammes] [nom_repas]")
		return
	}

	// Parsing des arguments
	fdcID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("fdcId invalide")
		return
	}

	grams, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Println("Quantité invalide")
		return
	}

	mealName := strings.Join(args[2:], " ") // permet "petit déjeuner", "collation matin" etc.

	// Récupération des données nutritionnelles depuis l'API
	food := api.GetFoodAsModel(fdcID)
	if food == nil {
		fmt.Println("Impossible de récupérer l'aliment.")
		return
	}

	// Mise à l’échelle en fonction de la quantité
	mealItem := models.MealItem{
		Food:     *food,
		Quantity: grams,
	}

	// Ajout dans le repas
	data.AddToMeal(mealName, mealItem)
	fmt.Printf("✅ Ajouté : %s (%.0fg) dans %s\n", food.Description, grams, mealName)
}

func handleShowMeals() {
	day := data.GetCurrentDayLog()
	if len(day.Meals) == 0 {
		fmt.Println("Aucun repas enregistré aujourd’hui.")
		return
	}

	fmt.Println("Repas du jour :")
	for _, meal := range day.Meals {
		fmt.Printf("▶ %s\n", meal.Name)
		for _, item := range meal.Items {
			fmt.Printf("   - %s (%.0fg)\n", item.Food.Description, item.Quantity)
		}
	}
}

func handleClearMeals() {
	data.ClearCurrentDay()
	fmt.Println("Repas du jour réinitialisés.")
}

func handleEditMeal(args []string) {
	if len(args) < 3 {
		fmt.Println("Utilisation : meal edit <repas> <fdc_id> <quantité>")
		return
	}

	// Prendre les deux derniers arguments
	qte, err1 := strconv.ParseFloat(args[len(args)-1], 64)
	fdcID, err2 := strconv.Atoi(args[len(args)-2])
	mealName := strings.Join(args[:len(args)-2], " ")

	fmt.Println("mealName :", mealName)
	fmt.Println("fdcID :", fdcID)
	fmt.Println("qte :", qte)

	if err1 != nil || err2 != nil {
		fmt.Println("Arguments invalides. Vérifie le fdc_id et la quantité.")
		return
	}

	ok := data.UpdateMealItemQuantity(mealName, fdcID, qte)
	if ok {
		fmt.Printf("Quantité mise à jour : %dg pour l’aliment %d dans le repas '%s'\n", int(qte), fdcID, mealName)
	} else {
		fmt.Println("Aucun aliment trouvé dans ce repas.")
	}
}

func HandleMealShowCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Utilisation : meal show <nom>")
		return
	}

	name := strings.Join(args, " ")
	meal, err := data.LoadSavedMeal(name)
	if err != nil {
		fmt.Printf("Repas introuvable : %v\n", err)
		return
	}

	fmt.Printf("Détails du repas enregistré '%s' :\n\n", name)

	var total models.Food

	for _, item := range meal.Items {
		fmt.Printf("- [%d] %s (%.0fg) → %.0f kcal | %.1fg prot | %.1fg gluc | %.1fg lip\n",
			item.Food.FdcID,
			item.Food.Description,
			item.Quantity,
			item.Food.Calories,
			item.Food.Proteins,
			item.Food.Carbs,
			item.Food.Fats,
		)

		// Agrégation
		total.Calories += item.Food.Calories
		total.Proteins += item.Food.Proteins
		total.Carbs += item.Food.Carbs
		total.Fats += item.Food.Fats
	}

	fmt.Println("\nTotal :")
	fmt.Printf("- Calories   : %.0f kcal\n", total.Calories)
	fmt.Printf("- Protéines  : %.1f g\n", total.Proteins)
	fmt.Printf("- Glucides   : %.1f g\n", total.Carbs)
	fmt.Printf("- Lipides    : %.1f g\n", total.Fats)
}

func HandleEditSavedMealCommand(args []string) {
	if len(args) < 3 {
		fmt.Println("Utilisation : meal edit-saved <nom> <fdc_id> <quantité>")
		return
	}

	qte, err1 := strconv.ParseFloat(args[len(args)-1], 64)
	fdcID, err2 := strconv.Atoi(args[len(args)-2])
	mealName := strings.Join(args[:len(args)-2], " ")

	if err1 != nil || err2 != nil {
		fmt.Println("Arguments invalides.")
		return
	}

	err := UpdateSavedMealQuantity(mealName, fdcID, qte)
	if err != nil {
		fmt.Println("Échec :", err)
	} else {
		fmt.Printf("Repas enregistré '%s' mis à jour : %dg pour l’aliment %d\n", mealName, int(qte), fdcID)
	}
}

func HandleRemoveItemFromSavedMeal(args []string) {
	if len(args) < 2 {
		fmt.Println("Utilisation : meal removeitem-saved <nom> <fdc_id>")
		return
	}

	fdcID, err := strconv.Atoi(args[len(args)-1])
	mealName := strings.Join(args[:len(args)-1], " ")

	if err != nil {
		fmt.Println("ID d’aliment invalide.")
		return
	}

	err = data.RemoveItemFromSavedMeal(mealName, fdcID)
	if err != nil {
		fmt.Println("Échec :", err)
	} else {
		fmt.Printf("Aliment %d supprimé du repas enregistré '%s'.\n", fdcID, mealName)
	}
}

func HandleRemoveItemFromCurrentMeal(args []string) {
	if len(args) < 2 {
		fmt.Println("Utilisation : meal removeitem <repas> <fdc_id>")
		return
	}

	fdcID, err := strconv.Atoi(args[len(args)-1])
	mealName := strings.Join(args[:len(args)-1], " ")

	if err != nil {
		fmt.Println("ID d’aliment invalide.")
		return
	}

	ok := data.RemoveItemFromMeal(mealName, fdcID)
	if ok {
		fmt.Printf("Aliment %d supprimé du repas '%s'.\n", fdcID, mealName)
	} else {
		fmt.Printf("Aucun aliment %d trouvé dans le repas '%s'.\n", fdcID, mealName)
	}
}

func HandleMealRenameCommand(args []string) {
	if len(args) < 2 {
		fmt.Println("Utilisation : meal rename <ancien_nom> <nouveau_nom>")
		return
	}

	oldName := args[0]
	newName := strings.Join(args[1:], " ")

	ok := data.RenameMeal(oldName, newName)
	if ok {
		fmt.Printf("Repas '%s' renommé en '%s'\n", oldName, newName)
	} else {
		fmt.Printf("Repas '%s' introuvable\n", oldName)
	}
}
