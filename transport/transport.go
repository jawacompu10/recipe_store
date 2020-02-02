package transport

// Business is the interface that defines the methods the business layer of the service should implement
type Business interface {
	GetRecipe(itemName string) (interface{} , error) 
	AddNewRecipe(interface{}) (interface{}, error)
	UpdateRecipe(interface{}) (interface{}, error)
}
