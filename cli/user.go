package cli

import (
	"bufio"
	"fmt"
	"macro-tracker/data"
	"os"
	"strconv"
	"strings"
)

func HandleUserCommand(args []string) {
	if len(args) == 0 {
		handleUserView()
		return
	}

	switch args[0] {
	case "edit":
		handleUserEdit()
	default:
		fmt.Println("Commande inconnue : user", args[0])
	}
}

// Affiche le profil utilisateur
func handleUserView() {
	user, err := data.CreateOrGetUser()
	if err != nil {
		fmt.Println("Erreur lors de la lecture du profil :", err)
		return
	}

	fmt.Println("Profil utilisateur")
	fmt.Printf("Nom        : %s %s\n", user.FirstName, user.LastName)
	fmt.Printf("Âge        : %d ans\n", user.Age)
	fmt.Printf("Genre      : %s\n", user.Gender)
	fmt.Printf("Poids      : %.1f kg\n", user.WeightKg)
	fmt.Printf("Taille     : %.1f cm\n", user.HeightCm)
	fmt.Printf("Taux de masse grasse   : %.1f %%\n", user.BodyFat)
	fmt.Printf("Objectifs  : %.0f%% glucides / %.0f%% protéines / %.0f%% lipides\n",
		user.CarbRatio*100, user.ProteinRatio*100, user.FatRatio*100)
	fmt.Printf("IMC (Indice de Masse Corporelle) : %.2f\n", user.BMI())
	fmt.Printf("Objectif : %s\n", user.Goal)

}

// Permet de modifier les champs un à un
func handleUserEdit() {
	user, err := data.CreateOrGetUser()
	if err != nil {
		fmt.Println("Impossible de charger le profil :", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Édition du profil : appuie sur [Entrée] pour conserver la valeur actuelle.")

	user.FirstName = askString(reader, "Prénom", user.FirstName)
	user.LastName = askString(reader, "Nom", user.LastName)
	user.Age = askInt(reader, "Âge", user.Age)
	gender := askString(reader, "Genre (homme / femme)", user.Gender)
	gender = strings.ToLower(strings.TrimSpace(gender))
	if gender != "homme" && gender != "femme" {
		fmt.Println("Genre invalide. Défini comme 'homme'.")
		gender = "homme"
	}
	user.Gender = gender
	user.WeightKg = askFloat(reader, "Poids (kg)", user.WeightKg)
	user.HeightCm = askFloat(reader, "Taille (cm)", user.HeightCm)
	user.BodyFat = askFloat(reader, "Taux de masse grasse (%) [facultatif]", user.BodyFat)
	user.CarbRatio = askFloat(reader, "Ratio glucides (ex: 0.4) [facultatif]", user.CarbRatio)
	user.ProteinRatio = askFloat(reader, "Ratio protéines (ex: 0.3) [facultatif]", user.ProteinRatio)
	user.FatRatio = askFloat(reader, "Ratio lipides (ex: 0.3) [facultatif]", user.FatRatio)
	validGoals := map[string]bool{
		"perte":    true,
		"maintien": true,
		"prise":    true,
	}

	var goal string
	for {
		goal = askString(reader, "Objectif (perte / maintien / prise)", user.Goal)
		goal = strings.ToLower(strings.TrimSpace(goal))
		if validGoals[goal] {
			break
		}
		fmt.Println("Objectif invalide. Choisis : perte / maintien / prise")
	}

	user.Goal = goal

	err = data.UpdateUser(user)
	if err != nil {
		fmt.Println("Échec de la mise à jour :", err)
		return
	}

	fmt.Println("Profil mis à jour avec succès.")
}

func askString(reader *bufio.Reader, label string, current string) string {
	fmt.Printf("%s [%s] : ", label, current)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return current
	}
	return text
}

func askInt(reader *bufio.Reader, label string, current int) int {
	fmt.Printf("%s [%d] : ", label, current)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return current
	}
	val, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("Entrée invalide, valeur conservée.")
		return current
	}
	return val
}

func askFloat(reader *bufio.Reader, label string, current float64) float64 {
	fmt.Printf("%s [%.2f] : ", label, current)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return current
	}
	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		fmt.Println("Entrée invalide, valeur conservée.")
		return current
	}
	return val
}
