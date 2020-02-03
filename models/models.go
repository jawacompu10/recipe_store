package models

// Recipe is the struct that holds a recipe
type Recipe struct {
	ID          string         `json:"_id,omitempty" bson:"_id,omitempty"`
	ItemName    string         `json:"item_name" bson:"item_name"`
	Quantity    int            `json:"quantity" bson:"quantity"`
	Ingredients map[string]int `json:"ingredients" bson:"ingredients"`
}
