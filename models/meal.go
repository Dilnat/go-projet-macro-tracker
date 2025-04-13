package models

type MealItem struct {
	Food     Food
	Quantity float64 // en grammes
}

// Un repas contient plusieurs aliments
type Meal struct {
	Name  string // ex: "Déjeuner"
	Items []MealItem
}

// Une journée est une liste de repas
type DayLog struct {
	Date  string // format YYYY-MM-DD
	Meals []Meal
}
