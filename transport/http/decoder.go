package http

import (
	"encoding/json"
	"io"

	"github.com/jawacompu10/juice_shop/recipe_store/models"
)

func decodeRecipeFromJSON(r io.Reader) (models.Recipe, error) {
	recipe := models.Recipe{}
	if err := json.NewDecoder(r).Decode(&recipe); err != nil {
		return recipe, err
	}
	return recipe, nil
}
