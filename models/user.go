package models

type User struct {
	ID           int
	FirstName    string
	LastName     string
	Age          int
	Gender       string
	WeightKg     float64
	HeightCm     float64
	BodyFat      float64
	CarbRatio    float64
	ProteinRatio float64
	FatRatio     float64
	Goal         string // "perte", "maintien", "prise"
}

// IMC = poids (kg) / taille² (m)
func (u *User) BMI() float64 {
	h := u.HeightCm / 100
	return u.WeightKg / (h * h)
}

func (u *User) ComputeTargetCalories() float64 {
	bmr := 10*u.WeightKg + 6.25*u.HeightCm - 5*float64(u.Age)
	if u.Gender == "homme" {
		bmr += 5
	} else {
		bmr -= 161
	}
	tdee := bmr * 1.5

	switch u.Goal {
	case "perte":
		return tdee * 0.85
	case "prise":
		return tdee * 1.10
	default:
		return tdee
	}
}

func (u *User) EstimatedBodyFat() float64 {
	bmi := u.BMI()
	if u.Gender == "homme" {
		return 1.20*bmi + 0.23*float64(u.Age) - 16.2
	} else {
		return 1.20*bmi + 0.23*float64(u.Age) - 5.4
	}
}
