package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func SearchFood(query string) {
	baseURL := "https://api.nal.usda.gov/fdc/v1/foods/search"
	params := url.Values{}
	params.Add("api_key", apiKey)

	// Construire l'URL avec la clé API
	fullURL := baseURL + "?" + params.Encode()

	// Créer le corps JSON pour la requête POST
	payload := map[string]interface{}{
		"query":    query,
		"pageSize": 5,
	}
	jsonBody, _ := json.Marshal(payload)

	// Créer la requête POST
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Erreur création requête :", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur HTTP :", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Erreur JSON :", err)
		fmt.Println("Réponse brute :", string(body))
		return
	}

	fmt.Println("Résultats pour :", query)
	if len(result.Foods) == 0 {
		fmt.Println("Aucun aliment trouvé.")
		return
	}
	for _, food := range result.Foods {
		fmt.Printf("- %s (ID: %d, Type: %s)\n", food.Description, food.FdcID, food.DataType)
	}
}
