package http

import (
	"encoding/json"
	"log"
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
	port := ":8080"
	if port, ok := opts["port"]; ok {
		port = ":" + port
	}
	log.Println("recipe_store: Starting to listen on", port)
	return http.ListenAndServe(port, r)
}

func (ht *Transport) buildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/recipe/{item}", ht.GetRecipe).Methods("GET")
	r.HandleFunc("/recipe/{id}", ht.UpdateRecipe).Methods("PUT")
	r.HandleFunc("/recipe", ht.AddNewRecipe).Methods("POST")

	r.Use(addContentTypeJSON)

	return r
}

func addContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

// GetRecipe returns a recipe for the given item
func (ht *Transport) GetRecipe(w http.ResponseWriter, req *http.Request) {
	item := mux.Vars(req)["item"]
	log.Println("GetRecipe request. Item: ", item)
	recipe, err := ht.business.GetRecipe(item)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

// AddNewRecipe adds a new recipe
func (ht *Transport) AddNewRecipe(w http.ResponseWriter, req *http.Request) {
	log.Println("AddNewRecipe request")
	recipe, err := decodeRecipeFromJSON(req.Body)
	if err != nil {
		log.Println("Failed to decode JSON")
		http.Error(w, "Invalid request:"+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Request body: %+v\n", recipe)
	recipe, err = ht.business.AddNewRecipe(recipe)
	if err != nil {
		log.Println("Failed to add new recipe: ", err)
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

// UpdateRecipe adds a new recipe
func (ht *Transport) UpdateRecipe(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	log.Println("UpdateRecipe request. Id: ", id)
	recipe, err := decodeRecipeFromJSON(req.Body)
	if err != nil {
		log.Println("Failed to decode JSON")
		http.Error(w, "Invalid request:"+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Request body: %+v\n", recipe)
	log.Println("Request recipe ID:", id)
	recipe.ID = id
	recipe, err = ht.business.UpdateRecipe(recipe)
	if err != nil {
		log.Println("Failed to update recipe: ", err)
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
