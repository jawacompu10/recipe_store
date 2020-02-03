package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jawacompu10/juice_shop/recipe_store/business"
	"github.com/jawacompu10/juice_shop/recipe_store/repo"
	"github.com/jawacompu10/juice_shop/recipe_store/transport/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting recipe_store")
	errorChan := make(chan error)

	db, err := repo.New(repo.DBInfo{
		ConnectionString: "mongodb://localhost:27017/",
		DBName:           "juice_shop",
		CollectionName:   "recipe_store",
	})
	if err != nil {
		log.Fatalln("Could not connect to DB")
	}
	go handleSignals(&errorChan)
	recipeStore := business.New(db)
	httpTransport := http.New(recipeStore)
	go func() {
		errorChan <- httpTransport.Start(map[string]string{
			"port": "8080",
		})
	}()

	log.Println("Stopping recipe_store: ", <-errorChan)
}

func handleSignals(errorChan *chan error) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	*errorChan <- fmt.Errorf("Signal recieved: %s", <-signalChan)
}
