package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mehradSaeedi/beauty-store/internal/database"
	"github.com/mehradSaeedi/beauty-store/internal/models"
	"github.com/mehradSaeedi/beauty-store/internal/routes"
)

func main() {

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, db)

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	fmt.Println("Server running on localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}
