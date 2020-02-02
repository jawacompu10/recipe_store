package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jawacompu10/juice_shop/recipe_store/transport"
)

// Transport implements the HTTP transport layer
type Transport struct {
	business transport.Business
}

// New creates and returns a new value for the http transport
func New(business transport.Business) transport.Transport {
	return &Transport{
		business: business,
	}
}

// Start starts listening to requests
func (ht *Transport) Start(opts map[string]string) error {
	r := ht.buildRouter()
	return http.ListenAndServe(":8080", r)
}

func (ht *Transport) buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/recipe/{item}", ht.GetRecipe).Methods("GET")

	return r
}

// GetRecipe returns a recipe for the given item
func (ht *Transport) GetRecipe(w http.ResponseWriter, req *http.Request) {
	var item string
	var ok bool
	if item, ok = mux.Vars(req)["item"]; !ok {
		http.Error(w, "Specify a valid item", http.StatusBadRequest)
		return
	}
	recipe, err := ht.business.GetRecipe(item)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

// AddNewRecipe adds a new recipe
func (ht *Transport) AddNewRecipe(w http.ResponseWriter, req *http.Request) {
	recipe, err := decodeRecipeFromJSON(req.Body)
	if err != nil {
		http.Error(w, "Invalid request:"+err.Error(), http.StatusBadRequest)
		return
	}
	recipe, err = ht.business.AddNewRecipe(recipe)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

// UpdateRecipe adds a new recipe
func (ht *Transport) UpdateRecipe(w http.ResponseWriter, req *http.Request) {
	recipe, err := decodeRecipeFromJSON(req.Body)
	if err != nil {
		http.Error(w, "Invalid request:"+err.Error(), http.StatusBadRequest)
		return
	}
	recipe, err = ht.business.UpdateRecipe(recipe)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

func getStatusCode(err error) int {
	code := http.StatusInternalServerError
	if sc, ok := getErrAsStatusCoder(err); ok {
		code = sc.GetStatusCode()
	}
	return code
}

// StatusCoder defines a method that returns a HTTP
// status code for an error
type StatusCoder interface {
	GetStatusCode() int
}

func getErrAsStatusCoder(err error) (StatusCoder, bool) {
	var i interface{} = err
	v, ok := i.(StatusCoder)
	return v, ok
}
