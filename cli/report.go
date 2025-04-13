package cli

import (
	"fmt"
	"macro-tracker/data"
	"macro-tracker/models"
	"strings"
)

func HandleReportCommand(args []string) {
	showChart := len(args) > 0 && args[0] == "--chart"

	// Charger la journée et le profil utilisateur
	day := data.GetCurrentDayLog()
	user, err := data.CreateOrGetUser()
	if err != nil {
		fmt.Println("Impossible de récupérer le profil :", err)
		return
	}

	if len(day.Meals) == 0 {
		fmt.Println("Aucun repas enregistré aujourd’hui.")
		return
	}

	// Agréger les macros de tous les repas
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

	calsFromCarbs := total.Carbs * 4
	calsFromProteins := total.Proteins * 4
	calsFromFats := total.Fats * 9
	totalMacrosCals := calsFromCarbs + calsFromProteins + calsFromFats

	ratioCarbs := (calsFromCarbs / totalMacrosCals) * 100
	ratioProteins := (calsFromProteins / totalMacrosCals) * 100
	ratioFats := (calsFromFats / totalMacrosCals) * 100

	fmt.Printf("\nBilan nutritionnel du %s\n", day.Date)

	// Affichage par repas
	fmt.Println("\nDétail des repas :")
	for _, meal := range day.Meals {
		fmt.Printf("\n▶ %s\n", meal.Name)
		for _, item := range meal.Items {
			fmt.Printf("- %-20s | %4.0f kcal | %.1fg prot | %.1fg gluc | %.1fg lip\n",
				item.Food.Description, item.Food.Calories, item.Food.Proteins, item.Food.Carbs, item.Food.Fats)
		}
	}

	// Résumé nutritionnel
	fmt.Println("\nTotaux nutritionnels :")
	fmt.Printf("- Calories  : %.0f kcal\n", total.Calories)
	fmt.Printf("- Protéines : %.1f g\n", total.Proteins)
	fmt.Printf("- Glucides  : %.1f g\n", total.Carbs)
	fmt.Printf("- Lipides   : %.1f g\n", total.Fats)

	// Objectifs
	fmt.Println("\nObjectifs :")
	fmt.Printf("- Glucides   : %.0f %%\n", user.CarbRatio*100)
	fmt.Printf("- Protéines  : %.0f %%\n", user.ProteinRatio*100)
	fmt.Printf("- Lipides    : %.0f %%\n", user.FatRatio*100)

	// IMC, objectifs, calories
	targetCalories := user.ComputeTargetCalories()
	fmt.Println("\nRapport nutritionnel :")
	fmt.Printf("Masse grasse estimée : %.0f%%\n", user.EstimatedBodyFat())
	fmt.Printf("IMC : %.1f | Objectif : %s\n", user.BMI(), user.Goal)
	fmt.Printf("Besoin estimé : %.0f kcal\n", targetCalories)
	fmt.Printf("Apport actuel : %.0f kcal\n", total.Calories)

	// Restant
	targetCarbs := (targetCalories * user.CarbRatio) / 4
	targetProteins := (targetCalories * user.ProteinRatio) / 4
	targetFats := (targetCalories * user.FatRatio) / 9
	restCarbs := targetCarbs - total.Carbs
	restProteins := targetProteins - total.Proteins
	restFats := targetFats - total.Fats
	restCalories := targetCalories - total.Calories

	fmt.Println("\nRestant pour aujourd’hui :")
	if restCalories > 0 {
		fmt.Printf("- Calories   : %.0f kcal restantes\n", restCalories)
	} else {
		fmt.Printf("- Calories   : Objectif atteint (%.0f kcal)\n", total.Calories)
	}

	if restCarbs > 0 {
		fmt.Printf("- Glucides   : %.1f g restants\n", restCarbs)
	} else {
		fmt.Printf("- Glucides   : Objectif atteint (%.1f g)\n", total.Carbs)
	}

	if restProteins > 0 {
		fmt.Printf("- Protéines  : %.1f g restants\n", restProteins)
	} else {
		fmt.Printf("- Protéines  : Objectif atteint (%.1f g)\n", total.Proteins)
	}

	if restFats > 0 {
		fmt.Printf("- Lipides    : %.1f g restants\n", restFats)
	} else {
		fmt.Printf("- Lipides    : Objectif atteint (%.1f g)\n", total.Fats)
	}

	// Graphe ASCII optionnel
	if showChart && totalMacrosCals > 0 {
		fmt.Println("\nRépartition des apports :")

		drawBar := func(label string, percent float64) {
			bars := int(percent / 5)
			barStr := strings.Repeat("█", bars)
			fmt.Printf("%-11s : %-20s %3.0f%%\n", label, barStr, percent)
		}

		drawBar("Glucides", ratioCarbs)
		drawBar("Protéines", ratioProteins)
		drawBar("Lipides", ratioFats)
	}
}
