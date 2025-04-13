package cli

import (
	"fmt"
	"macro-tracker/data"
	"strings"
)

func HandleHistoryCommand(args []string) {
	user, err := data.CreateOrGetUser()
	if err != nil {
		fmt.Println("Impossible de récupérer le profil :", err)
		return
	}

	measurements, err := data.GetMeasurements(user.ID)
	if err != nil {
		fmt.Println("Erreur récupération historique :", err)
		return
	}

	if len(measurements) == 0 {
		fmt.Println("Aucun historique enregistré.")
		return
	}

	fmt.Println("Historique des mesures :")
	fmt.Printf("%-12s | %-8s | %-9s | %-s\n", "Date", "Poids", "BodyFat", "Note")
	fmt.Println(strings.Repeat("-", 50))

	for _, m := range measurements {
		body := "-"
		if m.BodyFat > 0 {
			body = fmt.Sprintf("%.1f %%", m.BodyFat)
		}
		fmt.Printf("%-12s | %-8.1f | %-9s | %-s\n", m.Date.Format("2006-01-02"), m.Weight, body, m.Note)
	}
}

func HandleHistoryChartCommand(args []string) {
	user, err := data.CreateOrGetUser()
	if err != nil {
		fmt.Println("Impossible de récupérer le profil :", err)
		return
	}

	measurements, err := data.GetMeasurements(user.ID)
	if err != nil {
		fmt.Println("Erreur récupération historique :", err)
		return
	}

	if len(measurements) == 0 {
		fmt.Println("Aucun historique enregistré.")
		return
	}

	fmt.Println("\nÉvolution du poids :\n")

	// Trouver le poids max pour normaliser l’échelle
	maxWeight := measurements[0].Weight
	for _, m := range measurements {
		if m.Weight > maxWeight {
			maxWeight = m.Weight
		}
	}

	for i := len(measurements) - 1; i >= 0; i-- {
		m := measurements[i]
		barLength := int((m.Weight / maxWeight) * 20)
		bar := strings.Repeat("█", barLength)
		fmt.Printf("%s | %5.1f kg | %s\n", m.Date.Format("2006-01-02"), m.Weight, bar)
	}
}
