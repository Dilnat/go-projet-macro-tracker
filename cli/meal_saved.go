package cli

import (
	"context"
	"fmt"
	"macro-tracker/api"
	"macro-tracker/data"
	"macro-tracker/data/config"
	"macro-tracker/models"
	"strings"
)

func HandleSavedMealCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Utilisation :")
		fmt.Println("  meal save <nom>             → Sauvegarde le dernier repas ajouté (repas enregistré)")
		fmt.Println("  meal list                   → Liste des repas enregistrés")
		fmt.Println("  meal add <nom> <repas>      → Ajoute un repas enregistré à la journée (ex: meal add collation rapide dîner)")
		fmt.Println("  meal remove <repas>         → Supprimer un repas du jour (ex: meal remove dîner)")
		fmt.Println("  meal remove <nom_enregistré>→ Supprimer un repas enregistré (ex: meal remove petit-déj sportif)")
		fmt.Println("  meal                        → Afficher cette aide")
		return
	}

	switch args[0] {
	case "save":
		if len(args) < 2 {
			fmt.Println("Spécifie un nom pour sauvegarder le repas.")
			return
		}
		last := data.GetLastMeal()
		if last.Name == "" || len(last.Items) == 0 {
			fmt.Println("Aucun repas à sauvegarder.")
			return
		}
		name := strings.Join(args[1:], " ")
		err := data.SaveCurrentMealAs(name, last)
		if err != nil {
			fmt.Println("Erreur sauvegarde :", err)
		} else {
			fmt.Printf("Repas '%s' enregistré.\n", name)
		}

	case "list":
		names, err := data.ListSavedMeals()
		if err != nil {
			fmt.Println("Erreur récupération des repas :", err)
			return
		}
		if len(names) == 0 {
			fmt.Println("Aucun repas enregistré.")
			return
		}
		fmt.Println("Repas enregistrés :")
		for _, name := range names {
			fmt.Println("•", name)
		}

	case "add":
		handleAddMultple(args)
	default:
		fmt.Println("Commande inconnue :", args[0])
	}
}

func HandleMealRemoveCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Utilisation : meal remove <repas>  ou  meal remove <repas enregistré>")
		return
	}

	// Recomposer le nom complet du repas (supporte noms composés)
	mealName := strings.Join(args, " ")

	// 1. Tentative de suppression dans les repas du jour
	removed := data.RemoveMealByName(mealName)
	if removed {
		fmt.Printf("Repas '%s' supprimé de la journée.\n", mealName)
		return
	}

	// 2. Tentative de suppression dans les repas enregistrés
	err := data.DeleteSavedMeal(mealName)
	if err == nil {
		fmt.Printf("Repas enregistré '%s' supprimé.\n", mealName)
		return
	}

	// 3. Rien trouvé
	fmt.Printf("Aucun repas '%s' trouvé (ni du jour, ni enregistré).\n", mealName)
}

func UpdateSavedMealQuantity(name string, fdcID int, newQty float64) error {
	meal, err := data.LoadSavedMeal(name)
	if err != nil {
		return fmt.Errorf("repas introuvable : %w", err)
	}

	found := false
	for i, item := range meal.Items {
		if item.Food.FdcID == fdcID {
			meal.Items[i].Quantity = newQty
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("aliment %d non trouvé dans '%s'", fdcID, name)
	}

	tx, err := config.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// 1. Retrouver l’ID du repas enregistré
	var mealID int
	err = tx.QueryRow(context.Background(), `SELECT id FROM saved_meals WHERE name = $1`, name).Scan(&mealID)
	if err != nil {
		return fmt.Errorf("ID introuvable pour '%s' : %w", name, err)
	}

	// 2. Supprimer les anciens items
	_, err = tx.Exec(context.Background(), `DELETE FROM saved_meal_items WHERE saved_meal_id = $1`, mealID)
	if err != nil {
		return err
	}

	// 3. Réinsérer les nouveaux items
	for _, item := range meal.Items {
		_, err := tx.Exec(context.Background(), `
			INSERT INTO saved_meal_items (saved_meal_id, food_id, quantity)
			VALUES ($1, $2, $3)
		`, mealID, item.Food.FdcID, item.Quantity)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func handleAddMultple(args []string) {
	if len(args) < 3 {
		fmt.Println("Utilisation : meal add <repas_enregistré> <nom_du_repas>")
		return
	}

	mealName := args[len(args)-1]
	savedName := strings.Join(args[1:len(args)-1], " ")

	saved, err := data.LoadSavedMeal(savedName)
	if err != nil {
		fmt.Println("Repas introuvable :", err)
		return
	}

	// Extraire tous les fdcIDs
	var fdcIDs []int
	quantities := make(map[int]float64)

	for _, item := range saved.Items {
		fdcIDs = append(fdcIDs, item.Food.FdcID)
		quantities[item.Food.FdcID] = item.Quantity
	}

	// Chargement parallèle des aliments
	foods := api.GetMultipleFoodDetails(fdcIDs)

	// Ajout + affichage
	for _, food := range foods {
		qty := quantities[food.FdcID]
		data.AddToMeal(mealName, models.MealItem{
			Food:     food,
			Quantity: qty,
		})

		fmt.Printf("%.0fg de %s ajouté à '%s'\n", qty, food.Description, mealName)
	}

	fmt.Printf("Repas enregistré '%s' ajouté à la journée dans le repas '%s'\n", savedName, mealName)
}
