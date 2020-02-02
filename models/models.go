package models

// Recipe is the struct that holds a recipe
type Recipe struct {
	ID          string         `json:"_id"`
	ItemName    string         `json:"item_name"`
	Quantity    int            `json:"quantity"`
	Ingredients map[string]int `json:"ingredients"`
}
