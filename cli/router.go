package cli

import (
	"fmt"
)

func HandleCommand(args []string) {
	args_2 := args[1:]

	switch args[0] {
	case "help":
		printHelp()
	case "user":
		HandleUserCommand(args[1:])
	case "search", "details":
		HandleFoodCommand(args)
	case "add", "meals", "clear":
		HandleMealCommand(args)
	case "meal":
		if len(args) > 1 && args[1] == "show" {
			HandleMealShowCommand(args[2:])
		} else if len(args) > 1 && args[1] == "remove" {
			HandleMealRemoveCommand(args[2:])
		} else if len(args) > 1 && args[1] == "edit" {
			handleEditMeal(args[2:])
		} else if len(args) > 1 && args[1] == "edit-saved" {
			HandleEditSavedMealCommand(args[2:])
		} else if len(args) > 1 && args[1] == "removeitem-saved" {
			HandleRemoveItemFromSavedMeal(args[2:])
		} else if len(args) > 1 && args[1] == "removeitem" {
			HandleRemoveItemFromCurrentMeal(args[2:])
		} else if len(args) > 1 && args[1] == "rename" {
			HandleMealRenameCommand(args[2:])
		} else {
			HandleSavedMealCommand(args[1:])
		}
	case "report":
		HandleReportCommand(args[1:])
	case "save", "load":
		HandleStorageCommand(args)
	case "checkin":
		if len(args_2) > 0 && args_2[0] == "now" {
			HandleCheckinCommand([]string{"now"})
		} else {
			HandleCheckinCommand(nil)
		}
	case "history":
		if len(args) > 1 && args[1] == "chart" {
			HandleHistoryChartCommand(args[2:])
		} else {
			HandleHistoryCommand(args[1:])
		}
	case "exit", "quit":
		fmt.Println("À bientôt")
		return // quitte proprement la boucle
	default:
		fmt.Println("Commande inconnue.")
	}

	// Pas de return ici → la boucle continue !

}

func printHelp() {
	fmt.Println(`
Commandes disponibles :

Recherche & aliments
  search [mot]                Rechercher un aliment (ex: search riz)
  details [fdcId]             Voir les détails d’un aliment
  add [fdcId] [qte] [repas]   Ajouter un aliment à un repas (repas = petit-déjeuner, déjeuner, etc.)

Repas & journaux
  meals                       Afficher les repas du jour
  clear                       Réinitialiser les repas du jour
  save                        Sauvegarder les repas du jour en base
  load [date]                 Recharger une journée (ex: load 2025-04-08)

Repas enregistrés
  meal save [nom]             Sauvegarder le dernier repas ajouté
  meal list                   Lister les repas enregistrés
  meal add [nom] [repas]      Ajouter un repas enregistré dans un repas du jour

Bilan & utilisateur
  report                      Afficher le bilan nutritionnel du jour
  user                        Afficher le profil utilisateur
  user edit                   Modifier les infos utilisateur (nom, poids, objectifs...)
  history                     Afficher l’historique des check-ins
  history chart               Afficher l’historique des check-ins sous forme de graphe

Suivi physique
  checkin                     Enregistrer un poids / body fat (si dernier > 7j)
  checkin now                 Forcer une saisie de mesures (même jour)

Autres
  help                        Afficher ce menu d’aide
  exit / quit                 Quitter l’application
`)
}
