package main

import (
	"log"
	"net/http"

	"github.com/renatospaka/library/configs"
	"github.com/renatospaka/library/internal/entity"
	"github.com/renatospaka/library/internal/infra/database"
	"github.com/renatospaka/library/internal/infra/webservers/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.Println("iniciando a aplicação...")
	_, err := configs.LoadConfig(".")
	if err != nil {
		log.Panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	ProductHandler := handlers.NewProductHandler(productDB)

	log.Println("servidor escutando porta:", 8000)
	http.HandleFunc("/products", ProductHandler.CreateProduct)
	http.ListenAndServe(":8000", nil)
}
