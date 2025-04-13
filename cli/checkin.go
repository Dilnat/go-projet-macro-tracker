package cli

import (
	"bufio"
	"fmt"
	"macro-tracker/data"
	"macro-tracker/models"
	"os"
	"strings"
	"time"
)

func HandleCheckinCommand(args []string) {
	user, err := data.CreateOrGetUser()
	if err != nil {
		fmt.Println("Impossible de récupérer le profil :", err)
		return
	}

	date := time.Now()
	force := len(args) > 0 && args[0] == "now"

	if !force {
		// Vérifie la date du dernier check-in
		last, err := data.GetLastMeasurementDate(user.ID)
		if err == nil {
			daysSince := int(date.Sub(last).Hours() / 24)
			if daysSince < 7 {
				fmt.Printf("Dernier check-in : %s (%d jours)\n", last.Format("2006-01-02"), daysSince)
				fmt.Println("Utilise `checkin now` si tu veux forcer une nouvelle saisie.")
				return
			}
		}
	}

	if force {
		fmt.Println("⚠️  Mode forcé activé : même si un check-in récent existe.")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Saisie des mesures pour le %s\n", date.Format("2006-01-02"))

	weight := askFloat(reader, "Poids (kg)", user.WeightKg)
	bodyFat := askFloat(reader, "Taux de graisse (%) [facultatif]", 0.0)
	note := askString(reader, "Note [facultatif]", "")

	m := models.Measurement{
		UserID:  user.ID,
		Date:    date,
		Weight:  weight,
		BodyFat: bodyFat,
		Note:    note,
	}

	wasUpdate, err := data.InsertMeasurement(m)
	if err != nil {
		fmt.Println("Erreur lors de l’enregistrement :", err)
		return
	}

	if wasUpdate {
		fmt.Printf("Mesure mise à jour pour le %s.\n", m.Date.Format("2006-01-02"))
	} else {
		fmt.Printf("Mesure enregistrée pour le %s.\n", m.Date.Format("2006-01-02"))
	}

	// Calcul et affichage de l’IMC
	user.WeightKg = m.Weight
	imc := user.BMI()
	fmt.Printf("🧬 IMC calculé : %.1f\n", imc)

	// Interprétation de l’IMC
	switch {
	case imc < 18.5:
		fmt.Println("Interprétation : Maigreur")
	case imc < 25:
		fmt.Println("Interprétation : Corpulence normale")
	case imc < 30:
		fmt.Println("Interprétation : Surpoids")
	case imc < 35:
		fmt.Println("Interprétation : Obésité modérée")
	case imc < 40:
		fmt.Println("Interprétation : Obésité sévère")
	default:
		fmt.Println("Interprétation : Obésité morbide")
	}

	// Estimation body fat (soit saisi, soit estimé)
	if bodyFat == 0 {
		bodyFat = user.EstimatedBodyFat()
	}

	fmt.Printf("Masse grasse estimée : %.1f %%\n", bodyFat)

	// Interprétation
	if user.Gender == "homme" {
		switch {
		case bodyFat < 6:
			fmt.Println("Interprétation : en dessous de la norme")
		case bodyFat <= 13:
			fmt.Println("Interprétation : athlétique")
		case bodyFat <= 20:
			fmt.Println("Interprétation : normal")
		case bodyFat <= 24:
			fmt.Println("Interprétation : élevé")
		default:
			fmt.Println("Interprétation : obésité")
		}
	} else if user.Gender == "femme" {
		switch {
		case bodyFat < 14:
			fmt.Println("Interprétation : en dessous de la norme")
		case bodyFat <= 20:
			fmt.Println("Interprétation : athlétique")
		case bodyFat <= 28:
			fmt.Println("Interprétation : normal")
		case bodyFat <= 32:
			fmt.Println("Interprétation : élevé")
		default:
			fmt.Println("Interprétation : obésité")
		}
	}

}

func CheckAutoCheckin() {
	user, err := data.CreateOrGetUser()
	if err != nil {
		fmt.Println("Impossible de charger l'utilisateur :", err)
		return
	}

	lastDate, err := data.GetLastMeasurementDate(user.ID)
	if err != nil {
		// Aucun enregistrement ? On propose directement
		fmt.Println("Aucun check-in précédent trouvé. Tu veux en ajouter un maintenant ? (O/N)")
	} else {
		days := int(time.Since(lastDate).Hours() / 24)

		if days < 7 {
			fmt.Printf("Dernier check-in très récent (%d jours), pas besoin d’en refaire aujourd’hui.\n", days)
			return
		}
		fmt.Printf("Dernier check-in : %s (%d jours)\n", lastDate.Format("2006-01-02"), days)
		fmt.Print("Veux-tu enregistrer une nouvelle mesure maintenant ? (O/N) ")
	}

	// Demander à l’utilisateur
	var answer string
	fmt.Scanln(&answer)
	if strings.ToLower(answer) == "o" {
		HandleCheckinCommand(nil)
	}
}
