package transport

import (
	"github.com/jawacompu10/juice_shop/recipe_store/models"
)

// Business is the interface that defines the methods
// the business layer of the service should implement
type Business interface {
	GetRecipe(string) (models.Recipe, error)
	AddNewRecipe(models.Recipe) (models.Recipe, error)
	UpdateRecipe(models.Recipe) (models.Recipe, error)
}

// Transport is the interface that defines the methods
// the transport layer for the service should implement
type Transport interface {
	// Start starts listening for requests from clients
	Start(map[string]string) error
}
