package test

import (
	"context"
	"macro-tracker/data"
	"macro-tracker/data/config"
	"macro-tracker/models"
	"testing"
)

func setupUserTest(t *testing.T) {
	err := config.InitSchema()
	if err != nil {
		t.Fatalf("InitSchema échoué : %v", err)
	}

	// Nettoyage : supprime tous les utilisateurs
	_, err = config.DB.Exec(context.Background(), "DELETE FROM users")
	if err != nil {
		t.Fatalf("Échec du nettoyage : %v", err)
	}
}

func TestSchemaAndInsert(t *testing.T) {
	// Créer la table users
	err := config.InitSchema()
	if err != nil {
		t.Fatalf("Échec InitSchema : %v", err)
	}

	// Insérer un utilisateur
	err = data.InsertTestUser()
	if err != nil {
		t.Fatalf("Échec InsertTestUser : %v", err)
	}

	// Lire les utilisateurs
	users, err := data.GetUsers()
	if err != nil {
		t.Fatalf("Échec GetUsers : %v", err)
	}

	if len(users) == 0 {
		t.Fatal("Aucun utilisateur trouvé après insertion")
	}
}

func TestCreateOrGetUser_CreatesUserIfNoneExists(t *testing.T) {
	setupUserTest(t)

	user, err := data.CreateOrGetUser()
	if err != nil {
		t.Fatalf("Erreur CreateOrGetUser : %v", err)
	}

	if user.ID == 0 {
		t.Fatal("Utilisateur non créé correctement (ID == 0)")
	}
	if user.FirstName != "Inconnu" {
		t.Errorf("Nom attendu : Inconnu, obtenu : %s", user.FirstName)
	}
}

func TestCreateOrGetUser_ReturnsExistingUser(t *testing.T) {
	setupUserTest(t)

	// Création manuelle
	err := data.InsertTestUser()
	if err != nil {
		t.Fatalf("Erreur InsertTestUser : %v", err)
	}

	user1, err := data.CreateOrGetUser()
	if err != nil {
		t.Fatalf("Erreur CreateOrGetUser : %v", err)
	}

	user2, err := data.CreateOrGetUser()
	if err != nil {
		t.Fatalf("Erreur CreateOrGetUser : %v", err)
	}

	if user1.ID != user2.ID {
		t.Error("CreateOrGetUser devrait retourner le même utilisateur existant")
	}
}

func TestUpdateUser(t *testing.T) {
	setupUserTest(t)

	// Créer un utilisateur par défaut
	user, err := data.CreateOrGetUser()
	if err != nil {
		t.Fatalf("Erreur lors de la création du user : %v", err)
	}

	// Modifier certaines valeurs
	user.FirstName = "Alice"
	user.LastName = "Test"
	user.WeightKg = 80.0
	user.HeightCm = 175.0
	user.BodyFat = 22.5
	user.CarbRatio = 0.45
	user.ProteinRatio = 0.30
	user.FatRatio = 0.25
	user.Goal = "perte"

	// Appliquer la mise à jour
	err = data.UpdateUser(user)
	if err != nil {
		t.Fatalf("Échec UpdateUser : %v", err)
	}

	// Lire à nouveau
	updated, err := data.GetUser()
	if err != nil {
		t.Fatalf("Échec GetUser après update : %v", err)
	}

	if updated.FirstName != "Alice" || updated.WeightKg != 80.0 {
		t.Errorf("Mise à jour non prise en compte. Obtenu : %+v", updated)
	}
}

func TestUser_BMI(t *testing.T) {
	user := models.User{
		WeightKg: 70.0,
		HeightCm: 175.0,
	}

	expected := 22.86
	got := user.BMI()

	// Comparaison avec arrondi à 2 décimales
	if round(got, 2) != round(expected, 2) {
		t.Errorf("IMC incorrect : attendu %.2f, obtenu %.2f", expected, got)
	}
}

func round(val float64, decimals int) float64 {
	pow := 1.0
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int(val*pow+0.5)) / pow
}

func TestUser_RatiosSumToOne(t *testing.T) {
	user := models.User{
		CarbRatio:    0.45,
		ProteinRatio: 0.30,
		FatRatio:     0.25,
	}

	total := user.CarbRatio + user.ProteinRatio + user.FatRatio

	if round(total, 2) != 1.0 {
		t.Errorf("Les ratios nutritionnels ne totalisent pas 1.0 (100%%) : total = %.2f", total)
	}
}
