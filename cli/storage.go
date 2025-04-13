package cli

import (
	"fmt"
	"macro-tracker/data"
	"time"
)

func HandleStorageCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Commandes : save, load [date]")
		return
	}

	switch args[0] {
	case "save":
		handleSaveCommand()
	case "load":
		if len(args) < 2 {
			fmt.Println("Usage : load [YYYY-MM-DD]")
			return
		}
		handleLoadCommand(args[1])
	default:
		fmt.Println("Commande inconnue :", args[0])
	}
}

func handleSaveCommand() {
	user, _ := data.CreateOrGetUser()
	day := data.GetCurrentDayLog()
	fmt.Println("day_logs : ", day.Date)

	if day.Date == "" {
		day.Date = time.Now().Format("2006-01-02")
	}

	err := data.SaveDayLog(user.ID, day)
	if err != nil {
		fmt.Println("Erreur sauvegarde :", err)
	} else {
		fmt.Printf("Journée du %s enregistrée en base.\n", day.Date)
	}
}

func handleLoadCommand(date string) {
	user, _ := data.CreateOrGetUser()
	day, err := data.LoadDayLog(user.ID, date)
	if err != nil {
		fmt.Println("Impossible de charger la journée :", err)
		return
	}

	// Remplace currentDay en mémoire (triche simple ici)
	data.SetCurrentDayLog(day)

	fmt.Printf("Journée du %s chargée avec %d repas.\n", day.Date, len(day.Meals))
}
