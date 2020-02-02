package business

import (
	"github.com/jawacompu10/juice_shop/recipe_store/models"
	"github.com/jawacompu10/juice_shop/recipe_store/transport"
)

// Repository is the interface that defines the methods the
// database layer of the service should implement
type Repository interface {
	GetRecipeByItem(item string) (models.Recipe, error)
	AddRecipe(models.Recipe) (models.Recipe, error)
	UpdateRecipe(models.Recipe) (models.Recipe, error)
}

// RecipeService will implement the business layer of the service
type RecipeService struct {
	repo Repository
}

// New creates and returns a new value for the business layer implementation
func New(db Repository) transport.Business {
	return &RecipeService{db}
}

// GetRecipe gets a recipe for the given item name
func (rs *RecipeService) GetRecipe(itemName string) (models.Recipe, error) {
	return rs.repo.GetRecipeByItem(itemName)
}

// AddNewRecipe adds a new recipe
func (rs *RecipeService) AddNewRecipe(recipe models.Recipe) (models.Recipe, error) {
	return rs.repo.AddRecipe(recipe)
}

// UpdateRecipe updates an existing recipe
func (rs *RecipeService) UpdateRecipe(recipe models.Recipe) (models.Recipe, error) {
	return rs.repo.UpdateRecipe(recipe)
}
