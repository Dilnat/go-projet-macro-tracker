package cli

import (
	"fmt"
	"macro-tracker/api"
	"strconv"
)

func HandleFoodCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Commandes disponibles : search, details")
		return
	}

	switch args[0] {
	case "search":
		handleSearchCommand(args[1:])
	case "details":
		handleDetailsCommand(args[1:])
	default:
		fmt.Println("Commande food inconnue :", args[0])
	}
}

func handleSearchCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage : search [mot-clé]")
		return
	}
	query := ""
	for _, s := range args {
		query += s + " "
	}
	api.SearchFood(query)
}

func handleDetailsCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage : details [fdcId]")
		return
	}

	fdcID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("ID invalide")
		return
	}

	food := api.GetFoodDetails(fdcID)

	if food.FdcID == 0 {
		fmt.Println("❌ Aucun aliment trouvé.")
		return
	}

	fmt.Printf("\n%s (ID: %d)\n", food.Description, food.FdcID)
	fmt.Printf("🔹 Calories   : %.0f kcal\n", food.Calories)
	fmt.Printf("🔹 Protéines  : %.1f g\n", food.Proteins)
	fmt.Printf("🔹 Glucides   : %.1f g\n", food.Carbs)
	fmt.Printf("🔹 Lipides    : %.1f g\n", food.Fats)
	fmt.Printf("🔹 Fibres     : %.1f g\n", food.Fibers)
}
