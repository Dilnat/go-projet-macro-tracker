package api

const apiKey = "mNEe6Sg3FtR17Sge1iaG4S9cfYvMjzwcA0w2Zt7Z"

type Food struct {
	Description string `json:"description"`
	FdcID       int    `json:"fdcId"`
	DataType    string `json:"dataType"`
}

type SearchResponse struct {
	Foods []Food `json:"foods"`
}

type Nutrient struct {
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	UnitName string  `json:"unitName"`
}

type FoodDetails struct {
	Description    string `json:"description"`
	FdcID          int    `json:"fdcId"`
	LabelNutrients struct {
		Calories struct{ Value float64 } `json:"calories"`
		Protein  struct{ Value float64 } `json:"protein"`
		Fat      struct{ Value float64 } `json:"fat"`
		Carbs    struct{ Value float64 } `json:"carbohydrates"`
		Fiber    struct{ Value float64 } `json:"fiber"`
	} `json:"labelNutrients"`
	Nutrients []Nutrient `json:"foodNutrients"` // backup au cas où labelNutrients est vide
}
